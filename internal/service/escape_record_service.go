package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/lp/escape-room-ops/internal/constants"
	"github.com/lp/escape-room-ops/internal/dto"
	"github.com/lp/escape-room-ops/internal/model"
	"github.com/lp/escape-room-ops/internal/repository"
	"github.com/lp/escape-room-ops/internal/util"
)

// RankEntry is one leaderboard row per theme.
type RankEntry struct {
	ThemeRoomID    uint    `json:"theme_room_id"`
	ThemeName      string  `json:"theme_name"`
	FastestMinutes int     `json:"fastest_minutes"`
	TotalPlays     int     `json:"total_plays"`
	Escaped        int     `json:"escaped"`
	EscapeRate     float64 `json:"escape_rate"`
}

// EscapeRecordService records outcomes and builds the leaderboard.
type EscapeRecordService struct {
	records *repository.EscapeRecordRepository
	themes  *repository.ThemeRoomRepository
	sessions *repository.GameSessionRepository
	logger  *slog.Logger
}

// NewEscapeRecordService wires the escape record service dependencies.
func NewEscapeRecordService(records *repository.EscapeRecordRepository, themes *repository.ThemeRoomRepository, sessions *repository.GameSessionRepository, logger *slog.Logger) *EscapeRecordService {
	return &EscapeRecordService{records: records, themes: themes, sessions: sessions, logger: logger}
}

// Create stores a team outcome.
func (s *EscapeRecordService) Create(ctx context.Context, req *dto.CreateEscapeRecordRequest) (*model.EscapeRecord, error) {
	if _, err := s.sessions.FindByID(ctx, req.SessionID); err != nil {
		return nil, util.WrapAppError(fmt.Errorf("escape_record[session=%d] session lookup: %w", req.SessionID, err), 404, constants.CodeNotFound, constants.MsgNotFound)
	}
	if _, err := s.themes.FindByID(ctx, req.ThemeRoomID); err != nil {
		return nil, util.WrapAppError(fmt.Errorf("escape_record[theme=%d] theme lookup: %w", req.ThemeRoomID, err), 404, constants.CodeNotFound, constants.MsgNotFound)
	}
	e := &model.EscapeRecord{
		SessionID: req.SessionID, ThemeRoomID: req.ThemeRoomID, TeamName: req.TeamName,
		DurationMinutes: req.DurationMinutes, HintCount: req.HintCount, Escaped: req.Escaped,
	}
	if err := s.records.Create(ctx, e); err != nil {
		s.logger.Error(fmt.Sprintf(constants.LogEscapeRecordCreateFailed, req.SessionID, err))
		return nil, util.WrapAppError(fmt.Errorf("escape_record[session=%d] create: %w", req.SessionID, err), 500, constants.CodeInternalError, constants.MsgInternalError)
	}
	s.logger.Info(fmt.Sprintf(constants.LogEscapeRecordCreateSuccess, e.ID, req.SessionID))
	return e, nil
}

// Leaderboard aggregates per-theme fastest time and escape rate.
func (s *EscapeRecordService) Leaderboard(ctx context.Context) ([]RankEntry, error) {
	records, err := s.records.ListAll(ctx)
	if err != nil {
		return nil, util.WrapAppError(fmt.Errorf("escape_record leaderboard list: %w", err), 500, constants.CodeInternalError, constants.MsgInternalError)
	}
	stats := map[uint]*RankEntry{}
	for _, r := range records {
		entry, ok := stats[r.ThemeRoomID]
		if !ok {
			entry = &RankEntry{ThemeRoomID: r.ThemeRoomID, FastestMinutes: r.DurationMinutes}
			stats[r.ThemeRoomID] = entry
		}
		if r.DurationMinutes < entry.FastestMinutes {
			entry.FastestMinutes = r.DurationMinutes
		}
		entry.TotalPlays++
		if r.Escaped {
			entry.Escaped++
		}
	}
	for id, entry := range stats {
		if t, err := s.themes.FindByID(ctx, id); err == nil {
			entry.ThemeName = t.Name
		}
		entry.EscapeRate = util.CalcEscapeRate(entry.Escaped, entry.TotalPlays)
	}
	result := make([]RankEntry, 0, len(stats))
	for _, e := range stats {
		result = append(result, *e)
	}
	s.logger.Info(fmt.Sprintf(constants.LogLeaderboardBuildSuccess, len(records)))
	return result, nil
}
