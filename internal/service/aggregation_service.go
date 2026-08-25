package service

import (
	"math"
	"sort"
	"time"

	"metricscollector/internal/model"
	"metricscollector/pkg/idgen"
)

func (s *Service) CreateAggregation(input model.Aggregation) (*model.Aggregation, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	a := &model.Aggregation{
		ID:        idgen.Hex(),
		MetricID:  input.MetricID,
		WindowSec: input.WindowSec,
		Func:      input.Func,
		Result:    input.Result,
		Timestamp: input.Timestamp.UTC(),
		CreatedAt: time.Now().UTC(),
	}
	if err := s.store.CreateAggregation(a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *Service) ListAggregations(filter model.AggregationFilter, page, size int) ([]*model.Aggregation, int, error) {
	all := s.store.ListAggregations()
	matched := make([]*model.Aggregation, 0, len(all))
	for _, a := range all {
		if filter.Match(a) {
			matched = append(matched, a)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].Timestamp.After(matched[j].Timestamp)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Aggregation{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) GetAggregation(id string) (*model.Aggregation, error) {
	return s.store.GetAggregation(id)
}

func (s *Service) ComputeAggregation(metricID string, windowSec int, fn string) (*model.Aggregation, error) {
	if windowSec <= 0 {
		return nil, model.NewValidationError("window_sec", "窗口秒数必须大于 0")
	}
	if fn != model.AggregationFuncSum && fn != model.AggregationFuncAvg && fn != model.AggregationFuncMax &&
		fn != model.AggregationFuncMin && fn != model.AggregationFuncCount {
		return nil, model.NewValidationError("func", "聚合函数不合法，应为 sum/avg/max/min/count")
	}
	now := time.Now().UTC()
	windowStart := now.Add(-time.Duration(windowSec) * time.Second)

	samples := s.store.ListSamplesByMetricID(metricID)
	var values []float64
	for _, sa := range samples {
		if sa.Timestamp.After(windowStart) || sa.Timestamp.Equal(windowStart) {
			values = append(values, sa.Value)
		}
	}
	if len(values) == 0 {
		return nil, model.NewValidationError("window", "窗口内无样本数据")
	}

	var result float64
	switch fn {
	case model.AggregationFuncSum:
		for _, v := range values {
			result += v
		}
	case model.AggregationFuncAvg:
		for _, v := range values {
			result += v
		}
		result = result / float64(len(values))
	case model.AggregationFuncMax:
		result = values[0]
		for _, v := range values[1:] {
			if v > result {
				result = v
			}
		}
	case model.AggregationFuncMin:
		result = values[0]
		for _, v := range values[1:] {
			if v < result {
				result = v
			}
		}
	case model.AggregationFuncCount:
		result = float64(len(values))
	}

	agg := &model.Aggregation{
		ID:        idgen.Hex(),
		MetricID:  metricID,
		WindowSec: windowSec,
		Func:      fn,
		Result:    math.Round(result*1e6) / 1e6,
		Timestamp: now,
		CreatedAt: now,
	}
	if err := s.store.CreateAggregation(agg); err != nil {
		return nil, err
	}
	return agg, nil
}
