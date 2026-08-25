package model

import (
	"strings"
	"time"
)

const (
	MetricTypeCounter   = "counter"
	MetricTypeGauge     = "gauge"
	MetricTypeHistogram = "histogram"

	MetricStatusActive   = "active"
	MetricStatusInactive = "inactive"
)

type Metric struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Unit        string    `json:"unit"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (m *Metric) Validate() error {
	m.Name = strings.TrimSpace(m.Name)
	if m.Name == "" {
		return NewValidationError("name", "指标名称不能为空")
	}
	if m.Type == "" {
		m.Type = MetricTypeGauge
	}
	if m.Type != MetricTypeCounter && m.Type != MetricTypeGauge && m.Type != MetricTypeHistogram {
		return NewValidationError("type", "指标类型不合法，应为 counter/gauge/histogram")
	}
	if m.Status == "" {
		m.Status = MetricStatusActive
	}
	if m.Status != MetricStatusActive && m.Status != MetricStatusInactive {
		return NewValidationError("status", "指标状态不合法，应为 active/inactive")
	}
	return nil
}

type MetricFilter struct {
	Type    string
	Status  string
	Keyword string
}

func (f MetricFilter) Match(m *Metric) bool {
	if f.Type != "" && m.Type != f.Type {
		return false
	}
	if f.Status != "" && m.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(m.Name), k) &&
			!strings.Contains(strings.ToLower(m.Description), k) {
			return false
		}
	}
	return true
}
