package router

import (
	"github.com/gin-gonic/gin"
	"github.com/lp/escape-room-ops/internal/handler"
)

// RegisterPropRoutes registers prop inventory endpoints.
func RegisterPropRoutes(g *gin.RouterGroup, h *handler.PropHandler, auth, requireAdmin, apiLimiter gin.HandlerFunc) {
	props := g.Group("/props", auth, requireAdmin)
	{
		props.GET("", apiLimiter, h.List)
		props.POST("", apiLimiter, h.Create)
		props.POST("/:id/restock", apiLimiter, h.Restock)
		props.POST("/:id/consume", apiLimiter, h.Consume)
		props.POST("/:id/maintenance", apiLimiter, h.CreateMaintenance)
		props.GET("/:id/maintenance", apiLimiter, h.ListMaintenance)
	}
}
