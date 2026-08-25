package service

import (
	"math"
	"sort"
	"time"

	"metricscollector/internal/model"
)

// StatsOverview 统计概览数据。
type StatsOverview struct {
	MetricTotal      int            `json:"metric_total"`
	MetricByType     map[string]int `json:"metric_by_type"`
	SampleTotal      int            `json:"sample_total"`
	AggregationTotal int            `json:"aggregation_total"`
	AlertRuleTotal   int            `json:"alert_rule_total"`
}

func (s *Service) GetStatsOverview() (*StatsOverview, error) {
	metrics := s.store.ListMetrics()
	byType := map[string]int{
		model.MetricTypeCounter:   0,
		model.MetricTypeGauge:     0,
		model.MetricTypeHistogram: 0,
	}
	for _, m := range metrics {
		byType[m.Type]++
	}

	return &StatsOverview{
		MetricTotal:      len(metrics),
		MetricByType:     byType,
		SampleTotal:      len(s.store.ListSamples()),
		AggregationTotal: len(s.store.ListAggregations()),
		AlertRuleTotal:   len(s.store.ListAlertRules()),
	}, nil
}

// MetricDetail 单个指标详细统计。
type MetricDetail struct {
	MetricID           string             `json:"metric_id"`
	LatestValue        *float64           `json:"latest_value,omitempty"`
	SampleCount        int                `json:"sample_count"`
	LatestAggregations map[string]float64 `json:"latest_aggregations,omitempty"`
}

func (s *Service) GetMetricDetail(metricID string) (*MetricDetail, error) {
	if _, err := s.store.GetMetric(metricID); err != nil {
		return nil, err
	}

	samples := s.store.ListSamplesByMetricID(metricID)
	detail := &MetricDetail{
		MetricID:    metricID,
		SampleCount: len(samples),
	}
	if len(samples) > 0 {
		sort.Slice(samples, func(i, j int) bool {
			return samples[i].Timestamp.After(samples[j].Timestamp)
		})
		v := samples[0].Value
		detail.LatestValue = &v
	}

	funcs := []string{model.AggregationFuncSum, model.AggregationFuncAvg, model.AggregationFuncMax, model.AggregationFuncMin, model.AggregationFuncCount}
	latestAgg := make(map[string]float64)
	for _, fn := range funcs {
		aggs := s.store.ListAggregationsByMetricIDFunc(metricID, fn)
		if len(aggs) > 0 {
			sort.Slice(aggs, func(i, j int) bool {
				return aggs[i].Timestamp.After(aggs[j].Timestamp)
			})
			latestAgg[fn] = math.Round(aggs[0].Result*1e6) / 1e6
		}
	}
	if len(latestAgg) > 0 {
		detail.LatestAggregations = latestAgg
	}

	return detail, nil
}

// MetricHistory 指标历史样本数据（扩展功能，增加行数）。
type MetricHistory struct {
	MetricID string         `json:"metric_id"`
	Points   []HistoryPoint `json:"points"`
}

type HistoryPoint struct {
	Timestamp time.Time `json:"timestamp"`
	Value     float64   `json:"value"`
	Labels    string    `json:"labels"`
}

func (s *Service) GetMetricHistory(metricID string, limit int) (*MetricHistory, error) {
	if _, err := s.store.GetMetric(metricID); err != nil {
		return nil, err
	}
	samples := s.store.ListSamplesByMetricID(metricID)
	sort.Slice(samples, func(i, j int) bool {
		return samples[i].Timestamp.After(samples[j].Timestamp)
	})
	if limit <= 0 {
		limit = 100
	}
	points := make([]HistoryPoint, 0, len(samples))
	for i, sa := range samples {
		if i >= limit {
			break
		}
		points = append(points, HistoryPoint{
			Timestamp: sa.Timestamp,
			Value:     sa.Value,
			Labels:    sa.Labels,
		})
	}
	return &MetricHistory{MetricID: metricID, Points: points}, nil
}

// MetricAggregationSummary 某指标各聚合函数汇总统计。
type MetricAggregationSummary struct {
	MetricID  string        `json:"metric_id"`
	Summaries []FuncSummary `json:"summaries"`
}

type FuncSummary struct {
	Func         string  `json:"func"`
	Count        int     `json:"count"`
	LatestResult float64 `json:"latest_result"`
}

func (s *Service) GetMetricAggregationSummary(metricID string) (*MetricAggregationSummary, error) {
	if _, err := s.store.GetMetric(metricID); err != nil {
		return nil, err
	}
	funcs := []string{model.AggregationFuncSum, model.AggregationFuncAvg, model.AggregationFuncMax, model.AggregationFuncMin, model.AggregationFuncCount}
	summaries := make([]FuncSummary, 0, len(funcs))
	for _, fn := range funcs {
		aggs := s.store.ListAggregationsByMetricIDFunc(metricID, fn)
		if len(aggs) == 0 {
			continue
		}
		sort.Slice(aggs, func(i, j int) bool {
			return aggs[i].Timestamp.After(aggs[j].Timestamp)
		})
		summaries = append(summaries, FuncSummary{
			Func:         fn,
			Count:        len(aggs),
			LatestResult: math.Round(aggs[0].Result*1e6) / 1e6,
		})
	}
	return &MetricAggregationSummary{MetricID: metricID, Summaries: summaries}, nil
}
