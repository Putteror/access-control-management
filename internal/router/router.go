package router

import (
	"github.com/putteror/access-control-management/internal/app/handler"
	"github.com/putteror/access-control-management/internal/app/middleware"

	"github.com/gin-gonic/gin"
)

func NewRouter(
	accessControlDeviceHandler *handler.AccessControlDeviceHandler,
	accessControlGroupHandler *handler.AccessControlGroupHandler,
	accessControlRuleHandler *handler.AccessControlRuleHandler,
	accessControlServerHandler *handler.AccessControlServerHandler,
	accessRecordHandler *handler.AccessRecordHandler,
	attendanceHandler *handler.AttendanceHandler,
	authHandler *handler.AuthHandler,
	peopleHandler *handler.PersonHandler,
	userHandler *handler.UserHandler,
) *gin.Engine {
	router := gin.Default()
	router.POST("/login", authHandler.Login)

	api := router.Group("/api")
	api.Use(middleware.JWTAuthMiddleware())
	{

		// Access Control Device endpoints
		accessControlDevice := api.Group("/access-control-devices")
		{
			accessControlDevice.GET("/", accessControlDeviceHandler.GetAll)
			accessControlDevice.GET("/:id", accessControlDeviceHandler.GetByID)
			accessControlDevice.POST("/", accessControlDeviceHandler.Create)
			accessControlDevice.PUT("/:id", accessControlDeviceHandler.Update)
			accessControlDevice.PATCH("/:id", accessControlDeviceHandler.PartialUpdate)
			accessControlDevice.DELETE("/:id", accessControlDeviceHandler.Delete)
		}

		// Access Control Group endpoints
		accessControlGroup := api.Group("/access-control-groups")
		{
			accessControlGroup.GET("/", accessControlGroupHandler.GetAll)
			accessControlGroup.GET("/:id", accessControlGroupHandler.GetByID)
			accessControlGroup.POST("/", accessControlGroupHandler.Create)
			accessControlGroup.PUT("/:id", accessControlGroupHandler.Update)
			accessControlGroup.PATCH("/:id", accessControlGroupHandler.PartialUpdate)
			accessControlGroup.DELETE("/:id", accessControlGroupHandler.Delete)
		}

		// Access Control Rule endpoints
		accessControlRule := api.Group("/access-control-rules")
		{
			accessControlRule.GET("/", accessControlRuleHandler.GetAll)
			accessControlRule.GET("/:id", accessControlRuleHandler.GetByID)
			accessControlRule.POST("/", accessControlRuleHandler.Create)
			accessControlRule.PUT("/:id", accessControlRuleHandler.Update)
			accessControlRule.PATCH("/:id", accessControlRuleHandler.PartialUpdate)
			accessControlRule.DELETE("/:id", accessControlRuleHandler.Delete)
		}

		// Access Control Server endpoints
		accessControlServer := api.Group("/access-control-servers")
		{
			accessControlServer.GET("/", accessControlServerHandler.GetAll)
			accessControlServer.GET("/:id", accessControlServerHandler.GetByID)
			accessControlServer.POST("/", accessControlServerHandler.Create)
			accessControlServer.PUT("/:id", accessControlServerHandler.Update)
			accessControlServer.PATCH("/:id", accessControlServerHandler.PartialUpdate)
			accessControlServer.DELETE("/:id", accessControlServerHandler.Delete)
		}

		// Access record endpoints
		accessRecord := api.Group("/access-records")
		{
			accessRecord.GET("/", accessRecordHandler.GetAll)
			accessRecord.GET("/:id", accessRecordHandler.GetByID)
			accessRecord.POST("/", accessRecordHandler.Create)
			accessRecord.PUT("/:id", accessRecordHandler.Update)
			accessRecord.PATCH("/:id", accessRecordHandler.PartialUpdate)
			accessRecord.DELETE("/:id", accessRecordHandler.Delete)
		}

		// Attendance endpoints
		attendance := api.Group("/attendances")
		{
			attendance.GET("/", attendanceHandler.GetAll)
			attendance.GET("/:id", attendanceHandler.GetByID)
			attendance.POST("/", attendanceHandler.Create)
			attendance.PUT("/:id", attendanceHandler.Update)
			attendance.PATCH("/:id", attendanceHandler.PartialUpdate)
			attendance.DELETE("/:id", attendanceHandler.Delete)
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

		// User
		user := api.Group("/users")
		{
			user.GET("/", userHandler.GetAll)
			user.GET("/:id", userHandler.GetByID)
			user.POST("/", userHandler.Create)
			user.PUT("/:id", userHandler.Update)
			user.PATCH("/:id", userHandler.PartialUpdate)
			user.DELETE("/:id", userHandler.Delete)
		}

	}

	return router
}
