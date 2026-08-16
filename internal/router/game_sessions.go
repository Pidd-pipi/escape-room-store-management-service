package router

import (
	"github.com/gin-gonic/gin"
	"github.com/lp/escape-room-ops/internal/handler"
)

// RegisterGameSessionRoutes registers session endpoints.
func RegisterGameSessionRoutes(g *gin.RouterGroup, h *handler.GameSessionHandler, auth, requireAdmin, apiLimiter gin.HandlerFunc) {
	sessions := g.Group("/game-sessions")
	{
		sessions.GET("", apiLimiter, h.List)
		sessions.GET("/:id", apiLimiter, h.Get)
		authed := sessions.Group("", auth)
		{
			authed.POST("/:id/register", apiLimiter, h.Register)
			authed.POST("/:id/quit", apiLimiter, h.Quit)
		}
		admin := sessions.Group("", auth, requireAdmin)
		{
			admin.POST("", apiLimiter, h.Create)
			admin.POST("/:id/lock", apiLimiter, h.Lock)
			admin.POST("/:id/cancel", apiLimiter, h.Cancel)
			admin.POST("/:id/finish", apiLimiter, h.Finish)
		}
	}
}
