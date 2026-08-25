package model

import (
	"strings"
	"time"
)

type Sample struct {
	ID        string    `json:"id"`
	MetricID  string    `json:"metric_id"`
	Value     float64   `json:"value"`
	Labels    string    `json:"labels"`
	Timestamp time.Time `json:"timestamp"`
	CreatedAt time.Time `json:"created_at"`
}

func (s *Sample) Validate() error {
	if strings.TrimSpace(s.MetricID) == "" {
		return NewValidationError("metric_id", "指标 ID 不能为空")
	}
	if s.Timestamp.IsZero() {
		return NewValidationError("timestamp", "时间戳不能为空")
	}
	return nil
}

type SampleFilter struct {
	MetricID string
	Labels   string
}

func (f SampleFilter) Match(s *Sample) bool {
	if f.MetricID != "" && s.MetricID != f.MetricID {
		return false
	}
	if f.Labels != "" {
		want := strings.ToLower(strings.TrimSpace(f.Labels))
		if want != "" && !strings.Contains(strings.ToLower(s.Labels), want) {
			return false
		}
	}
	return true
}
