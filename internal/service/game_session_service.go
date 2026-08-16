package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/lp/escape-room-ops/internal/constants"
	"github.com/lp/escape-room-ops/internal/dto"
	"github.com/lp/escape-room-ops/internal/model"
	"github.com/lp/escape-room-ops/internal/repository"
	"github.com/lp/escape-room-ops/internal/util"
)

// GameSessionService manages scheduled sessions, registrations and state transitions.
type GameSessionService struct {
	sessions *repository.GameSessionRepository
	themes   *repository.ThemeRoomRepository
	users    *repository.UserRepository
	logger   *slog.Logger
}

// NewGameSessionService wires the game session service dependencies.
func NewGameSessionService(sessions *repository.GameSessionRepository, themes *repository.ThemeRoomRepository, users *repository.UserRepository, logger *slog.Logger) *GameSessionService {
	return &GameSessionService{sessions: sessions, themes: themes, users: users, logger: logger}
}

// Create schedules a new open session.
func (s *GameSessionService) Create(ctx context.Context, req *dto.CreateGameSessionRequest) (*model.GameSession, error) {
	theme, err := s.themes.FindByID(ctx, req.ThemeRoomID)
	if err != nil {
		return nil, util.WrapAppError(fmt.Errorf("game_session[theme=%d] theme lookup: %w", req.ThemeRoomID, err), 404, constants.CodeNotFound, constants.MsgNotFound)
	}
	maxPlayers := req.MaxPlayers
	if maxPlayers == 0 {
		maxPlayers = theme.MaxPlayers
	}
	session := &model.GameSession{
		ThemeRoomID: req.ThemeRoomID, StartTime: req.StartTime,
		MaxPlayers: maxPlayers, BookedCount: 0, Status: constants.GameSessionStatusOpen,
	}
	if err := s.sessions.Create(ctx, session); err != nil {
		s.logger.Error(fmt.Sprintf(constants.LogGameSessionCreateFailed, req.ThemeRoomID, err))
		return nil, util.WrapAppError(fmt.Errorf("game_session[theme=%d] create: %w", req.ThemeRoomID, err), 500, constants.CodeInternalError, constants.MsgInternalError)
	}
	s.logger.Info(fmt.Sprintf(constants.LogGameSessionCreateSuccess, session.ID, req.ThemeRoomID))
	return session, nil
}

// List returns sessions with filters.
func (s *GameSessionService) List(ctx context.Context, q *dto.ListSessionQuery) (*dto.PageResult, error) {
	q.Normalize()
	items, total, err := s.sessions.List(ctx, q.ThemeRoomID, q.Status, q.Page, q.PageSize)
	if err != nil {
		return nil, util.WrapAppError(fmt.Errorf("game_session list: %w", err), 500, constants.CodeInternalError, constants.MsgInternalError)
	}
	return &dto.PageResult{Items: items, Total: total, Page: q.Page, PageSize: q.PageSize}, nil
}

// Get returns one session.
func (s *GameSessionService) Get(ctx context.Context, id uint) (*model.GameSession, error) {
	session, err := s.sessions.FindByID(ctx, id)
	if err != nil {
		return nil, util.WrapAppError(fmt.Errorf("game_session[id=%d] get: %w", id, err), 404, constants.CodeNotFound, constants.MsgNotFound)
	}
	return session, nil
}

// Register lets a player sign up for an open session; locks when full.
func (s *GameSessionService) Register(ctx context.Context, user *model.User, sessionID uint) (*model.GameSession, error) {
	session, err := s.sessions.FindByID(ctx, sessionID)
	if err != nil {
		return nil, util.WrapAppError(fmt.Errorf("game_session[id=%d] register find: %w", sessionID, err), 404, constants.CodeNotFound, constants.MsgNotFound)
	}
	if session.Status != constants.GameSessionStatusOpen {
		return nil, util.NewAppError(409, constants.CodeConflict, constants.MsgSessionNotOpen, nil)
	}
	if session.BookedCount >= session.MaxPlayers {
		return nil, util.NewAppError(409, constants.CodeConflict, constants.MsgSessionFull, nil)
	}
	if _, err := s.sessions.FindRegistration(ctx, sessionID, user.ID); err == nil {
		return nil, util.NewAppError(409, constants.CodeConflict, constants.MsgAlreadyRegistered, nil)
	}
	if err := s.sessions.Transaction(ctx, func(txCtx context.Context) error {
		if err := s.sessions.IncrementBooked(txCtx, sessionID); err != nil {
			return err
		}
		reg := &model.SessionRegistration{SessionID: sessionID, UserID: user.ID}
		if err := s.sessions.CreateRegistration(txCtx, reg); err != nil {
			return err
		}
		return nil
	}); err != nil {
		if errors.Is(err, util.ErrConflict) {
			return nil, util.NewAppError(409, constants.CodeConflict, constants.MsgSessionFull, nil)
		}
		return nil, util.WrapAppError(fmt.Errorf("game_session[id=%d] register: %w", sessionID, err), 500, constants.CodeInternalError, constants.MsgInternalError)
	}
	s.logger.Info(fmt.Sprintf(constants.LogGameSessionRegisterSuccess, sessionID, user.ID))
	updated, err := s.sessions.FindByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	if updated.BookedCount >= updated.MaxPlayers && updated.Status == constants.GameSessionStatusOpen {
		if err := s.sessions.UpdateStatusIf(ctx, sessionID, constants.GameSessionStatusOpen, constants.GameSessionStatusLocked); err != nil {
			if errors.Is(err, util.ErrConflict) {
				// A concurrent registration already locked the session; return latest state.
				updated, err = s.sessions.FindByID(ctx, sessionID)
				if err != nil {
					return nil, err
				}
				return updated, nil
			}
			s.logger.Error(fmt.Sprintf(constants.LogGameSessionLockFailed, sessionID, err))
			return nil, util.WrapAppError(fmt.Errorf("game_session[id=%d] auto lock: %w", sessionID, err), 500, constants.CodeInternalError, constants.MsgInternalError)
		}
		s.logger.Info(fmt.Sprintf(constants.LogGameSessionLockSuccess, sessionID, updated.BookedCount))
		updated.Status = constants.GameSessionStatusLocked
	}
	return updated, nil
}

// Quit removes a player's registration from an open session.
func (s *GameSessionService) Quit(ctx context.Context, userID, sessionID uint) (*model.GameSession, error) {
	session, err := s.sessions.FindByID(ctx, sessionID)
	if err != nil {
		return nil, util.WrapAppError(fmt.Errorf("game_session[id=%d] quit find: %w", sessionID, err), 404, constants.CodeNotFound, constants.MsgNotFound)
	}
	if session.Status != constants.GameSessionStatusOpen {
		return nil, util.NewAppError(409, constants.CodeConflict, constants.MsgSessionNotOpen, nil)
	}
	if _, err := s.sessions.FindRegistration(ctx, sessionID, userID); err != nil {
		return nil, util.NewAppError(409, constants.CodeConflict, constants.MsgNotRegistered, nil)
	}
	if err := s.sessions.Transaction(ctx, func(txCtx context.Context) error {
		if err := s.sessions.DeleteRegistration(txCtx, sessionID, userID); err != nil {
			return err
		}
		if err := s.sessions.DecrementBooked(txCtx, sessionID); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, util.WrapAppError(fmt.Errorf("game_session[id=%d] quit: %w", sessionID, err), 500, constants.CodeInternalError, constants.MsgInternalError)
	}
	updated, err := s.sessions.FindByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

// Lock manually locks an open session (admin).
func (s *GameSessionService) Lock(ctx context.Context, sessionID uint) (*model.GameSession, error) {
	session, err := s.sessions.FindByID(ctx, sessionID)
	if err != nil {
		return nil, util.WrapAppError(fmt.Errorf("game_session[id=%d] lock find: %w", sessionID, err), 404, constants.CodeNotFound, constants.MsgNotFound)
	}
	if session.Status != constants.GameSessionStatusOpen {
		return nil, util.NewAppError(409, constants.CodeConflict, constants.MsgSessionStatusInvalid, nil)
	}
	if err := s.sessions.UpdateStatusIf(ctx, sessionID, constants.GameSessionStatusOpen, constants.GameSessionStatusLocked); err != nil {
		if errors.Is(err, util.ErrConflict) {
			return nil, util.NewAppError(409, constants.CodeConflict, constants.MsgSessionStatusInvalid, nil)
		}
		s.logger.Error(fmt.Sprintf(constants.LogGameSessionLockFailed, sessionID, err))
		return nil, util.WrapAppError(fmt.Errorf("game_session[id=%d] lock: %w", sessionID, err), 500, constants.CodeInternalError, constants.MsgInternalError)
	}
	s.logger.Info(fmt.Sprintf(constants.LogGameSessionLockSuccess, sessionID, session.BookedCount))
	session.Status = constants.GameSessionStatusLocked
	return session, nil
}

// Cancel cancels an open/locked session (admin).
func (s *GameSessionService) Cancel(ctx context.Context, sessionID uint) (*model.GameSession, error) {
	session, err := s.sessions.FindByID(ctx, sessionID)
	if err != nil {
		return nil, util.WrapAppError(fmt.Errorf("game_session[id=%d] cancel find: %w", sessionID, err), 404, constants.CodeNotFound, constants.MsgNotFound)
	}
	if session.Status != constants.GameSessionStatusOpen && session.Status != constants.GameSessionStatusLocked {
		return nil, util.NewAppError(409, constants.CodeConflict, constants.MsgSessionStatusInvalid, nil)
	}
	if err := s.sessions.UpdateStatusIf(ctx, sessionID, session.Status, constants.GameSessionStatusCancelled); err != nil {
		if errors.Is(err, util.ErrConflict) {
			return nil, util.NewAppError(409, constants.CodeConflict, constants.MsgSessionStatusInvalid, nil)
		}
		return nil, util.WrapAppError(fmt.Errorf("game_session[id=%d] cancel: %w", sessionID, err), 500, constants.CodeInternalError, constants.MsgInternalError)
	}
	session.Status = constants.GameSessionStatusCancelled
	return session, nil
}

// Finish marks a locked session as finished (admin) and records revenue.
func (s *GameSessionService) Finish(ctx context.Context, sessionID uint, unitPrice float64, revenueRepo *repository.RevenueRepository) (*model.GameSession, error) {
	session, err := s.sessions.FindByID(ctx, sessionID)
	if err != nil {
		return nil, util.WrapAppError(fmt.Errorf("game_session[id=%d] finish find: %w", sessionID, err), 404, constants.CodeNotFound, constants.MsgNotFound)
	}
	if session.Status != constants.GameSessionStatusLocked {
		return nil, util.NewAppError(409, constants.CodeConflict, constants.MsgSessionStatusInvalid, nil)
	}
	rev := &model.Revenue{
		SessionID: sessionID, ThemeRoomID: session.ThemeRoomID,
		Amount: unitPrice * float64(session.BookedCount), BookedCount: session.BookedCount,
		Date: time.Now(), CreatedAt: time.Now(),
	}
	if err := s.sessions.Transaction(ctx, func(txCtx context.Context) error {
		if err := s.sessions.UpdateStatusIf(txCtx, sessionID, constants.GameSessionStatusLocked, constants.GameSessionStatusFinished); err != nil {
			return err
		}
		if err := revenueRepo.Create(txCtx, rev); err != nil {
			return err
		}
		return nil
	}); err != nil {
		if errors.Is(err, util.ErrConflict) {
			return nil, util.NewAppError(409, constants.CodeConflict, constants.MsgSessionStatusInvalid, nil)
		}
		return nil, util.WrapAppError(fmt.Errorf("game_session[id=%d] finish: %w", sessionID, err), 500, constants.CodeInternalError, constants.MsgInternalError)
	}
	s.logger.Info(fmt.Sprintf(constants.LogGameSessionFinishSuccess, sessionID))
	s.logger.Info(fmt.Sprintf(constants.LogRevenueRecordSuccess, sessionID, rev.Amount))
	session.Status = constants.GameSessionStatusFinished
	return session, nil
}
