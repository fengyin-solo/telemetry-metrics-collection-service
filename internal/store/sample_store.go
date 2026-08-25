package store

import "metricscollector/internal/model"

func (s *MemoryStore) CreateSample(sa *model.Sample) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.samples[sa.ID] = sa
	return nil
}

func (s *MemoryStore) GetSample(id string) (*model.Sample, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	sa, ok := s.samples[id]
	if !ok {
		return nil, ErrNotFound
	}
	return sa, nil
}

func (s *MemoryStore) ListSamples() []*model.Sample {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Sample, 0, len(s.samples))
	for _, sa := range s.samples {
		list = append(list, sa)
	}
	return list
}

func (s *MemoryStore) ListSamplesByMetricID(metricID string) []*model.Sample {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Sample, 0)
	for _, sa := range s.samples {
		if sa.MetricID == metricID {
			list = append(list, sa)
		}
	}
	return list
}
