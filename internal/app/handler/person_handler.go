package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/putteror/access-control-management/internal/app/common"
	"github.com/putteror/access-control-management/internal/app/model"
	"github.com/putteror/access-control-management/internal/app/schema"
	"github.com/putteror/access-control-management/internal/app/service"

	"github.com/gin-gonic/gin"
)

type PersonHandler struct {
	service service.PersonService
}

func NewPersonHandler(service service.PersonService) *PersonHandler {
	return &PersonHandler{service: service}
}

func init() {
	validate = validator.New()
	validate.RegisterValidation("personType", validatePersonType)
}

func validatePersonType(fl validator.FieldLevel) bool {
	personType := strings.ToLower(fl.Field().String())
	for _, allowedType := range schema.PERSON_TYPE_LIST {
		if personType == allowedType {
			return true
		}
	}
	return false
}

func personHandleErrorResponse(c *gin.Context, err error, defaultMessage string) {
	if strings.Contains(err.Error(), "person name is already exist") {
		// นำไปใช้กับ format ที่กำหนดไว้ [2025-07-05]
		message := "ApplicationForm name is already exist"
		common.ErrorResponse(c, http.StatusBadRequest, message)
		return
	}
	common.ErrorResponse(c, http.StatusInternalServerError, defaultMessage)
}

// ---

// ## Get Operations

// GetAll retrieves all persons.
func (h *PersonHandler) GetAll(c *gin.Context) {
	var searchQuery schema.PersonSearchQuery
	if err := c.ShouldBindQuery(&searchQuery); err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, "Invalid search query parameter")
		return
	}

	if searchQuery.Page <= 0 {
		searchQuery.Page = common.DefaultPage
	}
	if searchQuery.Limit <= 0 {
		searchQuery.Limit = common.DefaultPageSize
	}

	persons, err := h.service.GetAll(searchQuery)
	if err != nil {
		personHandleErrorResponse(c, err, err.Error())
		return
	}

	personResponses := make([]schema.PersonResponse, len(persons))
	for i, person := range persons {
		response, err := h.service.ConvertToResponse(&person)
		if err != nil {
			common.ErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		personResponses[i] = *response
	}

	pageData := common.PageResponse{
		Page:      searchQuery.Page,
		Size:      searchQuery.Limit,
		Total:     len(persons),
		TotalPage: (len(persons) + searchQuery.Limit - 1) / searchQuery.Limit,
	}

	common.GetDataListResponse(c, "Success", personResponses, pageData)
}

// GetByID retrieves a person by its ID.
func (h *PersonHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	person, err := h.service.GetByID(id)
	if err != nil {
		common.ErrorResponse(c, http.StatusNotFound, "Person not found")
		return
	}

	personResponse, err := h.service.ConvertToResponse(person)
	if err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	common.SuccessResponse(c, "Success", personResponse)
}

// ---

// ## CRUD Operations

// Create creates a new person.
func (h *PersonHandler) Create(c *gin.Context) {
	var bodyRequest schema.PersonRequest
	if err := c.ShouldBind(&bodyRequest); err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := validate.Struct(bodyRequest); err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	faceImageFile, err := c.FormFile("faceImage")
	if err != nil && err != http.ErrMissingFile {
		common.ErrorResponse(c, http.StatusBadRequest, "Failed to get image file: "+err.Error())
		return
	}

	dob, err := parseTimeFromRequest(bodyRequest.DateOfBirth)
	if err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, "Invalid date of birth format")
		return
	}

	activate, err := parseTimeFromRequest(bodyRequest.ActiveAt)
	if err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, "Invalid active date format")
		return
	}

	expire, err := parseTimeFromRequest(bodyRequest.ExpireAt)
	if err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, "Invalid expire date format")
		return
	}
	if expire != nil && activate != nil && expire.Before(*activate) {
		common.ErrorResponse(c, http.StatusBadRequest, "expire date must be after active date")
		return
	}

	person := convertToModel(&bodyRequest, dob, activate, expire)

	if err := h.service.Save("", person, faceImageFile, bodyRequest.CardIDs, bodyRequest.LicensePlateTexts); err != nil {
		personHandleErrorResponse(c, err, err.Error())
		return
	}

	personResponse, err := h.service.ConvertToResponse(person)
	if err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	common.SuccessResponse(c, "Create person success", personResponse)
}

// Update updates an existing person.
func (h *PersonHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var bodyRequest schema.PersonRequest
	if err := c.ShouldBind(&bodyRequest); err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := validate.Struct(bodyRequest); err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	faceImageFile, err := c.FormFile("faceImage")
	if err != nil && err != http.ErrMissingFile {
		common.ErrorResponse(c, http.StatusBadRequest, "Failed to get image file: "+err.Error())
		return
	}

	dob, err := parseTimeFromRequest(bodyRequest.DateOfBirth)
	if err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, "Invalid date of birth format")
		return
	}

	activate, err := parseTimeFromRequest(bodyRequest.ActiveAt)
	if err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, "Invalid active date format")
		return
	}

	expire, err := parseTimeFromRequest(bodyRequest.ExpireAt)
	if err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, "Invalid expire date format")
		return
	}
	if expire != nil && activate != nil && expire.Before(*activate) {
		common.ErrorResponse(c, http.StatusBadRequest, "expire date must be after active date")
		return
	}

	person := convertToModel(&bodyRequest, dob, activate, expire)

	if err := h.service.Save(id, person, faceImageFile, bodyRequest.CardIDs, bodyRequest.LicensePlateTexts); err != nil {
		if strings.Contains(err.Error(), "not found") {
			common.ErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		personHandleErrorResponse(c, err, err.Error())
		return
	}

	personResponse, err := h.service.ConvertToResponse(person)
	if err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	common.SuccessResponse(c, "Update person success", personResponse)
}

// PartialUpdate updates an existing person.
func (h *PersonHandler) PartialUpdate(c *gin.Context) {
	id := c.Param("id")

	var bodyRequest schema.PersonRequest
	if err := c.ShouldBind(&bodyRequest); err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	faceImageFile, err := c.FormFile("faceImage")
	if err != nil && err != http.ErrMissingFile {
		common.ErrorResponse(c, http.StatusBadRequest, "Failed to get image file: "+err.Error())
		return
	}

	dob, err := parseTimeFromRequest(bodyRequest.DateOfBirth)
	if err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, "Invalid date of birth format")
		return
	}

	activate, err := parseTimeFromRequest(bodyRequest.ActiveAt)
	if err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, "Invalid active date format")
		return
	}

	expire, err := parseTimeFromRequest(bodyRequest.ExpireAt)
	if err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, "Invalid expire date format")
		return
	}

	person := convertToModel(&bodyRequest, dob, activate, expire)

	if err := h.service.PartialUpdate(id, person, faceImageFile, bodyRequest.CardIDs, bodyRequest.LicensePlateTexts); err != nil {
		if strings.Contains(err.Error(), "not found") {
			common.ErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		personHandleErrorResponse(c, err, err.Error())
		return
	}

	personResponse, err := h.service.ConvertToResponse(person)
	if err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	common.SuccessResponse(c, "Update person success", personResponse)
}

// Delete deletes a person by its ID.
func (h *PersonHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.Delete(id); err != nil {
		if strings.Contains(err.Error(), "not found") {
			common.ErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		common.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	common.SuccessResponse(c, "Person deleted successfully", nil)
}

// ---

// ## Helper Functions

// parseTimeFromRequest safely parses a string into a time.Time pointer.
func parseTimeFromRequest(timeString *string) (*time.Time, error) {
	if timeString == nil || *timeString == "" {
		return nil, nil
	}
	parsedTime, err := time.Parse("2006-01-02", *timeString)
	if err != nil {
		return nil, err
	}
	return &parsedTime, nil
}

// convertToModel converts a schema.PersonRequest to a model.Person.
func convertToModel(bodyRequest *schema.PersonRequest, dob, activate, expire *time.Time) *model.Person {
	return &model.Person{
		FirstName:           *bodyRequest.FirstName,
		MiddleName:          bodyRequest.MiddleName,
		LastName:            *bodyRequest.LastName,
		PersonType:          *bodyRequest.PersonType,
		PersonID:            bodyRequest.PersonID,
		Gender:              bodyRequest.Gender,
		DateOfBirth:         dob,
		Company:             bodyRequest.Company,
		Department:          bodyRequest.Department,
		JobPosition:         bodyRequest.JobPosition,
		Address:             bodyRequest.Address,
		MobileNumber:        bodyRequest.MobileNumber,
		Email:               bodyRequest.Email,
		ActiveAt:            activate,
		ExpireAt:            expire,
		AccessControlRuleID: bodyRequest.AccessControlRuleID,
		TimeAttendanceID:    bodyRequest.TimeAttendanceID,
	}
}
