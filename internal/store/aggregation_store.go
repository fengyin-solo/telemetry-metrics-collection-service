package store

import "metricscollector/internal/model"

func (s *MemoryStore) CreateAggregation(a *model.Aggregation) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.aggregations[a.ID] = a
	return nil
}

func (s *MemoryStore) GetAggregation(id string) (*model.Aggregation, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	a, ok := s.aggregations[id]
	if !ok {
		return nil, ErrNotFound
	}
	return a, nil
}

func (s *MemoryStore) ListAggregations() []*model.Aggregation {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Aggregation, 0, len(s.aggregations))
	for _, a := range s.aggregations {
		list = append(list, a)
	}
	return list
}

func (s *MemoryStore) ListAggregationsByMetricIDFunc(metricID, fn string) []*model.Aggregation {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Aggregation, 0)
	for _, a := range s.aggregations {
		if a.MetricID == metricID && (fn == "" || a.Func == fn) {
			list = append(list, a)
		}
	}
	return list
}
