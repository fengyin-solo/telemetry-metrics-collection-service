package store

import (
	"strings"

	"metricscollector/internal/model"
)

func (s *MemoryStore) CreateMetric(m *model.Metric) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.metrics {
		if strings.EqualFold(exist.Name, m.Name) {
			return ErrConflict
		}
	}
	s.metrics[m.ID] = m
	return nil
}

func (s *MemoryStore) GetMetric(id string) (*model.Metric, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	m, ok := s.metrics[id]
	if !ok {
		return nil, ErrNotFound
	}
	return m, nil
}

func (s *MemoryStore) GetMetricByName(name string) (*model.Metric, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, m := range s.metrics {
		if strings.EqualFold(m.Name, name) {
			return m, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListMetrics() []*model.Metric {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Metric, 0, len(s.metrics))
	for _, m := range s.metrics {
		list = append(list, m)
	}
	return list
}

func (s *MemoryStore) UpdateMetric(m *model.Metric) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.metrics[m.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.metrics {
		if exist.ID != m.ID && strings.EqualFold(exist.Name, m.Name) {
			return ErrConflict
		}
	}
	s.metrics[m.ID] = m
	return nil
}

func (s *MemoryStore) DeleteMetric(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.metrics[id]; !ok {
		return ErrNotFound
	}
	delete(s.metrics, id)
	return nil
}
