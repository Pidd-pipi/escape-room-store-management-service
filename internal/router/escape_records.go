package router

import (
	"github.com/gin-gonic/gin"
	"github.com/lp/escape-room-ops/internal/handler"
)

// RegisterEscapeRecordRoutes registers escape record and leaderboard endpoints.
func RegisterEscapeRecordRoutes(g *gin.RouterGroup, h *handler.EscapeRecordHandler, auth, requireAdmin, apiLimiter gin.HandlerFunc) {
	records := g.Group("/escape-records")
	{
		records.GET("/leaderboard", apiLimiter, h.Leaderboard)
		admin := records.Group("", auth, requireAdmin)
		{
			admin.POST("", apiLimiter, h.Create)
		}
	}
}
