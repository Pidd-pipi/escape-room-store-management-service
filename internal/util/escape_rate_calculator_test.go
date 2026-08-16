package util

import "testing"

func TestCalcEscapeRate(t *testing.T) {
	tests := []struct {
		escaped int
		total   int
		want    float64
	}{
		{escaped: 4, total: 5, want: 0.8},
		{escaped: 0, total: 0, want: 0},
		{escaped: 2, total: 4, want: 0.5},
	}
	for _, tt := range tests {
		if got := CalcEscapeRate(tt.escaped, tt.total); got != tt.want {
			t.Fatalf("CalcEscapeRate(%d,%d) = %f, want %f", tt.escaped, tt.total, got, tt.want)
		}
	}
}

func TestCalcAverageMinutes(t *testing.T) {
	if got := CalcAverageMinutes(300, 5); got != 60 {
		t.Fatalf("CalcAverageMinutes = %f, want 60", got)
	}
	if got := CalcAverageMinutes(100, 0); got != 0 {
		t.Fatalf("CalcAverageMinutes zero plays = %f, want 0", got)
	}
}
