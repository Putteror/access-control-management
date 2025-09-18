package service

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/putteror/access-control-management/internal/app/common"
	"github.com/putteror/access-control-management/internal/app/model"
	"github.com/putteror/access-control-management/internal/app/repository"
	"github.com/putteror/access-control-management/internal/app/schema"
	"gorm.io/gorm"
)

// AttendanceService defines the interface for attendance business logic.
type AttendanceService interface {
	GetAll(searchQuery schema.AttendanceSearchQuery) ([]model.Attendance, error)
	GetByID(id string) (*model.Attendance, error)
	Create(bodyRequest *schema.AttendanceRequest) (*model.Attendance, error)
	Update(id string, bodyRequest *schema.AttendanceRequest) (*model.Attendance, error)
	PartialUpdate(id string, bodyRequest *schema.AttendanceRequest) (*model.Attendance, error)
	Delete(id string) error
	ConvertToResponse(attendanceModel *model.Attendance) (*schema.AttendanceInfoResponse, error)
}

type attendanceServiceImpl struct {
	attendanceRepo repository.AttendanceRepository
	db             *gorm.DB
}

// NewAttendanceService creates a new instance of AttendanceService.
func NewAttendanceService(attendanceRepo repository.AttendanceRepository, db *gorm.DB) AttendanceService {
	return &attendanceServiceImpl{
		attendanceRepo: attendanceRepo,
		db:             db,
	}
}

// --- CRUD Operations ---

// GetAll retrieves all attendance records.
func (s *attendanceServiceImpl) GetAll(searchQuery schema.AttendanceSearchQuery) ([]model.Attendance, error) {
	return s.attendanceRepo.GetAll(searchQuery)
}

// GetByID retrieves an attendance record by its ID.
func (s *attendanceServiceImpl) GetByID(id string) (*model.Attendance, error) {
	idUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID")
	}
	return s.attendanceRepo.GetByID(idUUID)
}

// Create creates a new attendance record.
func (s *attendanceServiceImpl) Create(bodyRequest *schema.AttendanceRequest) (*model.Attendance, error) {

	// Set default value
	bodyRequest, err := s.validateAndSetDefaultValues(bodyRequest)
	if err != nil {
		return nil, err
	}
	// Validate (check duplicates)
	if err := s.validateBodyRequest(*bodyRequest, nil); err != nil {
		return nil, err
	}

	attendanceModel := &model.Attendance{
		Name: *bodyRequest.Name,
	}

	// ใช้ Transaction
	err = s.db.Transaction(func(tx *gorm.DB) error {
		txRepo := repository.NewAttendanceRepository(tx) // ต้องสร้าง txRepo ที่รับ tx
		// 1. Create Attendance
		if err := txRepo.Create(attendanceModel); err != nil {
			return fmt.Errorf("failed to create attendance: %w", err)
		}
		// 2. Create Attendance Schedules
		if len(bodyRequest.AttendanceSchedule) > 0 {
			schedules, err := s.createAttendanceScheduleModels(attendanceModel.ID.String(), bodyRequest.AttendanceSchedule)
			if err != nil {
				return err
			}
			if err := txRepo.CreateAttendanceSchedules(schedules); err != nil {
				return fmt.Errorf("failed to create attendance schedules: %w", err)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return attendanceModel, nil
}

// Update updates an existing attendance record (Full Replacement).
func (s *attendanceServiceImpl) Update(id string, bodyRequest *schema.AttendanceRequest) (*model.Attendance, error) {

	idUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID")
	}
	attendanceModel, err := s.attendanceRepo.GetByID(idUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing attendance: %w", err)
	}
	if attendanceModel == nil {
		return nil, fmt.Errorf("attendance with ID '%s' not found", id)
	}

	// Set default value
	bodyRequest, err = s.validateAndSetDefaultValues(bodyRequest)
	if err != nil {
		return nil, err
	}
	// Validate (check duplicates)
	if err := s.validateBodyRequest(*bodyRequest, attendanceModel); err != nil {
		return nil, err
	}

	// Update model
	attendanceModel.Name = *bodyRequest.Name

	err = s.db.Transaction(func(tx *gorm.DB) error {
		txRepo := repository.NewAttendanceRepository(tx)
		// 1. Update Attendance
		if err := txRepo.Update(attendanceModel); err != nil {
			return fmt.Errorf("failed to update attendance: %w", err)
		}

		// 2. Delete Schedules all & Recreate
		if err := txRepo.DeleteAttendanceSchedulesByAttendanceID(idUUID, tx); err != nil {
			return fmt.Errorf("failed to delete old attendance schedules: %w", err)
		}

		if len(bodyRequest.AttendanceSchedule) > 0 {
			schedules, err := s.createAttendanceScheduleModels(id, bodyRequest.AttendanceSchedule)
			if err != nil {
				return err
			}
			if err := txRepo.CreateAttendanceSchedules(schedules); err != nil {
				return fmt.Errorf("failed to create new attendance schedules: %w", err)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return attendanceModel, nil
}

// PartialUpdate performs a partial update on an existing attendance record.
func (s *attendanceServiceImpl) PartialUpdate(id string, bodyRequest *schema.AttendanceRequest) (*model.Attendance, error) {

	idUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID")
	}
	attendanceModel, err := s.attendanceRepo.GetByID(idUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing attendance: %w", err)
	}
	if attendanceModel == nil {
		return nil, fmt.Errorf("attendance with ID '%s' not found", id)
	}

	// Validate (check duplicates)
	// Note: ไม่เรียก validateAndSetDefaultValues เพราะไม่อยากตั้งค่า Default Schedules ถ้าไม่ได้ส่งมา
	if err := s.validateBodyRequest(*bodyRequest, attendanceModel); err != nil {
		return nil, err
	}

	// Update model Name (เฉพาะถ้ามีการส่ง Name มาและไม่ใช่ string ว่าง)
	if bodyRequest.Name != nil && *bodyRequest.Name != "" {
		attendanceModel.Name = *bodyRequest.Name
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		txRepo := repository.NewAttendanceRepository(tx)
		// 1. Update Attendance
		if err := txRepo.Update(attendanceModel); err != nil {
			return fmt.Errorf("failed to update attendance: %w", err)
		}

		// 2. Update Schedules (เฉพาะถ้ามีการส่ง AttendanceSchedule มา)
		if bodyRequest.AttendanceSchedule != nil {
			// Delete Schedules all & Recreate
			if err := txRepo.DeleteAttendanceSchedulesByAttendanceID(idUUID, tx); err != nil {
				return fmt.Errorf("failed to delete old attendance schedules: %w", err)
			}
			// Create Schedules
			if len(bodyRequest.AttendanceSchedule) > 0 {

				// ตรวจสอบและตั้งค่า Default ภายใน array schedule ที่ส่งมา (ตาม logic ใน validateAndSetDefaultValues)
				// Note: ต้องทำความสะอาด schedule array ก่อนสร้าง model
				validatedSchedules, err := s.validateAndCleanSchedules(bodyRequest.AttendanceSchedule)
				if err != nil {
					return err
				}

				schedules, err := s.createAttendanceScheduleModels(id, validatedSchedules)
				if err != nil {
					return err
				}
				if err := txRepo.CreateAttendanceSchedules(schedules); err != nil {
					return fmt.Errorf("failed to create new attendance schedules: %w", err)
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return attendanceModel, nil
}

// Delete deletes an attendance record by its ID.
func (s *attendanceServiceImpl) Delete(id string) error {
	idUUID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid ID")
	}

	// ตรวจสอบการมีอยู่ก่อนลบ (ตาม pattern ของ AccessControlGroupService)
	_, err = s.attendanceRepo.GetByID(idUUID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("attendance with ID '%s' not found", id)
		}
		return fmt.Errorf("failed to get attendance by ID: %w", err)
	}

	// การลบ Transaction ถูกจัดการภายใน AttendanceRepository.Delete แล้ว
	return s.attendanceRepo.Delete(idUUID)
}

// ConvertToResponse converts an attendance model to a response schema.
func (s *attendanceServiceImpl) ConvertToResponse(attendanceModel *model.Attendance) (*schema.AttendanceInfoResponse, error) {

	// 1. ดึง Schedule Info
	schedules, err := s.attendanceRepo.GetSchedulesByAttendanceID(attendanceModel.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get attendance schedules: %w", err)
	}

	scheduleResponses := make([]schema.AttendanceScheduleResponse, 0, len(schedules))
	for _, schedule := range schedules {
		scheduleResponses = append(scheduleResponses, schema.AttendanceScheduleResponse{
			ID:              schedule.ID.String(),
			DayOfWeek:       schedule.DayOfWeek,
			Date:            schedule.Date,
			StartTime:       schedule.StartTime,
			EndTime:         schedule.EndTime,
			EarlyInMinutes:  schedule.EarlyInMinutes,
			LateInMinutes:   schedule.LateInMinutes,
			EarlyOutMinutes: schedule.EarlyOutMinutes,
			LateOutMinutes:  schedule.LateOutMinutes,
		})
	}

	// 2. สร้าง Response
	response := &schema.AttendanceInfoResponse{
		ID:                  attendanceModel.ID.String(),
		Name:                attendanceModel.Name,
		AttendanceSchedules: scheduleResponses,

		// Note: หากมีการเพิ่มฟิลด์ AttendanceSchedules ใน schema.AttendanceInfoResponse ให้เพิ่มการ Map ที่นี่
	}

	return response, nil
}

// ----------> INNER FUNCTION <-----------------------//

// createAttendanceScheduleModels converts AttendanceScheduleRequest to AttendanceSchedule model.
func (s *attendanceServiceImpl) createAttendanceScheduleModels(attendanceID string, schedules []schema.AttendanceScheduleRequest) ([]model.AttendanceSchedule, error) {
	var attendanceSchedules []model.AttendanceSchedule
	for _, schedule := range schedules {

		// Note: ในโค้ดนี้เราจะสมมติว่า schedules ถูก validate และ clean มาแล้วจาก validateAndCleanSchedules หรือ validateAndSetDefaultValues

		// กำหนดค่าจาก Pointer
		var datePtr *string
		if schedule.Date != nil && *schedule.Date != "" {
			// ใช้ค่าจาก Request ถ้าไม่เป็น nil และไม่ใช่ string ว่าง
			datePtr = schedule.Date
		}

		// กำหนดค่าจาก Pointer (0 ถ้า nil)
		earlyInMinutes := 0
		if schedule.EarlyInMinutes != nil {
			earlyInMinutes = *schedule.EarlyInMinutes
		}
		lateInMinutes := 0
		if schedule.LateInMinutes != nil {
			lateInMinutes = *schedule.LateInMinutes
		}
		earlyOutMinutes := 0
		if schedule.EarlyOutMinutes != nil {
			earlyOutMinutes = *schedule.EarlyOutMinutes
		}
		lateOutMinutes := 0
		if schedule.LateOutMinutes != nil {
			lateOutMinutes = *schedule.LateOutMinutes
		}

		attendanceSchedules = append(attendanceSchedules, model.AttendanceSchedule{
			AttendanceID:    attendanceID,
			DayOfWeek:       *schedule.DayOfWeek,
			Date:            datePtr,
			StartTime:       *schedule.StartTime,
			EndTime:         *schedule.EndTime,
			EarlyInMinutes:  earlyInMinutes,
			LateInMinutes:   lateInMinutes,
			EarlyOutMinutes: earlyOutMinutes,
			LateOutMinutes:  lateOutMinutes,
		})
	}
	return attendanceSchedules, nil
}

// validateAndCleanSchedules performs validation and sets default values for schedule array only.
func (s *attendanceServiceImpl) validateAndCleanSchedules(schedules []schema.AttendanceScheduleRequest) ([]schema.AttendanceScheduleRequest, error) {

	// สร้างสำเนาเพื่อไม่ให้แก้ไข Request ที่ถูกส่งมาโดยตรง
	cleanedSchedules := make([]schema.AttendanceScheduleRequest, len(schedules))
	copy(cleanedSchedules, schedules)

	// ค่า Default (ใช้ new(string) เพื่อสร้าง pointer ไปยังค่า default)
	defaultStartTime := common.DefaultAttendanceStartTime
	defaultEndTime := common.DefaultAttendanceEndTime
	defaultZero := 0 // สำหรับสร้าง pointer ไปยัง 0

	for i, schedule := range cleanedSchedules {
		// 1. ตรวจสอบ DayOfWeek (required)
		if schedule.DayOfWeek == nil || *schedule.DayOfWeek < 1 || *schedule.DayOfWeek > 7 {
			return nil, fmt.Errorf("invalid day of week for schedule index %d", i)
		}

		// 2. ตรวจสอบ StartTime/EndTime (required แต่ตั้งค่า default ได้)
		if schedule.StartTime == nil || *schedule.StartTime == "" {
			cleanedSchedules[i].StartTime = &defaultStartTime
		}
		if schedule.EndTime == nil || *schedule.EndTime == "" {
			cleanedSchedules[i].EndTime = &defaultEndTime
		}

		// 3. ตั้งค่า Default สำหรับ Minutes (0 นาที ถ้าเป็น nil)
		if schedule.EarlyInMinutes == nil {
			cleanedSchedules[i].EarlyInMinutes = &defaultZero
		}
		if schedule.LateInMinutes == nil {
			cleanedSchedules[i].LateInMinutes = &defaultZero
		}
		if schedule.EarlyOutMinutes == nil {
			cleanedSchedules[i].EarlyOutMinutes = &defaultZero
		}
		if schedule.LateOutMinutes == nil {
			cleanedSchedules[i].LateOutMinutes = &defaultZero
		}
	}
	return cleanedSchedules, nil
}

// Validate and set default values
func (s *attendanceServiceImpl) validateAndSetDefaultValues(bodyRequest *schema.AttendanceRequest) (*schema.AttendanceRequest, error) {
	// 1. ตรวจสอบ Name
	if bodyRequest.Name == nil || *bodyRequest.Name == "" {
		return nil, fmt.Errorf("attendance name cannot be empty")
	}

	// 2. ตั้งค่า Default Schedules (24/7) หากไม่ได้ส่งมา
	if bodyRequest.AttendanceSchedule == nil {
		defaultSchedules := []schema.AttendanceScheduleRequest{}
		startTime := common.DefaultAttendanceStartTime
		endTime := common.DefaultAttendanceEndTime
		zero := 0 // ตัวแปรสำหรับ Pointer ที่มีค่าเป็น 0

		// วนลูปสร้าง Schedule สำหรับ DayOfWeek 1 (จันทร์) ถึง 7 (อาทิตย์)
		for day := 1; day <= 7; day++ {
			d := day

			defaultSchedules = append(defaultSchedules, schema.AttendanceScheduleRequest{
				DayOfWeek: &d,
				Date:      nil,
				StartTime: &startTime,
				EndTime:   &endTime,
				// ตั้งค่าความยืดหยุ่นทั้งหมดเป็น 0 นาที
				EarlyInMinutes:  &zero,
				LateInMinutes:   &zero,
				EarlyOutMinutes: &zero,
				LateOutMinutes:  &zero,
			})
		}
		bodyRequest.AttendanceSchedule = defaultSchedules
	} else {
		// 3. ถ้ามีการส่ง Schedules มา ให้ตรวจสอบและตั้งค่า Default ภายใน Array
		validatedSchedules, err := s.validateAndCleanSchedules(bodyRequest.AttendanceSchedule)
		if err != nil {
			return nil, err
		}
		bodyRequest.AttendanceSchedule = validatedSchedules
	}

	return bodyRequest, nil
}

// validateBodyRequest checks for duplicate attendance name.
func (s *attendanceServiceImpl) validateBodyRequest(bodyRequest schema.AttendanceRequest, attendanceModel *model.Attendance) error {

	excludeID := uuid.Nil
	if attendanceModel != nil {
		excludeID = attendanceModel.ID
	}

	// Check duplicate name
	isExistName, err := s.attendanceRepo.IsExistName(*bodyRequest.Name, excludeID)
	if err != nil {
		return fmt.Errorf("failed to check attendance name existence: %w", err)
	}
	if isExistName {
		return fmt.Errorf("attendance name is already exist")
	}

	return nil
}
