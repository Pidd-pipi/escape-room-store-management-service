package router

import (
	"github.com/gin-gonic/gin"
	"github.com/lp/escape-room-ops/internal/handler"
)

// RegisterThemeRoomRoutes registers theme room endpoints.
func RegisterThemeRoomRoutes(g *gin.RouterGroup, h *handler.ThemeRoomHandler, auth, requireAdmin, apiLimiter gin.HandlerFunc) {
	themes := g.Group("/theme-rooms")
	{
		themes.GET("", apiLimiter, h.List)
		themes.GET("/:id", apiLimiter, h.Get)
		admin := themes.Group("", auth, requireAdmin)
		{
			admin.POST("", apiLimiter, h.Create)
		}
	}
}
