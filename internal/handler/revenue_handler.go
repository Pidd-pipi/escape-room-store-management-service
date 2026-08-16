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

// RevenueHandler exposes revenue analytics endpoints.
type RevenueHandler struct {
	svc    *service.RevenueService
	logger *slog.Logger
}

// NewRevenueHandler wires the revenue handler dependencies.
func NewRevenueHandler(svc *service.RevenueService, logger *slog.Logger) *RevenueHandler {
	return &RevenueHandler{svc: svc, logger: logger}
}

// Analyze handles GET /revenues/analytics.
func (h *RevenueHandler) Analyze(c *gin.Context) {
	var q dto.RevenueQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeValidation, constants.MsgValidationFailed)
		return
	}
	result, err := h.svc.Analyze(c.Request.Context(), &q)
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, result)
}

// Export handles GET /revenues/export.
func (h *RevenueHandler) Export(c *gin.Context) {
	var q dto.RevenueQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeValidation, constants.MsgValidationFailed)
		return
	}
	rows, err := h.svc.Export(c.Request.Context(), &q)
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, rows)
}
