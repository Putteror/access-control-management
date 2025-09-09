package router

import (
	"github.com/putteror/access-control-management/internal/app/handler"

	"github.com/gin-gonic/gin"
)

func NewRouter(
	peopleHandler *handler.PersonHandler,
) *gin.Engine {
	router := gin.Default()

	api := router.Group("/api")
	{

		// People endpoints
		people := api.Group("/people")
		{
			people.GET("/", peopleHandler.FindAll)
			people.GET("/:id", peopleHandler.FindByID)
			people.POST("/", peopleHandler.Create)
			people.DELETE("/:id", peopleHandler.DeletePerson)
		}

	}

	return router
}
