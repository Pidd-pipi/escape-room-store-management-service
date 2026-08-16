package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/lp/escape-room-ops/internal/constants"
	"github.com/lp/escape-room-ops/internal/dto"
	"github.com/lp/escape-room-ops/internal/model"
	"github.com/lp/escape-room-ops/internal/repository"
	"github.com/lp/escape-room-ops/internal/util"
)

// PropService manages prop inventory, consumption and maintenance.
type PropService struct {
	props  *repository.PropRepository
	logger *slog.Logger
}

// NewPropService wires the prop service dependencies.
func NewPropService(props *repository.PropRepository, logger *slog.Logger) *PropService {
	return &PropService{props: props, logger: logger}
}

// Create adds a new prop.
func (s *PropService) Create(ctx context.Context, req *dto.CreatePropRequest) (*model.Prop, error) {
	p := &model.Prop{
		ThemeRoomID: req.ThemeRoomID, Name: req.Name, Category: req.Category,
		Stock: req.Stock, AlertThreshold: req.AlertThreshold, Status: constants.PropStatusNormal,
	}
	if err := s.props.Transaction(ctx, func(txCtx context.Context) error {
		if err := s.props.Create(txCtx, p); err != nil {
			return err
		}
		if err := s.props.RefreshStatus(txCtx, p.ID); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, util.WrapAppError(fmt.Errorf("prop[name=%s] create: %w", req.Name, err), 500, constants.CodeInternalError, constants.MsgInternalError)
	}
	s.logger.Info(fmt.Sprintf(constants.LogPropCreateSuccess, p.ID, p.Name))
	return s.props.FindByID(ctx, p.ID)
}

// List returns props filtered by theme.
func (s *PropService) List(ctx context.Context, themeRoomID uint, q *dto.PageQuery) (*dto.PageResult, error) {
	q.Normalize()
	items, total, err := s.props.List(ctx, themeRoomID, q.Page, 0)
	if err != nil {
		return nil, util.WrapAppError(fmt.Errorf("prop list: %w", err), 500, constants.CodeInternalError, constants.MsgInternalError)
	}
	return &dto.PageResult{Items: items, Total: total, Page: q.Page, PageSize: q.PageSize}, nil
}

// Restock increases stock and refreshes the status.
func (s *PropService) Restock(ctx context.Context, propID uint, stock int) (*model.Prop, error) {
	p, err := s.props.FindByID(ctx, propID)
	if err != nil {
		return nil, util.WrapAppError(fmt.Errorf("prop[id=%d] restock find: %w", propID, err), 404, constants.CodeNotFound, constants.MsgNotFound)
	}
	newStock := p.Stock + stock
	if err := s.props.Transaction(ctx, func(txCtx context.Context) error {
		if err := s.props.UpdateStock(txCtx, propID, newStock); err != nil {
			return err
		}
		if err := s.props.RefreshStatus(txCtx, propID); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, util.WrapAppError(fmt.Errorf("prop[id=%d] restock: %w", propID, err), 500, constants.CodeInternalError, constants.MsgInternalError)
	}
	s.logger.Info(fmt.Sprintf(constants.LogPropStockAlertSuccess, propID, newStock))
	return s.props.FindByID(ctx, propID)
}

// Consume decreases stock by quantity and alerts on low stock.
func (s *PropService) Consume(ctx context.Context, propID uint, quantity int) (*model.Prop, error) {
	p, err := s.props.FindByID(ctx, propID)
	if err != nil {
		return nil, util.WrapAppError(fmt.Errorf("prop[id=%d] consume find: %w", propID, err), 404, constants.CodeNotFound, constants.MsgNotFound)
	}
	if p.Stock < quantity {
		return nil, util.NewAppError(409, constants.CodeConflict, constants.MsgPropLowStock, nil)
	}
	if err := s.props.Transaction(ctx, func(txCtx context.Context) error {
		if err := s.props.Consume(txCtx, propID, quantity); err != nil {
			return err
		}
		if err := s.props.RefreshStatus(txCtx, propID); err != nil {
			return err
		}
		return nil
	}); err != nil {
		if errors.Is(err, util.ErrConflict) {
			return nil, util.NewAppError(409, constants.CodeConflict, constants.MsgPropLowStock, nil)
		}
		return nil, util.WrapAppError(fmt.Errorf("prop[id=%d] consume: %w", propID, err), 500, constants.CodeInternalError, constants.MsgInternalError)
	}
	updated, err := s.props.FindByID(ctx, propID)
	if err != nil {
		return nil, err
	}
	if updated.Status == constants.PropStatusLowStock {
		s.logger.Warn(fmt.Sprintf(constants.LogPropStockAlertSuccess, propID, updated.Stock))
	}
	return updated, nil
}

// CreateMaintenance records maintenance and marks the prop as repairing.
func (s *PropService) CreateMaintenance(ctx context.Context, propID uint, req *dto.CreateMaintenanceRequest) (*model.PropMaintenance, error) {
	if _, err := s.props.FindByID(ctx, propID); err != nil {
		return nil, util.WrapAppError(fmt.Errorf("prop[id=%d] maintenance find: %w", propID, err), 404, constants.CodeNotFound, constants.MsgNotFound)
	}
	m := &model.PropMaintenance{PropID: propID, Type: req.Type, Note: req.Note}
	if err := s.props.Transaction(ctx, func(txCtx context.Context) error {
		if err := s.props.CreateMaintenance(txCtx, m); err != nil {
			return err
		}
		if err := s.props.MarkRepairing(txCtx, propID); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, util.WrapAppError(fmt.Errorf("prop[id=%d] create maintenance: %w", propID, err), 500, constants.CodeInternalError, constants.MsgInternalError)
	}
	s.logger.Info(fmt.Sprintf(constants.LogPropMaintenanceCreate, propID, req.Type))
	return m, nil
}

// ListMaintenance returns maintenance records of a prop.
func (s *PropService) ListMaintenance(ctx context.Context, propID uint) ([]model.PropMaintenance, error) {
	items, err := s.props.ListMaintenance(ctx, propID)
	if err != nil {
		return nil, util.WrapAppError(fmt.Errorf("prop[id=%d] list maintenance: %w", propID, err), 500, constants.CodeInternalError, constants.MsgInternalError)
	}
	return items, nil
}
