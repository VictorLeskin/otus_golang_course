package storage

import (
	"testing"
	"time"
)

func TestEvent_IsStartTimeInInterval(t *testing.T) {
	base := time.Date(2026, 3, 1, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name  string
		event *Event
		from  time.Time
		to    time.Time
		want  bool
	}{
		{
			name:  "nil event",
			event: nil,
			from:  base,
			to:    base,
			want:  false,
		},
		{
			name:  "точно в интервале",
			event: &Event{StartTime: base},
			from:  base.Add(-1 * time.Hour),
			to:    base.Add(1 * time.Hour),
			want:  true,
		},
		{
			name:  "на левой границе",
			event: &Event{StartTime: base},
			from:  base,
			to:    base.Add(1 * time.Hour),
			want:  true,
		},
		{
			name:  "на правой границе",
			event: &Event{StartTime: base},
			from:  base.Add(-1 * time.Hour),
			to:    base,
			want:  true,
		},
		{
			name:  "до интервала",
			event: &Event{StartTime: base},
			from:  base.Add(1 * time.Hour),
			to:    base.Add(2 * time.Hour),
			want:  false,
		},
		{
			name:  "после интервала",
			event: &Event{StartTime: base},
			from:  base.Add(-2 * time.Hour),
			to:    base.Add(-1 * time.Hour),
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.event.IsStartTimeInInterval(tt.from, tt.to)
			if got != tt.want {
				t.Errorf("IsStartTimeInInterval(%v, %v) = %v; want %v",
					tt.from, tt.to, got, tt.want)
			}
		})
	}
}
