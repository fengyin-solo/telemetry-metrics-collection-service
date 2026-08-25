package store

import "metricscollector/internal/model"

func (s *MemoryStore) CreateExportJob(e *model.ExportJob) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.exportJobs[e.ID] = e
	return nil
}

func (s *MemoryStore) GetExportJob(id string) (*model.ExportJob, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.exportJobs[id]
	if !ok {
		return nil, ErrNotFound
	}
	return e, nil
}

func (s *MemoryStore) ListExportJobs() []*model.ExportJob {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.ExportJob, 0, len(s.exportJobs))
	for _, e := range s.exportJobs {
		list = append(list, e)
	}
	return list
}

func (s *MemoryStore) UpdateExportJob(e *model.ExportJob) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.exportJobs[e.ID]; !ok {
		return ErrNotFound
	}
	s.exportJobs[e.ID] = e
	return nil
}

func (s *MemoryStore) DeleteExportJob(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.exportJobs[id]; !ok {
		return ErrNotFound
	}
	delete(s.exportJobs, id)
	return nil
}
