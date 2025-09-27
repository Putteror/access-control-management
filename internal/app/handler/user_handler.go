package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/putteror/access-control-management/internal/app/common"
	"github.com/putteror/access-control-management/internal/app/schema"
	"github.com/putteror/access-control-management/internal/app/service"
)

// ต้องมั่นใจว่ามีการประกาศ var validate *validator.Validate และมีการเรียก init()
// เราจะสมมติว่ามีการประกาศและเรียกใช้ในไฟล์หลักของ package handler แล้ว
// var validate *validator.Validate
// func init() { validate = validator.New() }

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// ## Error Handling Utility
// ฟังก์ชันนี้ช่วยจัดการข้อผิดพลาดเฉพาะที่เกี่ยวกับชื่อซ้ำสำหรับ User
func handleUserErrorResponse(c *gin.Context, err error, defaultMessage string) {
	// ตรวจสอบข้อความ error สำหรับกรณีชื่อซ้ำที่กำหนดไว้ก่อนหน้า [2025-07-05]
	if strings.Contains(err.Error(), "ApplicationForm name is already exist") {
		// นำไปใช้กับ format ที่กำหนดไว้
		message := "User name is already exist" // เปลี่ยนให้สื่อถึง User
		common.ErrorResponse(c, http.StatusBadRequest, message)
		return
	}
	// สำหรับข้อผิดพลาดอื่น ๆ ที่มาจาก Service
	common.ErrorResponse(c, http.StatusInternalServerError, defaultMessage)
}

// ---

// ## Get Operations

// GetAll retrieves all user records.
func (h *UserHandler) GetAll(c *gin.Context) {

	var searchQuery schema.UserSearchQuery
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

	// UserService.GetAll คืนค่าเป็น []model.User
	userModels, err := h.service.GetAll(searchQuery)
	if err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	userResponses := make([]schema.UserResponse, 0, len(userModels))
	for _, userModel := range userModels {
		// ต้องแปลงเป็น Response Schema
		response, err := h.service.ConvertToResponse(&userModel)
		if err != nil {
			common.ErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		userResponses = append(userResponses, *response)
	}

	// Note: การคำนวณ TotalPage ด้วย len(userModels) อาจไม่ถูกต้องหากมีการใช้ pagination จริง
	// แต่จะคงไว้ตาม logic เดิมเพื่อ consistency
	pageData := common.PageResponse{
		Page:      searchQuery.Page,
		Size:      searchQuery.Limit,
		Total:     len(userModels),
		TotalPage: (len(userModels) + searchQuery.Limit - 1) / searchQuery.Limit,
	}

	common.GetDataListResponse(c, "Success", userResponses, pageData)
}

// GetByID retrieves a user record by its ID.
func (h *UserHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	// UserService.GetByID คืนค่าเป็น *model.User
	userModel, err := h.service.GetByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			common.ErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		common.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	userResponse, err := h.service.ConvertToResponse(userModel)
	if err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	common.SuccessResponse(c, "Success", userResponse)
}

// ---

// ## CRUD Operations

// Create creates a new user record.
func (h *UserHandler) Create(c *gin.Context) {
	var bodyRequest schema.UserRequest
	if err := c.ShouldBindJSON(&bodyRequest); err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// UserService.Create คืนค่าเป็น *model.User
	userModel, err := h.service.Create(&bodyRequest)
	if err != nil {
		// ใช้ handleUserErrorResponse เพื่อจัดการข้อผิดพลาดชื่อซ้ำ
		handleUserErrorResponse(c, err, err.Error())
		return
	}

	userResponse, err := h.service.ConvertToResponse(userModel)
	if err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	common.SuccessResponse(c, "Create user success", userResponse)
}

// Update updates an existing user record (Full Update).
func (h *UserHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var bodyRequest schema.UserRequest
	if err := c.ShouldBindJSON(&bodyRequest); err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// UserService.Update คืนค่าเป็น *model.User
	userModel, err := h.service.Update(id, &bodyRequest)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			common.ErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		handleUserErrorResponse(c, err, err.Error())
		return
	}

	userResponse, err := h.service.ConvertToResponse(userModel)
	if err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	common.SuccessResponse(c, "Update user success", userResponse)
}

// PartialUpdate updates an existing user record (Partial Update).
func (h *UserHandler) PartialUpdate(c *gin.Context) {
	id := c.Param("id")

	var bodyRequest schema.UserRequest
	if err := c.ShouldBindJSON(&bodyRequest); err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// UserService.PartialUpdate คืนค่าเป็น *model.User
	userModel, err := h.service.PartialUpdate(id, &bodyRequest)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			common.ErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		handleUserErrorResponse(c, err, err.Error())
		return
	}

	userResponse, err := h.service.ConvertToResponse(userModel)
	if err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	common.SuccessResponse(c, "Update user success", userResponse)
}

// Delete deletes a user record by its ID.
func (h *UserHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.Delete(id); err != nil {
		if strings.Contains(err.Error(), "not found") {
			common.ErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		common.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	common.SuccessResponse(c, "User record deleted successfully", nil)
}
