package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNaturalTime(t *testing.T) {
	base := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		value    time.Time
		expected string
	}{
		// Past times
		{
			name:     "a few seconds ago",
			value:    base.Add(-30 * time.Second),
			expected: "a few seconds ago",
		},
		{
			name:     "exactly 60 seconds ago",
			value:    base.Add(-60 * time.Second),
			expected: "a few seconds ago",
		},
		{
			name:     "a few minutes ago",
			value:    base.Add(-3 * time.Minute),
			expected: "a few minutes ago",
		},
		{
			name:     "10 minutes ago",
			value:    base.Add(-10 * time.Minute),
			expected: "10 minutes ago",
		},
		{
			name:     "59 minutes ago",
			value:    base.Add(-59 * time.Minute),
			expected: "59 minutes ago",
		},
		{
			name:     "2 hours ago (same day)",
			value:    base.Add(-2 * time.Hour),
			expected: "2 hours ago",
		},
		{
			name:     "yesterday",
			value:    time.Date(2024, 1, 14, 12, 0, 0, 0, time.UTC),
			expected: "yesterday",
		},
		{
			name:     "day before yesterday",
			value:    time.Date(2024, 1, 13, 12, 0, 0, 0, time.UTC),
			expected: "day before yesterday",
		},
		{
			name:     "5 days ago",
			value:    base.Add(-5 * 24 * time.Hour),
			expected: "5 days ago",
		},
		{
			name:     "29 days ago",
			value:    base.Add(-29 * 24 * time.Hour),
			expected: "29 days ago",
		},
		{
			name:     "last month",
			value:    base.Add(-30 * 24 * time.Hour),
			expected: "last month",
		},
		{
			name:     "3 months ago",
			value:    base.Add(-90 * 24 * time.Hour),
			expected: "3 months ago",
		},
		{
			name:     "last year",
			value:    base.Add(-365 * 24 * time.Hour),
			expected: "last year",
		},
		{
			name:     "3 years ago",
			value:    base.Add(-3 * 365 * 24 * time.Hour),
			expected: "3 years ago",
		},

		// Future times
		{
			name:     "in a few seconds",
			value:    base.Add(30 * time.Second),
			expected: "in a few seconds",
		},
		{
			name:     "in a few minutes",
			value:    base.Add(3 * time.Minute),
			expected: "in a few minutes",
		},
		{
			name:     "in 10 minutes",
			value:    base.Add(10 * time.Minute),
			expected: "in 10 minutes",
		},
		{
			name:     "in 2 hours",
			value:    base.Add(2 * time.Hour),
			expected: "in 2 hours",
		},
		{
			name:     "tomorrow",
			value:    base.Add(24 * time.Hour),
			expected: "tomorrow",
		},
		{
			name:     "day after tomorrow",
			value:    base.Add(48 * time.Hour),
			expected: "day after tomorrow",
		},
		{
			name:     "in 5 days",
			value:    base.Add(5 * 24 * time.Hour),
			expected: "in 5 days",
		},
		{
			name:     "next month",
			value:    base.Add(30 * 24 * time.Hour),
			expected: "next month",
		},
		{
			name:     "in 3 months",
			value:    base.Add(90 * 24 * time.Hour),
			expected: "in 3 months",
		},
		{
			name:     "next year",
			value:    base.Add(365 * 24 * time.Hour),
			expected: "next year",
		},
		{
			name:     "in 3 years",
			value:    base.Add(3 * 365 * 24 * time.Hour),
			expected: "in 3 years",
		},
		{
			name:     "same time (present)",
			value:    base,
			expected: "in a few seconds",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NaturalTime(base, tt.value)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPastNaturalTime(t *testing.T) {
	base := time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC)

	tests := []struct {
		name     string
		value    time.Time
		expected string
	}{
		{
			name:     "30 seconds ago",
			value:    base.Add(-30 * time.Second),
			expected: "a few seconds ago",
		},
		{
			name:     "1 minute ago",
			value:    base.Add(-1 * time.Minute),
			expected: "a few seconds ago",
		},
		{
			name:     "2 minutes ago",
			value:    base.Add(-2 * time.Minute),
			expected: "a few minutes ago",
		},
		{
			name:     "30 minutes ago",
			value:    base.Add(-30 * time.Minute),
			expected: "30 minutes ago",
		},
		{
			name:     "5 hours ago (same day)",
			value:    base.Add(-5 * time.Hour),
			expected: "5 hours ago",
		},
		{
			name:     "yesterday at same time",
			value:    time.Date(2024, 1, 14, 14, 30, 0, 0, time.UTC),
			expected: "yesterday",
		},
		{
			name:     "day before yesterday at same time",
			value:    time.Date(2024, 1, 13, 14, 30, 0, 0, time.UTC),
			expected: "day before yesterday",
		},
		{
			name:     "10 days ago",
			value:    base.Add(-10 * 24 * time.Hour),
			expected: "10 days ago",
		},
		{
			name:     "exactly 30 days ago",
			value:    base.Add(-30 * 24 * time.Hour),
			expected: "last month",
		},
		{
			name:     "6 months ago",
			value:    base.Add(-180 * 24 * time.Hour),
			expected: "6 months ago",
		},
		{
			name:     "exactly 12 months ago",
			value:    base.Add(-360 * 24 * time.Hour),
			expected: "last year",
		},
		{
			name:     "2 years ago",
			value:    base.Add(-2 * 365 * 24 * time.Hour),
			expected: "2 years ago",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := pastNaturalTime(base, tt.value)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFutureNaturalTime(t *testing.T) {
	base := time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC)

	tests := []struct {
		name     string
		value    time.Time
		expected string
	}{
		{
			name:     "in 30 seconds",
			value:    base.Add(30 * time.Second),
			expected: "in a few seconds",
		},
		{
			name:     "in 2 minutes",
			value:    base.Add(2 * time.Minute),
			expected: "in a few minutes",
		},
		{
			name:     "in 30 minutes",
			value:    base.Add(30 * time.Minute),
			expected: "in 30 minutes",
		},
		{
			name:     "in 5 hours",
			value:    base.Add(5 * time.Hour),
			expected: "in 5 hours",
		},
		{
			name:     "in exactly 24 hours (tomorrow)",
			value:    base.Add(24 * time.Hour),
			expected: "tomorrow",
		},
		{
			name:     "in 48 hours (day after tomorrow)",
			value:    base.Add(48 * time.Hour),
			expected: "day after tomorrow",
		},
		{
			name:     "in 10 days",
			value:    base.Add(10 * 24 * time.Hour),
			expected: "in 10 days",
		},
		{
			name:     "in exactly 30 days (next month)",
			value:    base.Add(30 * 24 * time.Hour),
			expected: "next month",
		},
		{
			name:     "in 6 months",
			value:    base.Add(180 * 24 * time.Hour),
			expected: "in 6 months",
		},
		{
			name:     "in 12 months (next year)",
			value:    base.Add(360 * 24 * time.Hour),
			expected: "next year",
		},
		{
			name:     "in 2 years",
			value:    base.Add(2 * 365 * 24 * time.Hour),
			expected: "in 2 years",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := futureNaturalTime(base, tt.value)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Edge case tests
func TestNaturalTimeEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		base     time.Time
		value    time.Time
		expected string
	}{
		{
			name:     "exactly at midnight boundary",
			base:     time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			value:    time.Date(2024, 1, 14, 23, 59, 59, 0, time.UTC),
			expected: "a few seconds ago",
		},
		{
			name:     "leap year consideration",
			base:     time.Date(2024, 3, 1, 12, 0, 0, 0, time.UTC),
			value:    time.Date(2024, 2, 29, 12, 0, 0, 0, time.UTC),
			expected: "yesterday",
		},
		{
			name:     "year boundary",
			base:     time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			value:    time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC),
			expected: "a few seconds ago",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NaturalTime(tt.base, tt.value)
			assert.Equal(t, tt.expected, result)
		})
	}
}
