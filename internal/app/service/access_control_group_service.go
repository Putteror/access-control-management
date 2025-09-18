package service

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/putteror/access-control-management/internal/app/model"
	"github.com/putteror/access-control-management/internal/app/repository"
	"github.com/putteror/access-control-management/internal/app/schema"
	"gorm.io/gorm"
)

// AccessControlGroupService defines the interface for access control group business logic.
type AccessControlGroupService interface {
	GetAll(searchQuery schema.AccessControlGroupSearchQuery) ([]model.AccessControlGroup, error)
	GetByID(id string) (*model.AccessControlGroup, error)
	Create(bodyRequest *schema.AccessControlGroupRequest) (*model.AccessControlGroup, error)
	Update(id string, bodyRequest *schema.AccessControlGroupRequest) (*model.AccessControlGroup, error)
	PartialUpdate(id string, bodyRequest *schema.AccessControlGroupRequest) (*model.AccessControlGroup, error)
	Delete(id string) error
	ConvertToResponse(groupModel *model.AccessControlGroup) (*schema.AccessControlGroupResponse, error)
	GetDevicesInfo(deviceIDs []string) ([]schema.AccessControlDeviceInfoResponse, error)
}

type accessControlGroupServiceImpl struct {
	accessControlGroupRepo  repository.AccessControlGroupRepository
	accessControlDeviceRepo repository.AccessControlDeviceRepository
	db                      *gorm.DB
}

// NewAccessControlGroupService creates a new instance of AccessControlGroupService.
func NewAccessControlGroupService(accessControlGroupRepo repository.AccessControlGroupRepository, accessControlDeviceRepo repository.AccessControlDeviceRepository, db *gorm.DB) AccessControlGroupService {
	return &accessControlGroupServiceImpl{
		accessControlGroupRepo:  accessControlGroupRepo,
		accessControlDeviceRepo: accessControlDeviceRepo,
		db:                      db,
	}
}

// GetAll retrieves all access control groups.
func (s *accessControlGroupServiceImpl) GetAll(searchQuery schema.AccessControlGroupSearchQuery) ([]model.AccessControlGroup, error) {
	return s.accessControlGroupRepo.GetAll(searchQuery)
}

// GetByID retrieves an access control group by its ID.
func (s *accessControlGroupServiceImpl) GetByID(id string) (*model.AccessControlGroup, error) {
	id_uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID")
	}
	return s.accessControlGroupRepo.GetByID(id_uuid)
}

// Create creates a new access control group.
func (s *accessControlGroupServiceImpl) Create(bodyRequest *schema.AccessControlGroupRequest) (*model.AccessControlGroup, error) {

	// Set default value
	bodyRequest, err := s.validateAndSetDefaultValues(bodyRequest)
	if err != nil {
		return nil, err
	}
	// Validate (Set default values, check duplicates)
	if err := s.validateBodyRequest(*bodyRequest, nil); err != nil {
		return nil, err
	}

	groupModel := &model.AccessControlGroup{
		Name: *bodyRequest.Name,
	}

	// ใช้ Transaction เพื่อให้แน่ใจว่าทั้ง Group และ Device ถูกสร้างหรือยกเลิกพร้อมกัน
	err = s.db.Transaction(func(tx *gorm.DB) error {
		txRepo := repository.NewAccessControlGroupRepository(tx)
		// 1. Create Group
		if err := txRepo.Create(groupModel); err != nil {
			return fmt.Errorf("failed to create group: %w", err)
		}
		// 2. Create Group Devices
		if len(bodyRequest.AccessControlDeviceIDs) > 0 {
			groupDevices, err := s.createGroupDeviceModels(groupModel.ID.String(), bodyRequest.AccessControlDeviceIDs)
			if err != nil {
				return err
			}
			// Add group device
			if err := txRepo.CreateGroupDevices(groupDevices); err != nil {
				return fmt.Errorf("failed to create group devices: %w", err)
			}
		}
		// 3. Create Group Schedules (ส่วนที่เพิ่มเข้ามา)
		if len(bodyRequest.AccessControlGroupSchedules) > 0 {
			// A. สร้าง Model สำหรับ Schedule
			groupSchedules, err := s.createGroupScheduleModels(groupModel.ID.String(), bodyRequest.AccessControlGroupSchedules)
			if err != nil {
				return err
			}
			// B. บันทึก Model ลงใน DB (ต้องเพิ่ม Method นี้ใน Repository)
			// *สมมติว่า accessControlGroupRepo มี Method CreateGroupSchedules*
			if err := txRepo.CreateAccessControlGroupSchedule(groupSchedules); err != nil {
				return fmt.Errorf("failed to create group schedules: %w", err)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return groupModel, nil
}

// Update updates an existing access control group.
func (s *accessControlGroupServiceImpl) Update(id string, bodyRequest *schema.AccessControlGroupRequest) (*model.AccessControlGroup, error) {

	id_uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID")
	}
	groupModel, err := s.accessControlGroupRepo.GetByID(id_uuid)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing group: %w", err)
	}
	if groupModel == nil {
		return nil, fmt.Errorf("group with ID '%s' not found", id)
	}

	// Set default value
	bodyRequest, err = s.validateAndSetDefaultValues(bodyRequest)
	if err != nil {
		return nil, err
	}
	// Validate (check duplicates)
	if err := s.validateBodyRequest(*bodyRequest, groupModel); err != nil {
		return nil, err
	}

	// Update model
	groupModel.Name = *bodyRequest.Name

	err = s.db.Transaction(func(tx *gorm.DB) error {
		txRepo := repository.NewAccessControlGroupRepository(tx)
		// 1. Update Group (ใช้ Repository)
		if err := txRepo.Update(groupModel); err != nil {
			return fmt.Errorf("failed to update group: %w", err)
		}
		// Delete GroupDevice all
		if err := txRepo.DeleteGroupDevicesByGroupID(id_uuid, tx); err != nil {
			return fmt.Errorf("failed to delete old group devices: %w", err)
		}
		// Create GroupDevice
		if len(bodyRequest.AccessControlDeviceIDs) > 0 {
			groupDevices, err := s.createGroupDeviceModels(id, bodyRequest.AccessControlDeviceIDs)
			if err != nil {
				return err
			}
			if err := txRepo.CreateGroupDevices(groupDevices); err != nil {
				return fmt.Errorf("failed to create new group devices: %w", err)
			}
		}
		// 3. [SCHEDULES] Delete GroupSchedule all & Recreate <--- ส่วนที่เพิ่ม
		if err := txRepo.DeleteAccessControlGroupScheduleByGroupID(id_uuid, tx); err != nil { // ต้องมี Method นี้ใน Repo
			return fmt.Errorf("failed to delete old group schedules: %w", err)
		}
		if len(bodyRequest.AccessControlGroupSchedules) > 0 {
			groupSchedules, err := s.createGroupScheduleModels(id, bodyRequest.AccessControlGroupSchedules)
			if err != nil {
				return err
			}
			if err := txRepo.CreateAccessControlGroupSchedule(groupSchedules); err != nil { // ใช้ txRepo
				return fmt.Errorf("failed to create new group schedules: %w", err)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return groupModel, nil
}

// PartialUpdate performs a partial update on an existing access control group.
func (s *accessControlGroupServiceImpl) PartialUpdate(id string, bodyRequest *schema.AccessControlGroupRequest) (*model.AccessControlGroup, error) {

	id_uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID")
	}
	groupModel, err := s.accessControlGroupRepo.GetByID(id_uuid)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing group: %w", err)
	}
	if groupModel == nil {
		return nil, fmt.Errorf("group with ID '%s' not found", id)
	}

	// Validate (check duplicates)
	if err := s.validateBodyRequest(*bodyRequest, groupModel); err != nil {
		return nil, err
	}

	// Update model
	if bodyRequest.Name != nil {
		groupModel.Name = *bodyRequest.Name
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		txRepo := repository.NewAccessControlGroupRepository(tx)
		// 1. Update Group (ใช้ Repository)
		if err := txRepo.Update(groupModel); err != nil {
			return fmt.Errorf("failed to update group: %w", err)
		}
		// 2. Update Group Devices (เฉพาะถ้ามีการส่ง AccessControlDeviceIDs มาใน body request)
		if bodyRequest.AccessControlDeviceIDs != nil {
			// Delete GroupDevice all
			if err := txRepo.DeleteGroupDevicesByGroupID(id_uuid, tx); err != nil {
				return fmt.Errorf("failed to delete old group devices: %w", err)
			}
			// Create GroupDevice
			if len(bodyRequest.AccessControlDeviceIDs) > 0 {
				groupDevices, err := s.createGroupDeviceModels(id, bodyRequest.AccessControlDeviceIDs)
				if err != nil {
					return err
				}
				if err := txRepo.CreateGroupDevices(groupDevices); err != nil {
					return fmt.Errorf("failed to create new group devices: %w", err)
				}
			}
		}
		// 3. [SCHEDULES] Update Group Schedules (เฉพาะถ้ามีการส่ง AccessControlGroupSchedules มา) <--- ส่วนที่เพิ่ม
		if bodyRequest.AccessControlGroupSchedules != nil {
			// Delete GroupSchedule all
			if err := txRepo.DeleteAccessControlGroupScheduleByGroupID(id_uuid, tx); err != nil { // ต้องมี Method นี้ใน Repo
				return fmt.Errorf("failed to delete old group schedules: %w", err)
			}
			// Create GroupSchedule
			if len(bodyRequest.AccessControlGroupSchedules) > 0 {
				groupSchedules, err := s.createGroupScheduleModels(id, bodyRequest.AccessControlGroupSchedules)
				if err != nil {
					return err
				}
				if err := txRepo.CreateAccessControlGroupSchedule(groupSchedules); err != nil { // ใช้ txRepo
					return fmt.Errorf("failed to create new group schedules: %w", err)
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return groupModel, nil
}

// Delete deletes an access control group by its ID.
func (s *accessControlGroupServiceImpl) Delete(id string) error {
	id_uuid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid ID")
	}
	_, err = s.accessControlGroupRepo.GetByID(id_uuid)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("group with ID '%s' not found", id)
		}
		return fmt.Errorf("failed to get group by ID: %w", err)
	}

	// การลบ Transaction ถูกจัดการภายใน AccessControlGroupRepository.Delete แล้ว
	return s.accessControlGroupRepo.Delete(id_uuid)
}

// ConvertToResponse converts a group model to a response schema.
func (s *accessControlGroupServiceImpl) ConvertToResponse(groupModel *model.AccessControlGroup) (*schema.AccessControlGroupResponse, error) {

	// 1. ดึง Device IDs ที่ผูกกับ Group นี้
	deviceIDs, err := s.accessControlGroupRepo.GetDeviceIDsByGroupID(groupModel.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get device IDs for group: %w", err)
	}

	// 2. ดึงข้อมูลอุปกรณ์เต็ม
	deviceResponses, err := s.GetDevicesInfo(deviceIDs)
	if err != nil {
		return nil, err
	}

	scheduleResponses, err := s.GetGroupScheduleInfo(groupModel.ID.String())
	if err != nil {
		return nil, err
	}

	// 3. สร้าง Response

	response := &schema.AccessControlGroupResponse{
		ID:                          groupModel.ID.String(),
		Name:                        groupModel.Name,
		AccessControlDevices:        deviceResponses,
		AccessControlGroupSchedules: scheduleResponses,
	}

	return response, nil
}

// ----------> RELATE FUNCTION <-----------------------//

// GetDevicesInfo retrieves full device info from device IDs.
func (s *accessControlGroupServiceImpl) GetDevicesInfo(deviceIDs []string) ([]schema.AccessControlDeviceInfoResponse, error) {
	if len(deviceIDs) == 0 {
		return []schema.AccessControlDeviceInfoResponse{}, nil
	}

	var devicesInfo []schema.AccessControlDeviceInfoResponse

	for _, idStr := range deviceIDs {
		id_uuid, err := uuid.Parse(idStr)
		if err != nil {
			// อาจจะ log warning แทนการ return error หาก ID ไม่ถูกต้อง
			continue
		}
		device, err := s.accessControlDeviceRepo.GetByID(id_uuid)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				// อุปกรณ์ไม่พบแล้ว ข้ามไป
				continue
			}
			return nil, fmt.Errorf("failed to get device by ID '%s': %w", idStr, err)
		}
		devicesInfo = append(devicesInfo, schema.AccessControlDeviceInfoResponse{
			ID:          device.ID.String(),
			Name:        device.Name,
			HostAddress: device.HostAddress,
		})
	}
	return devicesInfo, nil
}

// GetGroupScheduleIno
func (s *accessControlGroupServiceImpl) GetGroupScheduleInfo(groupID string) ([]schema.AccessControlGroupScheduleResponse, error) {
	groupSchedules, err := s.accessControlGroupRepo.GetAccessControlGroupScheduleByGroupID(groupID)
	if err != nil {
		return nil, err
	}

	var response []schema.AccessControlGroupScheduleResponse
	for _, schedule := range groupSchedules {
		response = append(response, schema.AccessControlGroupScheduleResponse{
			DayOfWeek: schedule.DayOfWeek,
			Date:      schedule.Date,
			StartTime: schedule.StartTime,
			EndTime:   schedule.EndTime,
		})
	}
	return response, nil
}

// ----------> INNER FUNCTION <-----------------------//

// createGroupDeviceModels converts device IDs to AccessControlGroupDevice models.
func (s *accessControlGroupServiceImpl) createGroupDeviceModels(groupID string, deviceIDs []string) ([]model.AccessControlGroupDevice, error) {
	var groupDevices []model.AccessControlGroupDevice
	for _, deviceID := range deviceIDs {
		// ตรวจสอบว่า Device ID นั้นมีอยู่จริงหรือไม่ (optional แต่แนะนำ)
		device_uuid, err := uuid.Parse(deviceID)
		if err != nil {
			return nil, fmt.Errorf("invalid device ID format: %s", deviceID)
		}
		_, err = s.accessControlDeviceRepo.GetByID(device_uuid)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, fmt.Errorf("device with ID '%s' not found", deviceID)
			}
			return nil, fmt.Errorf("failed to check device existence: %w", err)
		}

		groupDevices = append(groupDevices, model.AccessControlGroupDevice{
			AccessControlGroupID:  groupID,
			AccessControlDeviceID: deviceID,
		})
	}
	return groupDevices, nil
}

// createGroupScheduleModels converts AccessControlGroupScheduleRequest to AccessControlGroupSchedule model
func (s *accessControlGroupServiceImpl) createGroupScheduleModels(groupID string, schedules []schema.AccessControlGroupScheduleRequest) ([]model.AccessControlGroupSchedule, error) {
	var groupSchedules []model.AccessControlGroupSchedule
	for _, schedule := range schedules {
		var datePtr *string
		// 💡 การแก้ไข: เช็กค่าว่างและกำหนดให้เป็น nil ถ้าเป็นสตริงว่าง
		if schedule.Date != nil {
			datePtr = schedule.Date
		}
		groupSchedules = append(groupSchedules, model.AccessControlGroupSchedule{
			AccessControlGroupID: groupID,
			DayOfWeek:            schedule.DayOfWeek,
			Date:                 datePtr,
			StartTime:            schedule.StartTime,
			EndTime:              schedule.EndTime,
		})
	}
	return groupSchedules, nil
}

// Validate and set default
func (s *accessControlGroupServiceImpl) validateAndSetDefaultValues(bodyRequest *schema.AccessControlGroupRequest) (*schema.AccessControlGroupRequest, error) {
	if bodyRequest.Name == nil || *bodyRequest.Name == "" {
		return nil, fmt.Errorf("group name cannot be empty")
	}
	if bodyRequest.AccessControlDeviceIDs == nil {
		bodyRequest.AccessControlDeviceIDs = []string{}
	}
	if bodyRequest.AccessControlGroupSchedules == nil {
		defaultSchedules := []schema.AccessControlGroupScheduleRequest{}

		// วนลูปสร้าง Schedule สำหรับ DayOfWeek 1 (จันทร์) ถึง 7 (อาทิตย์)
		for day := 1; day <= 7; day++ {
			defaultSchedules = append(defaultSchedules, schema.AccessControlGroupScheduleRequest{
				DayOfWeek: day,
				// ตลอดวัน (00:00:00 - 23:59:59)
				Date:      nil, // สมมติว่า Date ถูกใช้สำหรับวันเฉพาะ ไม่ใช่วันซ้ำ
				StartTime: "00:00:00",
				EndTime:   "23:59:59",
			})
		}
		bodyRequest.AccessControlGroupSchedules = defaultSchedules
	}
	return bodyRequest, nil
}

// validateBodyRequest checks for duplicate group name.
func (s *accessControlGroupServiceImpl) validateBodyRequest(bodyRequest schema.AccessControlGroupRequest, groupModel *model.AccessControlGroup) error {

	excludeID := uuid.Nil
	if groupModel != nil {
		excludeID = groupModel.ID
	}

	// Check duplicate name
	if bodyRequest.Name != nil {
		isExistName, err := s.accessControlGroupRepo.IsExistName(*bodyRequest.Name, excludeID)
		if err != nil {
			return fmt.Errorf("failed to check group name existence: %w", err)
		}
		if isExistName {
			// แต่เนื่องจากเป็น Service Layer ควร return เป็น error เพื่อให้ Handler จัดการ
			return fmt.Errorf("accessControlGroup name is already exist: %s", *bodyRequest.Name)
		}
	}

	// Note: การตรวจสอบว่า AccessControlDeviceIDs มีอยู่จริงหรือไม่ ถูกย้ายไปทำใน createGroupDeviceModels
	return nil
}
