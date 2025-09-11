package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/putteror/access-control-management/internal/app/common"
	"github.com/putteror/access-control-management/internal/app/dto"
	"github.com/putteror/access-control-management/internal/app/model"
	"github.com/putteror/access-control-management/internal/app/service"

	"github.com/gin-gonic/gin"
)

var PERSON_TYPE_LIST = []string{"employee", "visitor"}

var validate *validator.Validate

func init() {
	validate = validator.New()
	validate.RegisterValidation("personType", validatePersonType)
}

type PersonRequest struct {
	FirstName        string `form:"firstName" validate:"required"`
	MiddleName       string `form:"middleName"`
	LastName         string `form:"lastName" validate:"required"`
	PersonType       string `form:"personType" validate:"required,personType"`
	PersonID         string `form:"personId"`
	Gender           string `form:"gender"`
	DateOfBirth      string `form:"dateOfBirth"`
	Company          string `form:"company"`
	Department       string `form:"department"`
	JobPosition      string `form:"jobPosition"`
	Address          string `form:"address"`
	MobileNumber     string `form:"mobileNumber"`
	Email            string `form:"email"`
	IsVerified       bool   `form:"isVerified"`
	ActiveAt         string `form:"activeAt"`
	ExpireAt         string `form:"expireAt"`
	RuleID           string `form:"ruleId"`
	TimeAttendanceID string `form:"timeAttendanceId"`
	// Face image will receive in function
}

// It checks if the value of the "personType" field is in our list of allowed values.
func validatePersonType(fl validator.FieldLevel) bool {
	personType := strings.ToLower(fl.Field().String())
	for _, allowedType := range PERSON_TYPE_LIST {
		if personType == allowedType {
			return true
		}
	}
	return false
}

type PersonHandler struct {
	service service.PersonService
}

func NewPersonHandler(service service.PersonService) *PersonHandler {
	return &PersonHandler{service: service}
}

func (h *PersonHandler) FindAll(c *gin.Context) {
	// Get query parameters for pagination using the dedicated function
	page, limit, err := common.GetPaginationParams(c)
	if err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, "Invalid page or limit parameter")
		return
	}

	// Retrieve persons from the service with pagination
	persons, err := h.service.PaginatedFindAllPeople(page, limit)
	if err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	personDTOs := make([]dto.PersonDTO, len(persons))
	for i, person := range persons {
		personDTOs[i] = dto.PersonDTO{
			ID:           person.ID,
			FirstName:    person.FirstName,
			LastName:     person.LastName,
			Email:        person.Email,
			DateOfBirth:  person.DateOfBirth,
			MobileNumber: person.MobileNumber,
		}
	}

	common.SuccessResponse(c, "All persons retrieved", personDTOs)
}

func (h *PersonHandler) FindByID(c *gin.Context) {
	id := c.Param("id")

	person, err := h.service.GetPersonByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Person not found"})
		return
	}

	c.JSON(http.StatusOK, person)
}

func (h *PersonHandler) Create(c *gin.Context) {
	var bodyRequest PersonRequest
	if err := c.ShouldBind(&bodyRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validate.Struct(bodyRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Face image upload
	faceImageFile, err := c.FormFile("faceImage")
	if err != nil {
		if err == http.ErrMissingFile {
			faceImageFile = nil
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get image file: " + err.Error()})
		return
	}

	dob, err := time.Parse("2006-01-02", bodyRequest.DateOfBirth)
	if err != nil {
		// หาก Format ของวันที่ไม่ถูกต้อง จะเกิด Error ที่นี่
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date of birth format"})
		return
	}

	active_date, err := time.Parse("2006-01-02", bodyRequest.ActiveAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid active date format"})
		return
	}

	expire_date, err := time.Parse("2006-01-02", bodyRequest.ExpireAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid expire date format"})
		return
	}

	person := model.Person{
		FirstName:        bodyRequest.FirstName,
		MiddleName:       bodyRequest.MiddleName,
		LastName:         bodyRequest.LastName,
		PersonType:       bodyRequest.PersonType,
		PersonID:         bodyRequest.PersonID,
		Gender:           bodyRequest.Gender,
		DateOfBirth:      &dob,
		Company:          bodyRequest.Company,
		Department:       bodyRequest.Department,
		JobPosition:      bodyRequest.JobPosition,
		Address:          bodyRequest.Address,
		MobileNumber:     bodyRequest.MobileNumber,
		Email:            bodyRequest.Email,
		IsVerified:       bodyRequest.IsVerified,
		ActiveAt:         active_date,
		ExpireAt:         expire_date,
		RuleID:           bodyRequest.RuleID,
		TimeAttendanceID: bodyRequest.TimeAttendanceID,
	}

	if err := h.service.CreatePerson(&person, faceImageFile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, person)
}

func (h *PersonHandler) DeletePerson(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.DeletePerson(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Person deleted successfully"})
}
