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

type AccessControlServerHandler struct {
	service service.AccessControlServerService
}

func NewAccessControlServerHandler(service service.AccessControlServerService) *AccessControlServerHandler {
	return &AccessControlServerHandler{service: service}
}

func init() {
	validate = validator.New()
}

// GetAll retrieves all access control servers.
func (h *AccessControlServerHandler) GetAll(c *gin.Context) {

	var searchQuery schema.AccessControlServerSearchQuery
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
	servers, err := h.service.GetAll(searchQuery)
	if err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	serverResponses := make([]schema.AccessControlServerResponse, len(servers))
	for i, server := range servers {
		response := schema.AccessControlServerResponse{
			ID:          server.ID,
			Name:        server.Name,
			HostAddress: server.HostAddress,
			Status:      server.Status,
		}
		serverResponses[i] = response
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

	var serverResponse schema.AccessControlServerResponse
	if server != nil {
		serverResponse = schema.AccessControlServerResponse{
			ID:          server.ID,
			Name:        server.Name,
			HostAddress: server.HostAddress,
			Status:      server.Status,
		}
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

	if err := validate.Struct(bodyRequest); err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	server := model.AccessControlServer{
		Name:        bodyRequest.Name,
		HostAddress: bodyRequest.HostAddress,
		Username:    bodyRequest.Username,
		Password:    bodyRequest.Password,
		AccessToken: bodyRequest.AccessToken,
		ApiToken:    bodyRequest.ApiToken,
		Status:      bodyRequest.Status,
	}

	if err := h.service.Save("", &server); err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	common.SuccessResponse(c, "Create server success", server)
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

	server := model.AccessControlServer{
		Name:        bodyRequest.Name,
		HostAddress: bodyRequest.HostAddress,
		Username:    bodyRequest.Username,
		Password:    bodyRequest.Password,
		AccessToken: bodyRequest.AccessToken,
		ApiToken:    bodyRequest.ApiToken,
		Status:      bodyRequest.Status,
	}
	server.ID = id

	if err := h.service.Save(id, &server); err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	common.SuccessResponse(c, "Update server success", server)
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
