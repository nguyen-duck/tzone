package route

import (
	"github.com/LuuDinhTheTai/tzone/internal/delivery/handler"
	"github.com/LuuDinhTheTai/tzone/internal/delivery/middleware"
	"github.com/LuuDinhTheTai/tzone/internal/service"
	"github.com/gin-gonic/gin"
)

func MapDeviceRoutes(r *gin.Engine, deviceHandler *handler.DeviceHandler, authService *service.AuthService, permissionService *service.PermissionService) {
	deviceGroup := r.Group("/api/v1/devices")
	{
		deviceGroup.GET("", deviceHandler.GetAllDevices)
		deviceGroup.GET("/:id", deviceHandler.GetDeviceById)

		// Protected endpoints
		deviceGroup.POST("", middleware.JWTAuth(authService), middleware.RBACAuth(permissionService), deviceHandler.CreateDevice)
		deviceGroup.PUT("/:id", middleware.JWTAuth(authService), middleware.RBACAuth(permissionService), deviceHandler.UpdateDevice)
		deviceGroup.DELETE("/:id", middleware.JWTAuth(authService), middleware.RBACAuth(permissionService), deviceHandler.DeleteDevice)
	}
}
