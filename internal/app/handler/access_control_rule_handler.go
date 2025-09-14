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

type AccessControlRuleHandler struct {
	service service.AccessControlRuleService
}

func NewAccessControlRuleHandler(service service.AccessControlRuleService) *AccessControlRuleHandler {
	return &AccessControlRuleHandler{service: service}
}

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func (h *AccessControlRuleHandler) GetAll(c *gin.Context) {

	var searchQuery schema.AccessControlRuleSearchQuery
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
	rules, err := h.service.GetAll(searchQuery)
	if err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	ruleResponses := make([]schema.AccessControlRuleResponse, len(rules))
	for i, rule := range rules {
		response := schema.AccessControlRuleResponse{
			ID:   rule.ID,
			Name: rule.Name,
		}
		ruleResponses[i] = response
	}

	pageData := common.PageResponse{
		Page:      searchQuery.Page,
		Size:      searchQuery.Limit,
		Total:     len(rules),
		TotalPage: (len(rules) + searchQuery.Limit - 1) / searchQuery.Limit,
	}

	common.GetDataListResponse(c, "Success", ruleResponses, pageData)
}

func (h *AccessControlRuleHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	rule, err := h.service.GetByID(id)
	if err != nil {
		common.ErrorResponse(c, http.StatusNotFound, "Rule not found")
		return
	}

	var ruleResponse schema.AccessControlRuleResponse
	if rule != nil {
		ruleResponse = schema.AccessControlRuleResponse{
			ID:   rule.ID,
			Name: rule.Name,
		}
	}

	common.SuccessResponse(c, "Success", ruleResponse)
}

func (h *AccessControlRuleHandler) Create(c *gin.Context) {
	var bodyRequest schema.AccessControlRuleRequest
	if err := c.ShouldBindJSON(&bodyRequest); err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, err.Error())
		common.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := validate.Struct(bodyRequest); err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	rule := model.AccessControlRule{
		Name: bodyRequest.Name,
	}

	if err := h.service.Save("", &rule); err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	common.SuccessResponse(c, "Create rule success", rule)
}

func (h *AccessControlRuleHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var bodyRequest schema.AccessControlRuleRequest
	if err := c.ShouldBind(&bodyRequest); err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := validate.Struct(bodyRequest); err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	rule := model.AccessControlRule{
		Name: bodyRequest.Name,
	}
	rule.ID = id

	if err := h.service.Save(id, &rule); err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	common.SuccessResponse(c, "Create rule success", rule)
}

func (h *AccessControlRuleHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.Delete(id); err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	common.SuccessResponse(c, "Rule deleted successfully", nil)
}
