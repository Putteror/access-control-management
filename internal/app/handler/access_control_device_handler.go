package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/putteror/access-control-management/internal/app/common"
	"github.com/putteror/access-control-management/internal/app/model"
	"github.com/putteror/access-control-management/internal/app/schema"
	"github.com/putteror/access-control-management/internal/app/service"
)

type AccessControlDeviceHandler struct {
	service service.AccessControlDeviceService
}

func NewAccessControlDeviceHandler(service service.AccessControlDeviceService) *AccessControlDeviceHandler {
	return &AccessControlDeviceHandler{service: service}
}

func init() {
	validate = validator.New()
}

// GetAll retrieves all access control devices.
func (h *AccessControlDeviceHandler) GetAll(c *gin.Context) {

	var searchQuery schema.AccessControlDeviceSearchQuery
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
	devices, err := h.service.GetAll(searchQuery)
	if err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	deviceResponses := make([]schema.AccessControlDeviceResponse, len(devices))
	for i, device := range devices {
		response, err := h.service.ConvertToResponse(&device)
		if err != nil {
			common.ErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		deviceResponses[i] = *response
	}

	pageData := common.PageResponse{
		Page:      searchQuery.Page,
		Size:      searchQuery.Limit,
		Total:     len(devices),
		TotalPage: (len(devices) + searchQuery.Limit - 1) / searchQuery.Limit,
	}

	common.GetDataListResponse(c, "Success", deviceResponses, pageData)
}

// GetByID retrieves an access control device by its ID.
func (h *AccessControlDeviceHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	device, err := h.service.GetByID(id)
	if err != nil {
		common.ErrorResponse(c, http.StatusNotFound, "Device not found")
		return
	}

	deviceResponse, err := h.service.ConvertToResponse(device)
	if err != nil {
		common.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	common.SuccessResponse(c, "Success", deviceResponse)
}

// Create creates a new access control device.
func (h *AccessControlDeviceHandler) Create(c *gin.Context) {
	var bodyRequest schema.AccessControlDeviceRequest
	if err := c.ShouldBindJSON(&bodyRequest); err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := validate.Struct(bodyRequest); err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	device := model.AccessControlDevice{
		Name:                  bodyRequest.Name,
		Type:                  bodyRequest.Type,
		HostAddress:           bodyRequest.HostAddress,
		Username:              bodyRequest.Username,
		Password:              bodyRequest.Password,
		AccessToken:           bodyRequest.AccessToken,
		ApiToken:              bodyRequest.ApiToken,
		RecordScan:            bodyRequest.RecordScan,
		RecordAttendance:      bodyRequest.RecordAttendance,
		AllowClockIn:          bodyRequest.AllowClockIn,
		AllowClockOut:         bodyRequest.AllowClockOut,
		Status:                bodyRequest.Status,
		AccessControlServerID: bodyRequest.AccessControlServerID,
	}

	if err := h.service.Save("", &device); err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	deviceResponse, err := h.service.ConvertToResponse(&device)
	if err != nil {
		common.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	common.SuccessResponse(c, "Create device success", deviceResponse)
}

// Update updates an existing access control device.
func (h *AccessControlDeviceHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var bodyRequest schema.AccessControlDeviceRequest
	if err := c.ShouldBind(&bodyRequest); err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := validate.Struct(bodyRequest); err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	device := model.AccessControlDevice{
		Name:                  bodyRequest.Name,
		Type:                  bodyRequest.Type,
		HostAddress:           bodyRequest.HostAddress,
		Username:              bodyRequest.Username,
		Password:              bodyRequest.Password,
		AccessToken:           bodyRequest.AccessToken,
		ApiToken:              bodyRequest.ApiToken,
		RecordScan:            bodyRequest.RecordScan,
		RecordAttendance:      bodyRequest.RecordAttendance,
		AllowClockIn:          bodyRequest.AllowClockIn,
		AllowClockOut:         bodyRequest.AllowClockOut,
		Status:                bodyRequest.Status,
		AccessControlServerID: bodyRequest.AccessControlServerID,
	}
	device.ID = id

	if err := h.service.Save(id, &device); err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	deviceResponse, err := h.service.ConvertToResponse(&device)
	if err != nil {
		common.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	common.SuccessResponse(c, "Update device success", deviceResponse)
}

// Delete deletes an access control device by its ID.
func (h *AccessControlDeviceHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.Delete(id); err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	common.SuccessResponse(c, "Device deleted successfully", nil)
}
