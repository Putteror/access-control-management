package router

import (
	"github.com/putteror/access-control-management/internal/app/handler"

	"github.com/gin-gonic/gin"
)

func NewRouter(
	deviceHandler *handler.DeviceHandler,
	peopleHandler *handler.PeopleHandler,
) *gin.Engine {
	router := gin.Default()

	api := router.Group("/api")
	{
		// Device endpoints
		devices := api.Group("/device")
		{
			devices.POST("/", deviceHandler.Create)
			devices.GET("/", deviceHandler.FindAll)
		}

		// People endpoints
		people := api.Group("/people")
		{
			people.POST("/", peopleHandler.Create)
			people.GET("/", peopleHandler.FindAll)
			people.GET("/:id", peopleHandler.FindByID)
		}

	}

	return router
}
