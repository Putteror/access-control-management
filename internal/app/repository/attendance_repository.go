package repository

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/putteror/access-control-management/internal/app/model"
	"github.com/putteror/access-control-management/internal/app/schema"
	"gorm.io/gorm"
)

// AttendanceRepository is the interface for attendance data access.
type AttendanceRepository interface {
	GetAll(searchQuery schema.AttendanceSearchQuery) ([]model.Attendance, error)
	GetByID(id uuid.UUID) (*model.Attendance, error)
	Create(attendance *model.Attendance) error
	Update(attendance *model.Attendance) error
	Delete(id uuid.UUID) error
	IsExistName(name string, excludeID uuid.UUID) (bool, error)

	// Schedule relationship methods
	GetSchedulesByAttendanceID(attendanceID uuid.UUID) ([]model.AttendanceSchedule, error)
	CreateAttendanceSchedules(schedules []model.AttendanceSchedule) error
	DeleteAttendanceSchedulesByAttendanceID(attendanceID uuid.UUID, tx *gorm.DB) error
}

// attendanceRepositoryImpl is the implementation of AttendanceRepository.
type attendanceRepositoryImpl struct {
	db *gorm.DB
}

// NewAttendanceRepository creates a new instance of AttendanceRepository.
func NewAttendanceRepository(db *gorm.DB) AttendanceRepository {
	return &attendanceRepositoryImpl{db: db}
}

// --- CRUD Operations ---

// GetAll retrieves all attendance records with pagination.
func (r *attendanceRepositoryImpl) GetAll(searchQuery schema.AttendanceSearchQuery) ([]model.Attendance, error) {
	var attendances []model.Attendance

	query := r.db.Model(&model.Attendance{})

	if searchQuery.Name != "" {
		query = query.Where("name ILIKE ?", "%"+searchQuery.Name+"%")
	}

	var page int = searchQuery.Page
	var limit int = searchQuery.Limit
	offset := (page - 1) * limit

	// Note: ไม่มีฟิลด์ 'All' ใน AttendanceSearchQuery แต่ใช้ pagination ตาม Group repo
	if page > 0 && limit > 0 {
		query = query.Offset(offset).Limit(limit)
	}

	if err := query.Find(&attendances).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve paginated attendance records: %w", err)
	}

	return attendances, nil
}

// GetByID retrieves an attendance record by its ID.
func (r *attendanceRepositoryImpl) GetByID(id uuid.UUID) (*model.Attendance, error) {
	var attendance model.Attendance
	if err := r.db.First(&attendance, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &attendance, nil
}

// Create creates a new attendance record.
func (r *attendanceRepositoryImpl) Create(attendance *model.Attendance) error {
	return r.db.Create(attendance).Error
}

// Update updates an existing attendance record.
func (r *attendanceRepositoryImpl) Update(attendance *model.Attendance) error {
	return r.db.Save(attendance).Error
}

// Delete deletes an attendance record by its ID.
func (r *attendanceRepositoryImpl) Delete(id uuid.UUID) error {
	// ใช้ Transaction เพื่อลบข้อมูลทั้งจากตารางหลักและตารางเชื่อมโยง
	err := r.db.Transaction(func(tx *gorm.DB) error {

		// 1. ลบข้อมูลจากตารางเชื่อมโยง (AttendanceSchedule)
		if err := r.DeleteAttendanceSchedulesByAttendanceID(id, tx); err != nil {
			return err
		}

		// 2. ลบข้อมูลจากตารางหลัก (Attendance)
		// ใช้ Unscoped เพื่อให้แน่ใจว่าลบได้แม้มี Soft Delete
		if err := tx.Unscoped().Where("id = ?", id).Delete(&model.Attendance{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// IsExistName checks if an attendance record with the given name exists.
func (r *attendanceRepositoryImpl) IsExistName(name string, excludeID uuid.UUID) (bool, error) {
	var count int64
	db := r.db.Model(&model.Attendance{}).Where("name = ? AND deleted_at IS NULL", name)
	if excludeID != uuid.Nil {
		db = db.Where("id != ?", excludeID)
	}
	if err := db.Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check attendance name existence: %w", err)
	}
	return count > 0, nil
}

// --- Schedule Relationship Methods ---

// GetSchedulesByAttendanceID retrieves all AttendanceSchedule records for an attendance ID.
func (r *attendanceRepositoryImpl) GetSchedulesByAttendanceID(attendanceID uuid.UUID) ([]model.AttendanceSchedule, error) {
	var schedules []model.AttendanceSchedule
	err := r.db.Where("attendance_id = ?", attendanceID).Find(&schedules).Error
	return schedules, err
}

// CreateAttendanceSchedules inserts multiple AttendanceSchedule records.
func (r *attendanceRepositoryImpl) CreateAttendanceSchedules(schedules []model.AttendanceSchedule) error {
	if len(schedules) == 0 {
		return nil
	}
	return r.db.Create(&schedules).Error
}

// DeleteAttendanceSchedulesByAttendanceID deletes all AttendanceSchedule records for an attendance ID.
func (r *attendanceRepositoryImpl) DeleteAttendanceSchedulesByAttendanceID(attendanceID uuid.UUID, tx *gorm.DB) error {
	// ใช้ tx ที่ส่งมาเพื่อรองรับการทำ Transaction
	if tx == nil {
		tx = r.db
	}
	return tx.Unscoped().Where("attendance_id = ?", attendanceID).
		Delete(&model.AttendanceSchedule{}).Error
}
