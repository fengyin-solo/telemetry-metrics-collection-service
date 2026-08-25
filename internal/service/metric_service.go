package service

import (
	"sort"
	"time"

	"metricscollector/internal/model"
	"metricscollector/pkg/idgen"
)

func (s *Service) CreateMetric(input model.Metric) (*model.Metric, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	m := &model.Metric{
		ID:          idgen.Hex(),
		Name:        input.Name,
		Type:        input.Type,
		Unit:        input.Unit,
		Description: input.Description,
		Status:      input.Status,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
	if err := s.store.CreateMetric(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *Service) ListMetrics(filter model.MetricFilter, page, size int) ([]*model.Metric, int, error) {
	all := s.store.ListMetrics()
	matched := make([]*model.Metric, 0, len(all))
	for _, m := range all {
		if filter.Match(m) {
			matched = append(matched, m)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Metric{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) GetMetric(id string) (*model.Metric, error) {
	return s.store.GetMetric(id)
}

func (s *Service) UpdateMetric(id string, input model.Metric) (*model.Metric, error) {
	m, err := s.store.GetMetric(id)
	if err != nil {
		return nil, err
	}
	m.Name = input.Name
	m.Type = input.Type
	m.Unit = input.Unit
	m.Description = input.Description
	m.Status = input.Status
	m.UpdatedAt = time.Now().UTC()
	if err := m.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateMetric(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *Service) DeleteMetric(id string) error {
	return s.store.DeleteMetric(id)
}
