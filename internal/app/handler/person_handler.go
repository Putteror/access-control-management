package handler

import (
	"mime/multipart"
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

var validate *validator.Validate

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

func (h *PersonHandler) FindAll(c *gin.Context) {

	var searchQuery schema.PersonSearchQuery
	if err := c.ShouldBindQuery(&searchQuery); err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, "Invalid search query parameter")
		return
	}
	if searchQuery.Page <= 0 {
		searchQuery.Page = 1
	}
	if searchQuery.Limit <= 0 {
		searchQuery.Limit = 10
	}

	persons, err := h.service.GetAll(searchQuery)
	if err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	personResponses := make([]schema.PersonResponse, len(persons))
	for i, person := range persons {
		personResponses[i] = schema.PersonResponse{
			ID:            person.ID,
			FirstName:     person.FirstName,
			MiddleName:    person.MiddleName,
			LastName:      person.LastName,
			PersonType:    person.PersonType,
			PersonID:      person.PersonID,
			Gender:        person.Gender,
			DateOfBirth:   person.DateOfBirth,
			Company:       person.Company,
			Department:    person.Department,
			JobPosition:   person.JobPosition,
			Address:       person.Address,
			MobileNumber:  person.MobileNumber,
			Email:         person.Email,
			FaceImagePath: person.FaceImagePath,
			ActiveAt:      person.ActiveAt,
			ExpireAt:      person.ExpireAt,
		}
	}

	pageData := common.PageResponse{
		Page:      searchQuery.Page,
		Size:      searchQuery.Limit,
		Total:     len(persons),
		TotalPage: (len(persons) + searchQuery.Limit - 1) / searchQuery.Limit,
	}

	common.GetDataListResponse(c, "Success", personResponses, pageData)
}

func (h *PersonHandler) FindByID(c *gin.Context) {
	id := c.Param("id")

	person, err := h.service.GetByID(id)
	if err != nil {
		common.ErrorResponse(c, http.StatusNotFound, "Person not found")
		return
	}

	common.SuccessResponse(c, "Success", person)
}

func (h *PersonHandler) Create(c *gin.Context) {
	var bodyRequest schema.PersonRequest
	if err := c.ShouldBind(&bodyRequest); err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, err.Error())
		common.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := validate.Struct(bodyRequest); err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var faceImageFile *multipart.FileHeader
	file, err := c.FormFile("faceImage")
	if err != nil {
		if err == http.ErrMissingFile {
			faceImageFile = nil
		} else {
			common.ErrorResponse(c, http.StatusBadRequest, "Failed to get image file: "+err.Error())
			return
		}
	} else {
		faceImageFile = file
	}

	var dob *time.Time
	if bodyRequest.DateOfBirth != "" {
		parsedDob, err := time.Parse("2006-01-02", bodyRequest.DateOfBirth)
		if err != nil {
			common.ErrorResponse(c, http.StatusBadRequest, "Invalid date of birth format")
			return
		}
		dob = &parsedDob
	}

	var activate *time.Time
	if bodyRequest.ActiveAt != "" {
		active_date, err := time.Parse("2006-01-02", bodyRequest.ActiveAt)
		if err != nil {
			common.ErrorResponse(c, http.StatusBadRequest, "Invalid active date format")
			return
		}
		activate = &active_date
	}

	var expire *time.Time
	if bodyRequest.ExpireAt != "" {
		expire_date, err := time.Parse("2006-01-02", bodyRequest.ExpireAt)
		if err != nil {
			common.ErrorResponse(c, http.StatusBadRequest, "Invalid expire date format")
			return
		}
		if expire_date.Before(*activate) {
			common.ErrorResponse(c, http.StatusBadRequest, "expire date must be after active date")
			return
		}
		expire = &expire_date
	}

	person := model.Person{
		FirstName:           bodyRequest.FirstName,
		MiddleName:          bodyRequest.MiddleName,
		LastName:            bodyRequest.LastName,
		PersonType:          bodyRequest.PersonType,
		PersonID:            bodyRequest.PersonID,
		Gender:              bodyRequest.Gender,
		DateOfBirth:         dob,
		Company:             bodyRequest.Company,
		Department:          bodyRequest.Department,
		JobPosition:         bodyRequest.JobPosition,
		Address:             bodyRequest.Address,
		MobileNumber:        bodyRequest.MobileNumber,
		Email:               bodyRequest.Email,
		IsVerified:          bodyRequest.IsVerified,
		ActiveAt:            activate,
		ExpireAt:            expire,
		AccessControlRuleID: bodyRequest.AccessControlRuleID,
		TimeAttendanceID:    bodyRequest.TimeAttendanceID,
	}

	if err := h.service.Save(nil, &person, faceImageFile); err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	common.SuccessResponse(c, "Create person success", person)
}

func (h *PersonHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.Delete(id); err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	common.SuccessResponse(c, "Person deleted successfully", nil)
}
