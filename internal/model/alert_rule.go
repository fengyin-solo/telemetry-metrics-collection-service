package model

import (
	"strings"
	"time"
)

const (
	AlertOperatorGT  = "gt"
	AlertOperatorLT  = "lt"
	AlertOperatorGTE = "gte"
	AlertOperatorLTE = "lte"

	AlertLevelWarn     = "warn"
	AlertLevelCritical = "critical"

	AlertRuleStatusActive   = "active"
	AlertRuleStatusInactive = "inactive"
)

type AlertRule struct {
	ID        string    `json:"id"`
	MetricID  string    `json:"metric_id"`
	Threshold float64   `json:"threshold"`
	Operator  string    `json:"operator"`
	Level     string    `json:"level"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (a *AlertRule) Validate() error {
	if strings.TrimSpace(a.MetricID) == "" {
		return NewValidationError("metric_id", "指标 ID 不能为空")
	}
	if a.Operator == "" {
		a.Operator = AlertOperatorGT
	}
	if a.Operator != AlertOperatorGT && a.Operator != AlertOperatorLT &&
		a.Operator != AlertOperatorGTE && a.Operator != AlertOperatorLTE {
		return NewValidationError("operator", "操作符不合法，应为 gt/lt/gte/lte")
	}
	if a.Level == "" {
		a.Level = AlertLevelWarn
	}
	if a.Level != AlertLevelWarn && a.Level != AlertLevelCritical {
		return NewValidationError("level", "告警级别不合法，应为 warn/critical")
	}
	if a.Status == "" {
		a.Status = AlertRuleStatusActive
	}
	if a.Status != AlertRuleStatusActive && a.Status != AlertRuleStatusInactive {
		return NewValidationError("status", "规则状态不合法，应为 active/inactive")
	}
	return nil
}

type AlertRuleFilter struct {
	MetricID string
	Level    string
	Status   string
}

func (f AlertRuleFilter) Match(a *AlertRule) bool {
	if f.MetricID != "" && a.MetricID != f.MetricID {
		return false
	}
	if f.Level != "" && a.Level != f.Level {
		return false
	}
	if f.Status != "" && a.Status != f.Status {
		return false
	}
	return true
}

func (a *AlertRule) Evaluate(value float64) bool {
	switch a.Operator {
	case AlertOperatorGT:
		return value > a.Threshold
	case AlertOperatorLT:
		return value < a.Threshold
	case AlertOperatorGTE:
		return value >= a.Threshold
	case AlertOperatorLTE:
		return value <= a.Threshold
	}
	return false
}
