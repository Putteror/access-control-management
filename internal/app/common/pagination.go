package common

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type PageResponse struct {
	Page      int `json:"page"`
	Size      int `json:"size"`
	Total     int `json:"total"`
	TotalPage int `json:"totalPage"`
}

var DefaultPage = 1
var DefaultPageSize = 10

// GetPaginationParams extracts and validates pagination parameters from the request.
// It returns the page and limit as integers, and an error if the parameters are invalid.
func GetPaginationParams(c *gin.Context) (page, limit int, err error) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err = strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		return 0, 0, err
	}
	limit, err = strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		return 0, 0, err
	}

	return page, limit, nil
}
