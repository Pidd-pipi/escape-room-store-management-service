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

// EscapeRecordHandler exposes escape record and leaderboard endpoints.
type EscapeRecordHandler struct {
	svc    *service.EscapeRecordService
	logger *slog.Logger
}

// NewEscapeRecordHandler wires the escape record handler dependencies.
func NewEscapeRecordHandler(svc *service.EscapeRecordService, logger *slog.Logger) *EscapeRecordHandler {
	return &EscapeRecordHandler{svc: svc, logger: logger}
}

// Create handles POST /escape-records (admin).
func (h *EscapeRecordHandler) Create(c *gin.Context) {
	var req dto.CreateEscapeRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeValidation, constants.MsgValidationFailed)
		return
	}
	_, err := h.svc.Create(c.Request.Context(), &req)
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, nil)
}

// Leaderboard handles GET /escape-records/leaderboard.
func (h *EscapeRecordHandler) Leaderboard(c *gin.Context) {
	result, err := h.svc.Leaderboard(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, result)
}
