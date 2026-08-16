package util

import "testing"

func TestCalcEscapeRateFraction(t *testing.T) {
	if got := CalcEscapeRate(3, 4); got != 0.75 {
		t.Fatalf("CalcEscapeRate(3,4) = %v, want 0.75", got)
	}
}

func TestCalcAverageMinutesFraction(t *testing.T) {
	if got := CalcAverageMinutes(100, 4); got != 25 {
		t.Fatalf("CalcAverageMinutes(100,4) = %v, want 25", got)
	}
}
