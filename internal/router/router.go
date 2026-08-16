// Package router assembles the Gin engine and all route groups for escape-room-ops.
package router

import (
	"log/slog"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/lp/escape-room-ops/internal/config"
	"github.com/lp/escape-room-ops/internal/constants"
	"github.com/lp/escape-room-ops/internal/handler"
	"github.com/lp/escape-room-ops/internal/middleware"
	"github.com/lp/escape-room-ops/internal/repository"
	"github.com/lp/escape-room-ops/internal/service"
	"github.com/lp/escape-room-ops/internal/util"
	"gorm.io/gorm"
)

// New builds the Gin engine with all dependencies wired.
func New(cfg *config.Config, db *gorm.DB, logger *slog.Logger) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.RequestID(logger))
	r.Use(middleware.AccessLog(logger))
	r.Use(cors.New(cors.Config{
		AllowOrigins:  cfg.CORSOrigins,
		AllowMethods:  []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:  []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders: []string{"X-Request-ID"},
	}))
	r.Use(middleware.ErrorHandler())

	r.GET("/healthz", func(c *gin.Context) { util.OK(c, gin.H{"status": "ok"}) })

	// repositories
	userRepo := repository.NewUserRepository(db)
	themeRepo := repository.NewThemeRoomRepository(db)
	sessionRepo := repository.NewGameSessionRepository(db)
	propRepo := repository.NewPropRepository(db)
	escapeRepo := repository.NewEscapeRecordRepository(db)
	revenueRepo := repository.NewRevenueRepository(db)

	// services
	userSvc := service.NewUserService(userRepo, cfg.JWTSecret, cfg.JWTExpireHours, logger)
	themeSvc := service.NewThemeRoomService(themeRepo, logger)
	sessionSvc := service.NewGameSessionService(sessionRepo, themeRepo, userRepo, logger)
	propSvc := service.NewPropService(propRepo, logger)
	escapeSvc := service.NewEscapeRecordService(escapeRepo, themeRepo, sessionRepo, logger)
	revenueSvc := service.NewRevenueService(revenueRepo, sessionRepo, logger)

	// handlers
	userH := handler.NewUserHandler(userSvc, logger)
	themeH := handler.NewThemeRoomHandler(themeSvc, logger)
	sessionH := handler.NewGameSessionHandler(sessionSvc, userSvc, revenueRepo, logger)
	propH := handler.NewPropHandler(propSvc, logger)
	escapeH := handler.NewEscapeRecordHandler(escapeSvc, logger)
	revenueH := handler.NewRevenueHandler(revenueSvc, logger)

	auth := middleware.AuthRequired(cfg.JWTSecret, logger)
	requireAdmin := middleware.RequireRole(logger, constants.UserRoleAdmin)
	loginLimiter := middleware.RateLimit(middleware.NewRateLimiter(cfg.LoginRateLimit, time.Minute), logger)
	apiLimiter := middleware.RateLimit(middleware.NewRateLimiter(cfg.RateLimitPerMin, time.Minute), logger)

	v1 := r.Group("/api/v1")
	{
		RegisterUserRoutes(v1, userH, auth, loginLimiter, apiLimiter)
		RegisterThemeRoomRoutes(v1, themeH, auth, requireAdmin, apiLimiter)
		RegisterGameSessionRoutes(v1, sessionH, auth, requireAdmin, apiLimiter)
		RegisterPropRoutes(v1, propH, auth, requireAdmin, apiLimiter)
		RegisterEscapeRecordRoutes(v1, escapeH, auth, requireAdmin, apiLimiter)
		RegisterRevenueRoutes(v1, revenueH, auth, requireAdmin, apiLimiter)
	}
	return r
}
