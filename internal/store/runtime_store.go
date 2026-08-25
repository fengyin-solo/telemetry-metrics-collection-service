package store

import (
	"context"
	"errors"
	"sync"
	"time"

	"metricscollector/internal/model"
)

type RetryQueue struct {
	delay time.Duration
}

func NewRetryQueue(delay time.Duration) *RetryQueue { return &RetryQueue{delay: delay} }

func (q *RetryQueue) Wait(ctx context.Context) error {
	timer := time.NewTimer(q.delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func StreamValues(values []float64, failAt int) (<-chan float64, <-chan error) {
	items := make(chan float64)
	errs := make(chan error, 1)
	go func() {
		defer close(items)
		defer close(errs)
		for i, value := range values {
			if i == failAt {
				errs <- errors.New("sample stream interrupted")
				return
			}
			items <- value
		}
	}()
	return items, errs
}

type LabelRegistry struct {
	mu     sync.RWMutex
	series map[string][]float64
}

func NewLabelRegistry() *LabelRegistry {
	return &LabelRegistry{series: make(map[string][]float64)}
}

func (r *LabelRegistry) Replace(key string, values []float64) {
	r.mu.Lock()
	r.series[key] = append([]float64(nil), values...)
	r.mu.Unlock()
}

func (r *LabelRegistry) Snapshot(key string) []float64 {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return append([]float64(nil), r.series[key]...)
}

type BatchCache struct {
	mu      sync.RWMutex
	batches map[string]model.ParsedBatch
}

func NewBatchCache() *BatchCache { return &BatchCache{batches: make(map[string]model.ParsedBatch)} }

func (c *BatchCache) Put(key string, batch model.ParsedBatch) {
	c.mu.Lock()
	c.batches[key] = batch
	c.mu.Unlock()
}

func (c *BatchCache) Get(key string) (model.ParsedBatch, bool) {
	c.mu.RLock()
	batch, ok := c.batches[key]
	c.mu.RUnlock()
	return batch, ok
}

type SnapshotCache struct {
	mu        sync.RWMutex
	snapshots map[string]*model.AlertSnapshot
}

func NewSnapshotCache() *SnapshotCache {
	return &SnapshotCache{snapshots: make(map[string]*model.AlertSnapshot)}
}

func (c *SnapshotCache) Put(key string, snapshot *model.AlertSnapshot) error {
	if snapshot == nil || !snapshot.Ready {
		return errors.New("alert snapshot is incomplete")
	}
	c.mu.Lock()
	c.snapshots[key] = snapshot.Clone()
	c.mu.Unlock()
	return nil
}

func (c *SnapshotCache) Get(key string) (*model.AlertSnapshot, bool) {
	c.mu.RLock()
	snapshot, ok := c.snapshots[key]
	c.mu.RUnlock()
	return snapshot.Clone(), ok
}

type TaskStore struct {
	mu    sync.Mutex
	tasks map[string]model.TaskUpdate
}

func NewTaskStore() *TaskStore { return &TaskStore{tasks: make(map[string]model.TaskUpdate)} }

func (s *TaskStore) Apply(update model.TaskUpdate) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	current, ok := s.tasks[update.ID]
	if ok && !update.NewerThan(current) {
		return false
	}
	s.tasks[update.ID] = update
	return true
}

func (s *TaskStore) Get(id string) (model.TaskUpdate, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	value, ok := s.tasks[id]
	return value, ok
}

type ConfigStore struct {
	mu     sync.Mutex
	values map[string]string
}

func NewConfigStore(defaults map[string]string) *ConfigStore {
	values := make(map[string]string, len(defaults))
	for key, value := range defaults {
		values[key] = value
	}
	return &ConfigStore{values: values}
}

func (s *ConfigStore) Merge(values map[string]string) map[string]string {
	s.mu.Lock()
	defer s.mu.Unlock()
	for key, value := range values {
		s.values[key] = value
	}
	result := make(map[string]string, len(s.values))
	for key, value := range s.values {
		result[key] = value
	}
	return result
}

type EventTx struct {
	CommitErr   error
	RollbackErr error
	committed   bool
}

func (tx *EventTx) Commit() error {
	if tx.CommitErr != nil {
		return tx.CommitErr
	}
	tx.committed = true
	return nil
}

func (tx *EventTx) Rollback() error {
	tx.committed = false
	return tx.RollbackErr
}

func (tx *EventTx) Committed() bool { return tx.committed }

type BatchTx struct {
	mu        sync.Mutex
	committed []string
}

func (tx *BatchTx) Commit(items []string) {
	tx.mu.Lock()
	tx.committed = append(tx.committed, items...)
	tx.mu.Unlock()
}

func (tx *BatchTx) Items() []string {
	tx.mu.Lock()
	defer tx.mu.Unlock()
	return append([]string(nil), tx.committed...)
}
