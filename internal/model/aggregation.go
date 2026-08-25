package model

import (
	"strings"
	"time"
)

const (
	AggregationFuncSum   = "sum"
	AggregationFuncAvg   = "avg"
	AggregationFuncMax   = "max"
	AggregationFuncMin   = "min"
	AggregationFuncCount = "count"
)

type Aggregation struct {
	ID        string    `json:"id"`
	MetricID  string    `json:"metric_id"`
	WindowSec int       `json:"window_sec"`
	Func      string    `json:"func"`
	Result    float64   `json:"result"`
	Timestamp time.Time `json:"timestamp"`
	CreatedAt time.Time `json:"created_at"`
}

func (a *Aggregation) Validate() error {
	if strings.TrimSpace(a.MetricID) == "" {
		return NewValidationError("metric_id", "指标 ID 不能为空")
	}
	if a.WindowSec <= 0 {
		return NewValidationError("window_sec", "窗口秒数必须大于 0")
	}
	if a.Func == "" {
		return NewValidationError("func", "聚合函数不能为空")
	}
	if a.Func != AggregationFuncSum && a.Func != AggregationFuncAvg && a.Func != AggregationFuncMax &&
		a.Func != AggregationFuncMin && a.Func != AggregationFuncCount {
		return NewValidationError("func", "聚合函数不合法，应为 sum/avg/max/min/count")
	}
	if a.Timestamp.IsZero() {
		return NewValidationError("timestamp", "时间戳不能为空")
	}
	return nil
}

type AggregationFilter struct {
	MetricID string
	Func     string
}

func (f AggregationFilter) Match(a *Aggregation) bool {
	if f.MetricID != "" && a.MetricID != f.MetricID {
		return false
	}
	if f.Func != "" && a.Func != f.Func {
		return false
	}
	return true
}
