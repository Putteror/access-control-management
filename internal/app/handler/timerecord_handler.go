package handler

import (
	"net/http"

	"github.com/putteror/access-control-management/internal/app/service"

	"github.com/gin-gonic/gin"
)

type TimeRecordHandler struct {
	service service.TimeRecordService
}

func NewTimeRecordHandler(service service.TimeRecordService) *TimeRecordHandler {
	return &TimeRecordHandler{service: service}
}

func (h *TimeRecordHandler) ClockIn(c *gin.Context) {
	var req struct {
		PeopleID uint `json:"people_id"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.ClockIn(req.PeopleID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Clock-in successful"})
}

func (h *TimeRecordHandler) ClockOut(c *gin.Context) {
	var req struct {
		PeopleID uint `json:"people_id"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.ClockOut(req.PeopleID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Clock-out successful"})
}
