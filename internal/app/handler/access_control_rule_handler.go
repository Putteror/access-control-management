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

type AccessControlRuleHandler struct {
	service service.AccessControlRuleService
}

func NewAccessControlRuleHandler(service service.AccessControlRuleService) *AccessControlRuleHandler {
	// ต้องมั่นใจว่า service ที่ถูก inject เป็น service.AccessControlRuleService
	return &AccessControlRuleHandler{service: service}
}

// ## Error Handling Utility
// ฟังก์ชันนี้ช่วยจัดการข้อผิดพลาดเฉพาะที่เกี่ยวกับชื่อซ้ำสำหรับ Rule

func handleRuleErrorResponse(c *gin.Context, err error, defaultMessage string) {
	// ตรวจสอบข้อความ error สำหรับกรณีชื่อซ้ำที่กำหนดไว้ก่อนหน้า
	if strings.Contains(err.Error(), "accessControlRule name is already exist") {
		// นำไปใช้กับ format ที่กำหนดไว้ [2025-07-05]
		message := "Rule name is already exist" // ใช้ข้อความเดิมตามที่ user request ให้จำ
		common.ErrorResponse(c, http.StatusBadRequest, message)
		return
	}
	// สำหรับข้อผิดพลาดอื่น ๆ ที่มาจาก Service
	common.ErrorResponse(c, http.StatusInternalServerError, defaultMessage)
}

// ---

// ## Get Operations

// GetAll retrieves all access control rules.
func (h *AccessControlRuleHandler) GetAll(c *gin.Context) {

	var searchQuery schema.AccessControlRuleSearchQuery
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
	rules, err := h.service.GetAll(searchQuery)
	if err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	ruleResponses := make([]schema.AccessControlRuleResponse, 0, len(rules))
	for _, rule := range rules {
		response, err := h.service.ConvertToResponse(&rule)
		if err != nil {
			common.ErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		ruleResponses = append(ruleResponses, *response)
	}

	// Note: การคำนวณ TotalPage ด้วย len(rules) อาจไม่ถูกต้องหากมีการใช้ pagination จริงใน Repo
	// แต่จะคงไว้ตาม logic เดิมเพื่อ consistency
	pageData := common.PageResponse{
		Page:      searchQuery.Page,
		Size:      searchQuery.Limit,
		Total:     len(rules),
		TotalPage: (len(rules) + searchQuery.Limit - 1) / searchQuery.Limit,
	}

	common.GetDataListResponse(c, "Success", ruleResponses, pageData)
}

// GetByID retrieves an access control rule by its ID.
func (h *AccessControlRuleHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	rule, err := h.service.GetByID(id)
	if err != nil {
		common.ErrorResponse(c, http.StatusNotFound, "Rule not found")
		return
	}

	ruleResponse, err := h.service.ConvertToResponse(rule)
	if err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	common.SuccessResponse(c, "Success", ruleResponse)
}

// ---

// ## CRUD Operations

// Create creates a new access control rule.
func (h *AccessControlRuleHandler) Create(c *gin.Context) {
	var bodyRequest schema.AccessControlRuleRequest
	if err := c.ShouldBindJSON(&bodyRequest); err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// ใช้ validate.Struct() เพื่อตรวจสอบ Name (required)
	if validate == nil {
		validate = validator.New()
	}
	if err := validate.Struct(bodyRequest); err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	ruleModel, err := h.service.Create(&bodyRequest)
	if err != nil {
		// ใช้ handleRuleErrorResponse เพื่อจัดการข้อผิดพลาดชื่อซ้ำ
		handleRuleErrorResponse(c, err, err.Error())
		return
	}

	ruleResponse, err := h.service.ConvertToResponse(ruleModel)
	if err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	common.SuccessResponse(c, "Create rule success", ruleResponse)
}

// Update updates an existing access control rule (Full Update).
func (h *AccessControlRuleHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var bodyRequest schema.AccessControlRuleRequest
	if err := c.ShouldBindJSON(&bodyRequest); err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// ใช้ validate.Struct() เพื่อตรวจสอบ Name (required)
	if validate == nil {
		validate = validator.New()
	}
	if err := validate.Struct(bodyRequest); err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	ruleModel, err := h.service.Update(id, &bodyRequest)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			common.ErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		handleRuleErrorResponse(c, err, err.Error())
		return
	}

	ruleResponse, err := h.service.ConvertToResponse(ruleModel)
	if err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	common.SuccessResponse(c, "Update rule success", ruleResponse)
}

// PartialUpdate updates an existing access control rule (Partial Update).
func (h *AccessControlRuleHandler) PartialUpdate(c *gin.Context) {
	id := c.Param("id")

	var bodyRequest schema.AccessControlRuleRequest
	if err := c.ShouldBindJSON(&bodyRequest); err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	// ใน Partial Update เราจะไม่เรียก validate.Struct() เต็มรูปแบบ เพราะ Name อาจไม่ถูกส่งมา
	// แต่ Name ใน schema ถูกตั้งเป็น required ทำให้เกิดความย้อนแย้งเล็กน้อย
	// เราจะสมมติว่าถ้า Name ไม่ได้ถูกส่งมาใน JSON จะเป็น string ว่าง และ Service จะจัดการเอง

	ruleModel, err := h.service.PartialUpdate(id, &bodyRequest)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			common.ErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		handleRuleErrorResponse(c, err, err.Error())
		return
	}

	ruleResponse, err := h.service.ConvertToResponse(ruleModel)
	if err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	common.SuccessResponse(c, "Update rule success", ruleResponse)
}

// Delete deletes an access control rule by its ID.
func (h *AccessControlRuleHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.Delete(id); err != nil {
		if strings.Contains(err.Error(), "not found") {
			common.ErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		common.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	common.SuccessResponse(c, "Rule deleted successfully", nil)
}
