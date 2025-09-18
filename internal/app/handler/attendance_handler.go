package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/putteror/access-control-management/internal/app/common"
	"github.com/putteror/access-control-management/internal/app/schema"
	"github.com/putteror/access-control-management/internal/app/service"
)

// ต้องมั่นใจว่ามีการประกาศ var validate *validator.Validate และมีการเรียก init()
// เราจะสมมติว่ามีการประกาศและเรียกใช้ในไฟล์หลักของ package handler แล้ว
// var validate *validator.Validate
// func init() { validate = validator.New() }

type AttendanceHandler struct {
	service service.AttendanceService
}

func NewAttendanceHandler(service service.AttendanceService) *AttendanceHandler {
	return &AttendanceHandler{service: service}
}

// ## Error Handling Utility
// ฟังก์ชันนี้ช่วยจัดการข้อผิดพลาดเฉพาะที่เกี่ยวกับชื่อซ้ำสำหรับ Attendance
func handleAttendanceErrorResponse(c *gin.Context, err error, defaultMessage string) {
	// ตรวจสอบข้อความ error สำหรับกรณีชื่อซ้ำที่กำหนดไว้ก่อนหน้า
	if strings.Contains(err.Error(), "ApplicationForm name is already exist") {
		// นำไปใช้กับ format ที่กำหนดไว้ [2025-07-05]
		message := "Attendance name is already exist"
		common.ErrorResponse(c, http.StatusBadRequest, message)
		return
	}
	// สำหรับข้อผิดพลาดอื่น ๆ ที่มาจาก Service
	common.ErrorResponse(c, http.StatusInternalServerError, defaultMessage)
}

// ---

// ## Get Operations

// GetAll retrieves all attendance records.
func (h *AttendanceHandler) GetAll(c *gin.Context) {

	var searchQuery schema.AttendanceSearchQuery
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
	attendances, err := h.service.GetAll(searchQuery)
	if err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	attendanceResponses := make([]schema.AttendanceInfoResponse, 0, len(attendances))
	for _, attendance := range attendances {
		response, err := h.service.ConvertToResponse(&attendance)
		if err != nil {
			common.ErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		attendanceResponses = append(attendanceResponses, *response)
	}

	// Note: การคำนวณ TotalPage ด้วย len(attendances) อาจไม่ถูกต้องหากมีการใช้ pagination จริง
	// แต่จะคงไว้ตาม logic เดิมเพื่อ consistency
	pageData := common.PageResponse{
		Page:      searchQuery.Page,
		Size:      searchQuery.Limit,
		Total:     len(attendances),
		TotalPage: (len(attendances) + searchQuery.Limit - 1) / searchQuery.Limit,
	}

	common.GetDataListResponse(c, "Success", attendanceResponses, pageData)
}

// GetByID retrieves an attendance record by its ID.
func (h *AttendanceHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	attendance, err := h.service.GetByID(id)
	if err != nil {
		common.ErrorResponse(c, http.StatusNotFound, "Attendance record not found")
		return
	}

	attendanceResponse, err := h.service.ConvertToResponse(attendance)
	if err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	common.SuccessResponse(c, "Success", attendanceResponse)
}

// ---

// ## CRUD Operations

// Create creates a new attendance record.
func (h *AttendanceHandler) Create(c *gin.Context) {
	var bodyRequest schema.AttendanceRequest
	if err := c.ShouldBindJSON(&bodyRequest); err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// ตรวจสอบโครงสร้าง
	if validate == nil {
		validate = validator.New()
	}
	if err := validate.Struct(bodyRequest); err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	attendanceModel, err := h.service.Create(&bodyRequest)
	if err != nil {
		// ใช้ handleAttendanceErrorResponse เพื่อจัดการข้อผิดพลาดชื่อซ้ำ
		handleAttendanceErrorResponse(c, err, err.Error())
		return
	}

	attendanceResponse, err := h.service.ConvertToResponse(attendanceModel)
	if err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	common.SuccessResponse(c, "Create attendance success", attendanceResponse)
}

// Update updates an existing attendance record (Full Update).
func (h *AttendanceHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var bodyRequest schema.AttendanceRequest
	if err := c.ShouldBindJSON(&bodyRequest); err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// ตรวจสอบโครงสร้าง
	if validate == nil {
		validate = validator.New()
	}
	if err := validate.Struct(bodyRequest); err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	attendanceModel, err := h.service.Update(id, &bodyRequest)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			common.ErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		handleAttendanceErrorResponse(c, err, err.Error())
		return
	}

	attendanceResponse, err := h.service.ConvertToResponse(attendanceModel)
	if err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	common.SuccessResponse(c, "Update attendance success", attendanceResponse)
}

// PartialUpdate updates an existing attendance record (Partial Update).
func (h *AttendanceHandler) PartialUpdate(c *gin.Context) {
	id := c.Param("id")

	var bodyRequest schema.AttendanceRequest
	if err := c.ShouldBindJSON(&bodyRequest); err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	// ใน Partial Update มักจะไม่เรียก validate.Struct() เต็มรูปแบบ

	attendanceModel, err := h.service.PartialUpdate(id, &bodyRequest)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			common.ErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		handleAttendanceErrorResponse(c, err, err.Error())
		return
	}

	attendanceResponse, err := h.service.ConvertToResponse(attendanceModel)
	if err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	common.SuccessResponse(c, "Update attendance success", attendanceResponse)
}

// Delete deletes an attendance record by its ID.
func (h *AttendanceHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.Delete(id); err != nil {
		if strings.Contains(err.Error(), "not found") {
			common.ErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		common.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	common.SuccessResponse(c, "Attendance record deleted successfully", nil)
}
