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

type AccessRecordHandler struct {
	service service.AccessRecordService
}

func NewAccessRecordHandler(service service.AccessRecordService) *AccessRecordHandler {
	return &AccessRecordHandler{service: service}
}

func init() {
	validate = validator.New()
}

// Get All
func (h *AccessRecordHandler) GetAll(c *gin.Context) {

	var searchQuery schema.AccessRecordSearchQuery
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
	records, err := h.service.GetAll(searchQuery)
	if err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	recordResponses := make([]schema.AccessRecordResponse, len(records))
	for i, record := range records {
		response, err := h.service.ConvertToResponse(&record)
		if err != nil {
			common.ErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		recordResponses[i] = *response
	}

	pageData := common.PageResponse{
		Page:      searchQuery.Page,
		Size:      searchQuery.Limit,
		Total:     len(records),
		TotalPage: (len(records) + searchQuery.Limit - 1) / searchQuery.Limit,
	}

	common.GetDataListResponse(c, "Success", recordResponses, pageData)

}

// GetByID retrieves an record by its ID.
func (h *AccessRecordHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	itemModel, err := h.service.GetByID(id)
	if err != nil {
		common.ErrorResponse(c, http.StatusNotFound, "Record not found")
		return
	}

	itemResponse, err := h.service.ConvertToResponse(itemModel)
	if err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	common.SuccessResponse(c, "Success", itemResponse)
}

func (h *AccessRecordHandler) Create(c *gin.Context) {
	var bodyRequest schema.AccessRecordRequest
	if err := c.ShouldBindJSON(&bodyRequest); err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := validate.Struct(bodyRequest); err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	itemModel, err := h.service.Create(&bodyRequest)
	if err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	itemResponse, err := h.service.ConvertToResponse(itemModel)
	if err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	common.SuccessResponse(c, "Create record success", itemResponse)
}

// Update updates an existing item (Full Update).
func (h *AccessRecordHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var bodyRequest schema.AccessRecordRequest
	if err := c.ShouldBindJSON(&bodyRequest); err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := validate.Struct(bodyRequest); err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	itemModel, err := h.service.Update(id, &bodyRequest)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			common.ErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		handleErrorResponse(c, err, err.Error())
		return
	}

	itemResponse, err := h.service.ConvertToResponse(itemModel)
	if err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	common.SuccessResponse(c, "Update item success", itemResponse)
}

// PartialUpdate updates an existing item (Partial Update).
func (h *AccessRecordHandler) PartialUpdate(c *gin.Context) {
	id := c.Param("id")

	var bodyRequest schema.AccessRecordRequest
	if err := c.ShouldBindJSON(&bodyRequest); err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	itemModel, err := h.service.PartialUpdate(id, &bodyRequest)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			common.ErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		handleErrorResponse(c, err, err.Error())
		return
	}

	itemResponse, err := h.service.ConvertToResponse(itemModel)
	if err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	common.SuccessResponse(c, "Update item success", itemResponse)
}

// Delete deletes an item by its ID.
func (h *AccessRecordHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.Delete(id); err != nil {
		if strings.Contains(err.Error(), "not found") {
			common.ErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		common.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	common.SuccessResponse(c, "Item deleted successfully", nil)
}
