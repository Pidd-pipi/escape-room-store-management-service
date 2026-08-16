package handler

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lp/escape-room-ops/internal/constants"
	"github.com/lp/escape-room-ops/internal/dto"
	"github.com/lp/escape-room-ops/internal/service"
	"github.com/lp/escape-room-ops/internal/util"
)

// ThemeRoomHandler exposes theme room endpoints.
type ThemeRoomHandler struct {
	svc    *service.ThemeRoomService
	logger *slog.Logger
}

// NewThemeRoomHandler wires the theme room handler dependencies.
func NewThemeRoomHandler(svc *service.ThemeRoomService, logger *slog.Logger) *ThemeRoomHandler {
	return &ThemeRoomHandler{svc: svc, logger: logger}
}

// Create handles POST /theme-rooms (admin).
func (h *ThemeRoomHandler) Create(c *gin.Context) {
	var req dto.CreateThemeRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeValidation, constants.MsgValidationFailed)
		return
	}
	theme, err := h.svc.Create(c.Request.Context(), &req)
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, theme)
}

// Get handles GET /theme-rooms/:id.
func (h *ThemeRoomHandler) Get(c *gin.Context) {
	id, err := parseUint(c.Param("id"))
	if err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "主题ID不合法")
		return
	}
	theme, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, theme)
}

// List handles GET /theme-rooms.
func (h *ThemeRoomHandler) List(c *gin.Context) {
	var q dto.ListThemeQuery
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
