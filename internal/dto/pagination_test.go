package dto

import "testing"

func TestPageQueryNormalizeDefaults(t *testing.T) {
	q := &PageQuery{Page: 0, PageSize: 0}
	q.Normalize()
	if q.Page != 1 {
		t.Fatalf("Page = %d, want 1", q.Page)
	}
	if q.PageSize != 10 {
		t.Fatalf("PageSize = %d, want 10", q.PageSize)
	}
}

func TestPageQueryNormalizeClampsLargeSize(t *testing.T) {
	q := &PageQuery{Page: -1, PageSize: 200}
	q.Normalize()
	if q.Page != 1 {
		t.Fatalf("Page = %d, want 1", q.Page)
	}
	if q.PageSize != 10 {
		t.Fatalf("PageSize = %d, want 10", q.PageSize)
	}
}
