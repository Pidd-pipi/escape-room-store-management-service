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

// PropHandler exposes prop inventory endpoints.
type PropHandler struct {
	svc    *service.PropService
	logger *slog.Logger
}

// NewPropHandler wires the prop handler dependencies.
func NewPropHandler(svc *service.PropService, logger *slog.Logger) *PropHandler {
	return &PropHandler{svc: svc, logger: logger}
}

// Create handles POST /props (admin).
func (h *PropHandler) Create(c *gin.Context) {
	var req dto.CreatePropRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeValidation, constants.MsgValidationFailed)
		return
	}
	p, err := h.svc.Create(c.Request.Context(), &req)
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, p)
}

// List handles GET /props.
func (h *PropHandler) List(c *gin.Context) {
	var q dto.PageQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeValidation, constants.MsgValidationFailed)
		return
	}
	themeID, _ := parseUint(c.Query("theme_room_id"))
	result, err := h.svc.List(c.Request.Context(), themeID, &q)
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, result)
}

// Restock handles POST /props/:id/restock (admin).
func (h *PropHandler) Restock(c *gin.Context) {
	id, err := parseUint(c.Param("id"))
	if err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "道具ID不合法")
		return
	}
	var req struct {
		Stock int `json:"stock" binding:"gte=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeValidation, constants.MsgValidationFailed)
		return
	}
	p, err := h.svc.Restock(c.Request.Context(), id, req.Stock)
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, p)
}

// Consume handles POST /props/:id/consume (admin).
func (h *PropHandler) Consume(c *gin.Context) {
	id, err := parseUint(c.Param("id"))
	if err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "道具ID不合法")
		return
	}
	var req dto.ConsumeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeValidation, constants.MsgValidationFailed)
		return
	}
	p, err := h.svc.Consume(c.Request.Context(), id, req.Quantity)
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, p)
}

// CreateMaintenance handles POST /props/:id/maintenance (admin).
func (h *PropHandler) CreateMaintenance(c *gin.Context) {
	id, err := parseUint(c.Param("id"))
	if err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "道具ID不合法")
		return
	}
	var req dto.CreateMaintenanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeValidation, constants.MsgValidationFailed)
		return
	}
	m, err := h.svc.CreateMaintenance(c.Request.Context(), id, &req)
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, m)
}

// ListMaintenance handles GET /props/:id/maintenance.
func (h *PropHandler) ListMaintenance(c *gin.Context) {
	id, err := parseUint(c.Param("id"))
	if err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "道具ID不合法")
		return
	}
	items, err := h.svc.ListMaintenance(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}
	util.OK(c, items)
}
