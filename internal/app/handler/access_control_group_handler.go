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
// แต่เพื่อให้สอดคล้องกับไฟล์ต้นฉบับ จะขอเว้นไว้ และสมมติว่ามีการประกาศไว้แล้ว

type AccessControlGroupHandler struct {
	service service.AccessControlGroupService
}

func NewAccessControlGroupHandler(service service.AccessControlGroupService) *AccessControlGroupHandler {
	return &AccessControlGroupHandler{service: service}
}

func init() {
	validate = validator.New()
}

// ## Error Handling Utility
// ฟังก์ชันนี้ช่วยจัดการข้อผิดพลาดเฉพาะที่เกี่ยวกับชื่อซ้ำ

func handleErrorResponse(c *gin.Context, err error, defaultMessage string) {
	// ตรวจสอบข้อความ error สำหรับกรณีชื่อซ้ำที่กำหนดไว้ก่อนหน้า
	if strings.Contains(err.Error(), "accessControlGroup name is already exist") {
		// นำไปใช้กับ format ที่กำหนดไว้ [2025-07-05]
		message := "Group name is already exist"
		common.ErrorResponse(c, http.StatusBadRequest, message)
		return
	}
	// สำหรับข้อผิดพลาดอื่น ๆ
	common.ErrorResponse(c, http.StatusInternalServerError, defaultMessage)
}

// ---

// ## Get Operations

// GetAll retrieves all access control groups.
func (h *AccessControlGroupHandler) GetAll(c *gin.Context) {

	var searchQuery schema.AccessControlGroupSearchQuery
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
	groups, err := h.service.GetAll(searchQuery)
	if err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	groupResponses := make([]schema.AccessControlGroupResponse, len(groups))
	for i, group := range groups {
		response, err := h.service.ConvertToResponse(&group)
		if err != nil {
			common.ErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		groupResponses[i] = *response
	}

	pageData := common.PageResponse{
		Page:      searchQuery.Page,
		Size:      searchQuery.Limit,
		Total:     len(groups),
		TotalPage: (len(groups) + searchQuery.Limit - 1) / searchQuery.Limit,
	}

	common.GetDataListResponse(c, "Success", groupResponses, pageData)
}

// GetByID retrieves an access control group by its ID.
func (h *AccessControlGroupHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	group, err := h.service.GetByID(id)
	if err != nil {
		common.ErrorResponse(c, http.StatusNotFound, "Group not found")
		return
	}

	groupResponse, err := h.service.ConvertToResponse(group)
	if err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	common.SuccessResponse(c, "Success", groupResponse)
}

// ---

// ## CRUD Operations

// Create creates a new access control group.
func (h *AccessControlGroupHandler) Create(c *gin.Context) {
	var bodyRequest schema.AccessControlGroupRequest
	if err := c.ShouldBindJSON(&bodyRequest); err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := validate.Struct(bodyRequest); err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	groupModel, err := h.service.Create(&bodyRequest)
	if err != nil {
		// ใช้ handleErrorResponse เพื่อจัดการข้อผิดพลาดชื่อซ้ำ
		handleErrorResponse(c, err, err.Error())
		return
	}

	groupResponse, err := h.service.ConvertToResponse(groupModel)
	if err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	common.SuccessResponse(c, "Create group success", groupResponse)
}

// Update updates an existing access control group (Full Update).
func (h *AccessControlGroupHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var bodyRequest schema.AccessControlGroupRequest
	// ควรใช้ ShouldBindJSON แทน ShouldBind เผื่อการตรวจสอบ Content-Type ที่เข้มงวดกว่า
	if err := c.ShouldBindJSON(&bodyRequest); err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := validate.Struct(bodyRequest); err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	groupModel, err := h.service.Update(id, &bodyRequest)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			common.ErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		handleErrorResponse(c, err, err.Error())
		return
	}

	groupResponse, err := h.service.ConvertToResponse(groupModel)
	if err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	common.SuccessResponse(c, "Update group success", groupResponse)
}

// PartialUpdate updates an existing access control group (Partial Update).
func (h *AccessControlGroupHandler) PartialUpdate(c *gin.Context) {
	id := c.Param("id")

	var bodyRequest schema.AccessControlGroupRequest
	if err := c.ShouldBindJSON(&bodyRequest); err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	// ใน Partial Update มักจะไม่เรียก validate.Struct() เว้นแต่คุณจะใช้ validation เฉพาะบางฟิลด์

	groupModel, err := h.service.PartialUpdate(id, &bodyRequest)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			common.ErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		handleErrorResponse(c, err, err.Error())
		return
	}

	groupResponse, err := h.service.ConvertToResponse(groupModel)
	if err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	common.SuccessResponse(c, "Update group success", groupResponse)
}

// Delete deletes an access control group by its ID.
func (h *AccessControlGroupHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.Delete(id); err != nil {
		if strings.Contains(err.Error(), "not found") {
			common.ErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		common.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	common.SuccessResponse(c, "Group deleted successfully", nil)
}
