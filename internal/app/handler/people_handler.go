package handler

import (
	"net/http"
	"strconv"

	"github.com/putteror/access-control-management/internal/app/model"
	"github.com/putteror/access-control-management/internal/app/service"

	"github.com/gin-gonic/gin"
)

type PeopleHandler struct {
	service service.PeopleService
}

func NewPeopleHandler(service service.PeopleService) *PeopleHandler {
	return &PeopleHandler{service: service}
}

func (h *PeopleHandler) Create(c *gin.Context) {
	var people model.People
	if err := c.BindJSON(&people); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreatePeople(&people); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, people)
}

func (h *PeopleHandler) FindAll(c *gin.Context) {
	peoples, err := h.service.GetAllPeople()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, peoples)
}

func (h *PeopleHandler) FindByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	people, err := h.service.GetPeopleByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "People not found"})
		return
	}

	c.JSON(http.StatusOK, people)
}
