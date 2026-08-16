package router

import (
	"github.com/gin-gonic/gin"
	"github.com/lp/escape-room-ops/internal/handler"
)

// RegisterRevenueRoutes registers revenue analytics endpoints (admin).
func RegisterRevenueRoutes(g *gin.RouterGroup, h *handler.RevenueHandler, auth, requireAdmin, apiLimiter gin.HandlerFunc) {
	revenues := g.Group("/revenues", auth, requireAdmin)
	{
		revenues.GET("/analytics", apiLimiter, h.Analyze)
		revenues.GET("/export", apiLimiter, h.Export)
	}
}
