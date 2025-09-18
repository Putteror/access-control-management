package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/putteror/access-control-management/internal/app/common"
	"github.com/putteror/access-control-management/internal/app/schema"
	"github.com/putteror/access-control-management/internal/app/service"
)

// ตัวแปร validate ถูกกำหนดไว้แล้วในไฟล์อื่น (น่าจะอยู่ในไฟล์เดียวกันกับ AccessControlDeviceHandler หรือไฟล์ที่ import)
// หากต้องการให้โค้ดนี้รันได้โดยสมบูรณ์ อาจต้องเพิ่ม var validate *validator.Validate
// แต่เพื่อให้สอดคล้องกับไฟล์ต้นฉบับ จะขอเว้นไว้ และสมมติว่ามีการประกาศไว้แล้ว

type AccessControlServerHandler struct {
	service service.AccessControlServerService
}

func NewAccessControlServerHandler(service service.AccessControlServerService) *AccessControlServerHandler {
	return &AccessControlServerHandler{service: service}
}

// GetAll retrieves all access control servers.
func (h *AccessControlServerHandler) GetAll(c *gin.Context) {

	var searchQuery schema.AccessControlServerSearchQuery
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
	servers, err := h.service.GetAll(searchQuery)
	if err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	serverResponses := make([]schema.AccessControlServerResponse, len(servers))
	for i, server := range servers {
		response, err := h.service.ConvertToResponse(&server)
		if err != nil {
			common.ErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		serverResponses[i] = *response
	}

	pageData := common.PageResponse{
		Page:      searchQuery.Page,
		Size:      searchQuery.Limit,
		Total:     len(servers),
		TotalPage: (len(servers) + searchQuery.Limit - 1) / searchQuery.Limit,
	}

	common.GetDataListResponse(c, "Success", serverResponses, pageData)
}

// GetByID retrieves an access control server by its ID.
func (h *AccessControlServerHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	server, err := h.service.GetByID(id)
	if err != nil {
		common.ErrorResponse(c, http.StatusNotFound, "Server not found")
		return
	}

	serverResponse, err := h.service.ConvertToResponse(server)
	if err != nil {
		common.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	common.SuccessResponse(c, "Success", serverResponse)
}

// Create creates a new access control server.
func (h *AccessControlServerHandler) Create(c *gin.Context) {
	var bodyRequest schema.AccessControlServerRequest
	if err := c.ShouldBindJSON(&bodyRequest); err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// ต้องมั่นใจว่ามีการประกาศ var validate = validator.New() ไว้ในไฟล์เดียวกันหรือไฟล์ที่ import
	if err := validate.Struct(bodyRequest); err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	serverModel, err := h.service.Create(&bodyRequest)
	if err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	serverResponse, err := h.service.ConvertToResponse(serverModel)
	if err != nil {
		common.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	common.SuccessResponse(c, "Create server success", serverResponse)
}

// Update updates an existing access control server.
func (h *AccessControlServerHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var bodyRequest schema.AccessControlServerRequest
	if err := c.ShouldBind(&bodyRequest); err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := validate.Struct(bodyRequest); err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	serverModel, err := h.service.Update(id, &bodyRequest)
	if err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	serverResponse, err := h.service.ConvertToResponse(serverModel)
	if err != nil {
		common.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	common.SuccessResponse(c, "Update server success", serverResponse)
}

// PartialUpdate performs a partial update on an existing access control server.
func (h *AccessControlServerHandler) PartialUpdate(c *gin.Context) {
	id := c.Param("id")

	var bodyRequest schema.AccessControlServerRequest
	if err := c.ShouldBindJSON(&bodyRequest); err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	serverModel, err := h.service.PartialUpdate(id, &bodyRequest)
	if err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	serverResponse, err := h.service.ConvertToResponse(serverModel)
	if err != nil {
		common.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	common.SuccessResponse(c, "Update server success", serverResponse)
}

// Delete deletes an access control server by its ID.
func (h *AccessControlServerHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.Delete(id); err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	common.SuccessResponse(c, "Server deleted successfully", nil)
}
