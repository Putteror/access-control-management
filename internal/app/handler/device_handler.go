package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/putteror/access-control-management/internal/app/model"
	"github.com/putteror/access-control-management/internal/app/service"
)

type DeviceHandler struct {
	service service.DeviceService
}

func NewDeviceHandler(service service.DeviceService) *DeviceHandler {
	return &DeviceHandler{service: service}
}

func (h *DeviceHandler) Create(c *gin.Context) {
	var device model.Device
	if err := c.BindJSON(&device); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateDevice(&device); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, device)
}

func (h *DeviceHandler) FindAll(c *gin.Context) {
	devices, err := h.service.GetAllDevice()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, devices)
}
