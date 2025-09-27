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

type AccessRecordService interface {
	GetAll(searchQuery schema.AccessRecordSearchQuery) ([]model.AccessRecord, error)
	GetByID(id string) (*model.AccessRecord, error)
	Create(bodyRequest *schema.AccessRecordRequest) (*model.AccessRecord, error)
	Update(id string, bodyRequest *schema.AccessRecordRequest) (*model.AccessRecord, error)
	PartialUpdate(id string, bodyRequest *schema.AccessRecordRequest) (*model.AccessRecord, error)
	Delete(id string) error
	ConvertToResponse(accessRecordModel *model.AccessRecord) (*schema.AccessRecordResponse, error)
}

type AccessRecordServiceImpl struct {
	accessRecordRepo repository.AccessRecordRepository
	personRepo       repository.PersonRepository
	deviceRepo       repository.AccessControlDeviceRepository
}

func NewAccessRecordService(accessRecordRepo repository.AccessRecordRepository, personRepo repository.PersonRepository, deviceRepo repository.AccessControlDeviceRepository) AccessRecordService {
	return &AccessRecordServiceImpl{
		accessRecordRepo: accessRecordRepo,
		personRepo:       personRepo,
		deviceRepo:       deviceRepo,
	}
}

func (s *AccessRecordServiceImpl) GetAll(searchQuery schema.AccessRecordSearchQuery) ([]model.AccessRecord, error) {
	return s.accessRecordRepo.GetAll(searchQuery)
}

func (s *AccessRecordServiceImpl) GetByID(id string) (*model.AccessRecord, error) {
	idUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID")
	}
	return s.accessRecordRepo.GetByID(idUUID)
}

func (s *AccessRecordServiceImpl) Create(bodyRequest *schema.AccessRecordRequest) (*model.AccessRecord, error) {

	bodyRequest, err := s.validateAndSetDefaultValues(bodyRequest)
	if err != nil {
		return nil, err
	}

	access_time, err := common.ConvertTimeStrToTime(*bodyRequest.AccessTime)
	if err != nil {
		return nil, err
	}

	// convert Time
	accessRecordModel := &model.AccessRecord{
		AccessControlDeviceID: bodyRequest.AccessControlDeviceID,
		PersonID:              bodyRequest.PersonID,
		Type:                  *bodyRequest.Type,
		Result:                *bodyRequest.Result,
		AccessTime:            access_time,
	}
	err = s.accessRecordRepo.Create(accessRecordModel)
	if err != nil {
		return nil, err
	}
	return accessRecordModel, nil
}

func (s *AccessRecordServiceImpl) Update(id string, bodyRequest *schema.AccessRecordRequest) (*model.AccessRecord, error) {

	id_uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID")
	}
	accessRecordModel, err := s.accessRecordRepo.GetByID(id_uuid)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing access record: %w", err)
	}
	if accessRecordModel == nil {
		return nil, fmt.Errorf("access record with ID '%s' not found", id)
	}

	// Set default value
	bodyRequest, err = s.validateAndSetDefaultValues(bodyRequest)
	if err != nil {
		return nil, err
	}
	// Validate
	if err := s.validateBodyRequest(*bodyRequest, accessRecordModel); err != nil {
		return nil, err
	}
	access_time, err := common.ConvertTimeStrToTime(*bodyRequest.AccessTime)
	if err != nil {
		return nil, err
	}

	// Update model
	accessRecordModel.AccessControlDeviceID = bodyRequest.AccessControlDeviceID
	accessRecordModel.PersonID = bodyRequest.PersonID
	accessRecordModel.Type = *bodyRequest.Type
	accessRecordModel.Result = *bodyRequest.Result
	accessRecordModel.AccessTime = access_time

	err = s.accessRecordRepo.Update(accessRecordModel)
	if err != nil {
		return nil, err
	}

	return accessRecordModel, nil
}

func (s *AccessRecordServiceImpl) PartialUpdate(id string, bodyRequest *schema.AccessRecordRequest) (*model.AccessRecord, error) {
	id_uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID")
	}
	accessRecordModel, err := s.accessRecordRepo.GetByID(id_uuid)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing access record: %w", err)
	}
	if accessRecordModel == nil {
		return nil, fmt.Errorf("access record with ID '%s' not found", id)
	}

	if err := s.validateBodyRequest(*bodyRequest, accessRecordModel); err != nil {
		return nil, err
	}
	if bodyRequest.AccessControlDeviceID != nil {
		accessRecordModel.AccessControlDeviceID = bodyRequest.AccessControlDeviceID
	}
	if bodyRequest.PersonID != nil {
		accessRecordModel.PersonID = bodyRequest.PersonID
	}
	if bodyRequest.Type != nil {
		accessRecordModel.Type = *bodyRequest.Type
	}
	if bodyRequest.Result != nil {
		accessRecordModel.Result = *bodyRequest.Result
	}
	if bodyRequest.AccessTime != nil {
		access_time, err := common.ConvertTimeStrToTime(*bodyRequest.AccessTime)
		if err != nil {
			return nil, err
		}
		accessRecordModel.AccessTime = access_time
	}

	err = s.accessRecordRepo.Update(accessRecordModel)
	if err != nil {
		return nil, err
	}

	return accessRecordModel, nil

}

func (s *AccessRecordServiceImpl) Delete(id string) error {
	id_uuid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid ID")
	}
	_, err = s.accessRecordRepo.GetByID(id_uuid)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("access record with ID '%s' not found", id)
		}
		return fmt.Errorf("failed to get access record by ID: %w", err)
	}
	return s.accessRecordRepo.Delete(id_uuid)
}

// ConvertToResponse converts a group model to a response schema.
func (s *AccessRecordServiceImpl) ConvertToResponse(accessRecordModel *model.AccessRecord) (*schema.AccessRecordResponse, error) {

	var personResponse *schema.AccessRecordPersonResponse
	var deviceResponse *schema.AccessRecordDeviceResponse

	if accessRecordModel.PersonID != nil && *accessRecordModel.PersonID != "" {

		person_uuid, err := uuid.Parse(*accessRecordModel.PersonID)
		if err != nil {
			return nil, err
		}
		personModel, err := s.personRepo.GetByID(person_uuid)
		if err != nil {
			return nil, err
		}

		personResponse = &schema.AccessRecordPersonResponse{
			ID:          personModel.ID.String(),
			FirstName:   personModel.FirstName,
			LastName:    personModel.LastName,
			Company:     personModel.Company,
			Department:  personModel.Department,
			JobPosition: personModel.JobPosition,
		}
	}

	if accessRecordModel.AccessControlDeviceID != nil {
		device_uuid, err := uuid.Parse(*accessRecordModel.AccessControlDeviceID)
		if err != nil {
			return nil, err
		}
		deviceModel, err := s.deviceRepo.GetByID(device_uuid)
		if err != nil {
			return nil, err
		}

		deviceResponse = &schema.AccessRecordDeviceResponse{
			ID:          deviceModel.ID.String(),
			Name:        deviceModel.Name,
			HostAddress: deviceModel.HostAddress,
			Type:        deviceModel.Type,
		}
	}

	response := &schema.AccessRecordResponse{
		ID:                  accessRecordModel.ID.String(),
		AccessControlDevice: deviceResponse,
		Person:              personResponse,
		Type:                accessRecordModel.Type,
		Result:              accessRecordModel.Result,
		AccessTime:          accessRecordModel.AccessTime.Format("2006-01-02 15:04:05"),
	}
	return response, nil
}

// ----------> INNER FUNCTION <-----------------------//

// Validate for pass whole requestBody to model
func (s *AccessRecordServiceImpl) validateAndSetDefaultValues(bodyRequest *schema.AccessRecordRequest) (*schema.AccessRecordRequest, error) {

	if bodyRequest.Type == nil || *bodyRequest.Type == "" {
		return nil, fmt.Errorf("type cannot be empty")
	}

	if bodyRequest.Result == nil || *bodyRequest.Result == "" {
		return nil, fmt.Errorf("result cannot be empty")
	}

	if bodyRequest.AccessTime == nil || *bodyRequest.AccessTime == "" {
		return nil, fmt.Errorf("access time cannot be empty")
	}

	return bodyRequest, nil
}

func (s *AccessRecordServiceImpl) validateBodyRequest(bodyRequest schema.AccessRecordRequest, accessRecordModel *model.AccessRecord) error {

	if bodyRequest.Type != nil && *bodyRequest.Type != "" {
		if !common.ValidateAccessRecordType(*bodyRequest.Type) {
			return fmt.Errorf("type must be 'in' or 'out'")
		}
	}

	if bodyRequest.Result != nil && *bodyRequest.Result != "" {
		if !common.ValidateAccessRecordResult(*bodyRequest.Result) {
			return fmt.Errorf("result must be 'success' or 'failed' or 'unknown'")
		}
	}

	return nil
}
