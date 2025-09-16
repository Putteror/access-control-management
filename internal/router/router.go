package router

import (
	"github.com/putteror/access-control-management/internal/app/handler"

	"github.com/gin-gonic/gin"
)

func NewRouter(
	accessControlDeviceHandler *handler.AccessControlDeviceHandler,
	accessControlRuleHandler *handler.AccessControlRuleHandler,
	accessControlServerHandler *handler.AccessControlServerHandler,
	peopleHandler *handler.PersonHandler,
) *gin.Engine {
	router := gin.Default()

	api := router.Group("/api")
	{

		// Access Control Device endpoints
		accessControlDevice := api.Group("/access-control-devices")
		{
			accessControlDevice.GET("/", accessControlDeviceHandler.GetAll)
			accessControlDevice.GET("/:id", accessControlDeviceHandler.GetByID)
			accessControlDevice.POST("/", accessControlDeviceHandler.Create)
			accessControlDevice.PUT("/:id", accessControlDeviceHandler.Update)
			accessControlDevice.DELETE("/:id", accessControlDeviceHandler.Delete)
		}

		// Access Control Rule endpoints
		accessControlRule := api.Group("/access-control-rules")
		{
			accessControlRule.GET("/", accessControlRuleHandler.GetAll)
			accessControlRule.GET("/:id", accessControlRuleHandler.GetByID)
			accessControlRule.POST("/", accessControlRuleHandler.Create)
			accessControlRule.PUT("/:id", accessControlRuleHandler.Update)
			accessControlRule.DELETE("/:id", accessControlRuleHandler.Delete)
		}

		// Access Control Server endpoints
		accessControlServer := api.Group("/access-control-servers")
		{
			accessControlServer.GET("/", accessControlServerHandler.GetAll)
			accessControlServer.GET("/:id", accessControlServerHandler.GetByID)
			accessControlServer.POST("/", accessControlServerHandler.Create)
			accessControlServer.PUT("/:id", accessControlServerHandler.Update)
			accessControlServer.DELETE("/:id", accessControlServerHandler.Delete)
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
