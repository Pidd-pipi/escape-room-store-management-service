package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/lp/escape-room-ops/internal/constants"
	"github.com/lp/escape-room-ops/internal/dto"
	"github.com/lp/escape-room-ops/internal/model"
	"github.com/lp/escape-room-ops/internal/repository"
	"github.com/lp/escape-room-ops/internal/util"
)

// RevenuePoint is one aggregated analytics point.
type RevenuePoint struct {
	Date        string  `json:"date"`
	Amount      float64 `json:"amount"`
	BookedCount int     `json:"booked_count"`
}

// RevenueSummary is the analytics payload.
type RevenueSummary struct {
	Range       string         `json:"range"`
	TotalAmount float64        `json:"total_amount"`
	TotalBooked int            `json:"total_booked"`
	Occupancy   float64        `json:"occupancy"`
	Points      []RevenuePoint `json:"points"`
	RepeatRate  float64        `json:"repeat_rate"`
	PeakHour    int            `json:"peak_hour"`
}

// RevenueService provides revenue and occupancy analytics.
type RevenueService struct {
	revenue  *repository.RevenueRepository
	sessions *repository.GameSessionRepository
	logger   *slog.Logger
}

// NewRevenueService wires the revenue service dependencies.
func NewRevenueService(revenue *repository.RevenueRepository, sessions *repository.GameSessionRepository, logger *slog.Logger) *RevenueService {
	return &RevenueService{revenue: revenue, sessions: sessions, logger: logger}
}

// Analyze returns revenue analytics for the given range (day/week/month).
func (s *RevenueService) Analyze(ctx context.Context, q *dto.RevenueQuery) (*RevenueSummary, error) {
	rangeName := q.Range
	if rangeName == "" {
		rangeName = "month"
	}
	now := time.Now()
	var since time.Time
	switch rangeName {
	case "day":
		since = now.Add(-24 * time.Hour)
	case "week":
		since = now.AddDate(0, 0, -7)
	default:
		since = now.AddDate(0, -1, 0)
	}
	records, err := s.revenue.ListSince(ctx, since)
	if err != nil {
		return nil, util.WrapAppError(fmt.Errorf("revenue analyze list: %w", err), 500, constants.CodeInternalError, constants.MsgInternalError)
	}
	summary := &RevenueSummary{Range: rangeName, Points: []RevenuePoint{}}
	byDay := map[string]*RevenuePoint{}
	peakHourCounts := map[int]int{}
	sessionIDs := map[uint]bool{}
	for _, r := range records {
		summary.TotalAmount += r.Amount
		summary.TotalBooked += r.BookedCount
		sessionIDs[r.SessionID] = true
		key := r.Date.Format("2006-01-02")
		pt, ok := byDay[key]
		if !ok {
			pt = &RevenuePoint{Date: key}
			byDay[key] = pt
		}
		pt.Amount += r.Amount
		pt.BookedCount += r.BookedCount
		peakHourCounts[r.Date.Hour()]++
	}
	for _, pt := range byDay {
		pt.Amount = round2(pt.Amount)
		summary.Points = append(summary.Points, *pt)
	}
	summary.TotalAmount = round2(summary.TotalAmount)
	if len(sessionIDs) > 0 {
		ids := make([]uint, 0, len(sessionIDs))
		for id := range sessionIDs {
			ids = append(ids, id)
		}
		var sessions []model.GameSession
		if err := s.sessions.QueryByIDs(ctx, ids, &sessions); err == nil {
			totalCapacity := 0
			for _, sess := range sessions {
				totalCapacity += sess.MaxPlayers
			}
			if totalCapacity > 0 {
				summary.Occupancy = round2(float64(summary.TotalBooked) / float64(totalCapacity))
			}
		}
	}
	peak := 0
	max := 0
	for h, c := range peakHourCounts {
		if c > max {
			max = c
			peak = h
		}
	}
	summary.PeakHour = peak
	s.logger.Info(fmt.Sprintf(constants.LogRevenueExportSuccess, len(records)))
	return summary, nil
}

// Export returns raw revenue rows for CSV-like export.
func (s *RevenueService) Export(ctx context.Context, q *dto.RevenueQuery) ([]model.Revenue, error) {
	if q.Range == "" || q.Range == "all" {
		return s.revenue.ListAll(ctx)
	}
	return s.revenue.ListSince(ctx, time.Now().AddDate(0, -1, 0))
}

func round2(v float64) float64 {
	return float64(int(v*100+0.5)) / 100
}
