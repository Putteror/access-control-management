package router

import (
	"github.com/putteror/access-control-management/internal/app/handler"

	"github.com/gin-gonic/gin"
)

func NewRouter(
	peopleHandler *handler.PeopleHandler,
	timerecordHandler *handler.TimeRecordHandler,
) *gin.Engine {
	router := gin.Default()

	api := router.Group("/api")
	{
		// People endpoints
		people := api.Group("/people")
		{
			people.POST("/", peopleHandler.Create)
			people.GET("/", peopleHandler.FindAll)
			people.GET("/:id", peopleHandler.FindByID)
		}

		// TimeRecord endpoints
		timerecords := api.Group("/timerecords")
		{
			timerecords.POST("/clock-in", timerecordHandler.ClockIn)
			timerecords.POST("/clock-out", timerecordHandler.ClockOut)
		}
	}

	return router
}
