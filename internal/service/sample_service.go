package service

import (
	"sort"
	"time"

	"metricscollector/internal/model"
	"metricscollector/internal/store"
	"metricscollector/pkg/idgen"
)

func (s *Service) CreateSample(input model.Sample) (*model.Sample, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	m, err := s.store.GetMetric(input.MetricID)
	if err != nil {
		if err == store.ErrNotFound {
			return nil, model.NewValidationError("metric_id", "关联指标不存在")
		}
		return nil, err
	}
	if m.Status != model.MetricStatusActive {
		return nil, model.NewValidationError("metric_id", "关联指标未激活")
	}
	sa := &model.Sample{
		ID:        idgen.Hex(),
		MetricID:  input.MetricID,
		Value:     input.Value,
		Labels:    input.Labels,
		Timestamp: input.Timestamp.UTC(),
		CreatedAt: time.Now().UTC(),
	}
	if err := s.store.CreateSample(sa); err != nil {
		return nil, err
	}
	return sa, nil
}

func (s *Service) ListSamples(filter model.SampleFilter, page, size int) ([]*model.Sample, int, error) {
	all := s.store.ListSamples()
	matched := make([]*model.Sample, 0, len(all))
	for _, sa := range all {
		if filter.Match(sa) {
			matched = append(matched, sa)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].Timestamp.After(matched[j].Timestamp)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Sample{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) GetSample(id string) (*model.Sample, error) {
	return s.store.GetSample(id)
}
