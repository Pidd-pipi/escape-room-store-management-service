package handler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lp/escape-room-ops/internal/constants"
	"github.com/lp/escape-room-ops/internal/dto"
	"github.com/lp/escape-room-ops/internal/middleware"
	"github.com/lp/escape-room-ops/internal/model"
	"github.com/lp/escape-room-ops/internal/repository"
	"github.com/lp/escape-room-ops/internal/service"
	"github.com/lp/escape-room-ops/internal/util"
)

// GameSessionHandler exposes session endpoints.
type GameSessionHandler struct {
	svc     *service.GameSessionService
	users   *service.UserService
	revenue *repository.RevenueRepository
	logger  *slog.Logger
}

// NewGameSessionHandler wires the game session handler dependencies.
func NewGameSessionHandler(svc *service.GameSessionService, users *service.UserService, revenue *repository.RevenueRepository, logger *slog.Logger) *GameSessionHandler {
	return &GameSessionHandler{svc: svc, users: users, revenue: revenue, logger: logger}
}

// Create handles POST /game-sessions (admin).
func (h *GameSessionHandler) Create(c *gin.Context) {
	var req dto.CreateGameSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeValidation, constants.MsgValidationFailed)
		return
	}
	session, err := h.svc.Create(c.Request.Context(), &req)
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, session)
}

// List handles GET /game-sessions.
func (h *GameSessionHandler) List(c *gin.Context) {
	var q dto.ListSessionQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeValidation, constants.MsgValidationFailed)
		return
	}
	result, err := h.svc.List(c.Request.Context(), &q)
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, result)
}

// Get handles GET /game-sessions/:id.
func (h *GameSessionHandler) Get(c *gin.Context) {
	id, err := parseUint(c.Param("id"))
	if err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "场次ID不合法")
		return
	}
	session, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, session)
}

// Register handles POST /game-sessions/:id/register.
func (h *GameSessionHandler) Register(c *gin.Context) {
	user := h.requireUser(c)
	if user == nil {
		return
	}
	id, err := parseUint(c.Param("id"))
	if err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "场次ID不合法")
		return
	}
	session, err := h.svc.Register(c.Request.Context(), user, id)
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, session)
}

// Quit handles POST /game-sessions/:id/quit.
func (h *GameSessionHandler) Quit(c *gin.Context) {
	userID, err := middleware.CurrentUserID(c)
	if err != nil {
		util.Fail(c, http.StatusUnauthorized, constants.CodeUnauthorized, constants.MsgUnauthorized)
		return
	}
	id, err := parseUint(c.Param("id"))
	if err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "场次ID不合法")
		return
	}
	session, err := h.svc.Quit(c.Request.Context(), userID, id)
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, session)
}

// Lock handles POST /game-sessions/:id/lock (admin).
func (h *GameSessionHandler) Lock(c *gin.Context) {
	h.act(c, h.svc.Lock)
}

// Cancel handles POST /game-sessions/:id/cancel (admin).
func (h *GameSessionHandler) Cancel(c *gin.Context) {
	h.act(c, h.svc.Cancel)
}

// Finish handles POST /game-sessions/:id/finish (admin).
func (h *GameSessionHandler) Finish(c *gin.Context) {
	id, err := parseUint(c.Param("id"))
	if err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "场次ID不合法")
		return
	}
	session, err := h.svc.Finish(c.Request.Context(), id, 68, h.revenue)
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, session)
}

func (h *GameSessionHandler) act(c *gin.Context, fn func(ctx context.Context, id uint) (*model.GameSession, error)) {
	id, err := parseUint(c.Param("id"))
	if err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "场次ID不合法")
		return
	}
	session, err := fn(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, session)
}

func (h *GameSessionHandler) requireUser(c *gin.Context) *model.User {
	userID, err := middleware.CurrentUserID(c)
	if err != nil {
		util.Fail(c, http.StatusUnauthorized, constants.CodeUnauthorized, constants.MsgUnauthorized)
		return nil
	}
	user, err := h.users.GetProfile(c.Request.Context(), userID)
	if err != nil {
		c.Error(err)
		return nil
	}
	return user
}
