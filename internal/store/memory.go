package store

import (
	"sync"

	"metricscollector/internal/model"
)

type MemoryStore struct {
	mu           sync.RWMutex
	metrics      map[string]*model.Metric
	samples      map[string]*model.Sample
	aggregations map[string]*model.Aggregation
	exportJobs   map[string]*model.ExportJob
	alertRules   map[string]*model.AlertRule
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		metrics:      make(map[string]*model.Metric),
		samples:      make(map[string]*model.Sample),
		aggregations: make(map[string]*model.Aggregation),
		exportJobs:   make(map[string]*model.ExportJob),
		alertRules:   make(map[string]*model.AlertRule),
	}
}

var _ Store = (*MemoryStore)(nil)
