package util

import "math"

// EscapeRateResult is the aggregated leaderboard payload.
type EscapeRateResult struct {
	AverageMinutes float64 `json:"average_minutes"`
	TotalPlays     int     `json:"total_plays"`
	Escaped        int     `json:"escaped"`
	EscapeRate     float64 `json:"escape_rate"`
}

// CalcEscapeRate returns the escape success rate (0-1).
func CalcEscapeRate(escaped, total int) float64 {
	if total <= 0 {
		return 0
	}
	return round2(float64(escaped) / float64(total))
}

// CalcAverageMinutes returns the average clear time.
func CalcAverageMinutes(totalMinutes, plays int) float64 {
	if plays <= 0 {
		return 0
	}
	return round2(float64(totalMinutes) / float64(plays))
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}
