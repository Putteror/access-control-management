package router

import (
	"github.com/putteror/access-control-management/internal/app/handler"

	"github.com/gin-gonic/gin"
)

func NewRouter(
	accessControlRuleHandler *handler.AccessControlRuleHandler,
	peopleHandler *handler.PersonHandler,
) *gin.Engine {
	router := gin.Default()

	api := router.Group("/api")
	{
		// Access Control Rule endpoints
		accessControlRule := api.Group("/access-control-rules")
		{
			accessControlRule.GET("/", accessControlRuleHandler.GetAll)
			accessControlRule.GET("/:id", accessControlRuleHandler.GetByID)
			accessControlRule.POST("/", accessControlRuleHandler.Create)
			accessControlRule.PUT("/:id", accessControlRuleHandler.Update)
			accessControlRule.DELETE("/:id", accessControlRuleHandler.Delete)
		}

		// People endpoints
		people := api.Group("/people")
		{
			people.GET("/", peopleHandler.GetAll)
			people.GET("/:id", peopleHandler.GetByID)
			people.POST("/", peopleHandler.Create)
			people.PUT("/:id", peopleHandler.Update)
			people.DELETE("/:id", peopleHandler.Delete)
		}

	}

	return router
}
