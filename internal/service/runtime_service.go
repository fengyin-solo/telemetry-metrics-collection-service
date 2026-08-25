package service

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"metricscollector/internal/model"
	"metricscollector/internal/store"
)

func RunExportWithRetry(ctx context.Context, queue *store.RetryQueue, attempts int, run func(context.Context) error) error {
	var last error
	for attempt := 0; attempt < attempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		last = run(ctx)
		if last == nil {
			return nil
		}
		if attempt+1 < attempts {
			if err := queue.Wait(ctx); err != nil {
				return err
			}
		}
	}
	return last
}

func DispatchEnvelope(pool *model.EnvelopePool, tenant, metric string, payload []byte, hold <-chan struct{}, observed chan<- model.SampleEnvelope) {
	envelope := pool.Acquire()
	envelope.Tenant = tenant
	envelope.Metric = metric
	envelope.Payload = append(envelope.Payload[:0], payload...)
	envelope.Labels["tenant"] = tenant
	snapshot := envelope.Clone()
	pool.Release(envelope)
	go func() {
		<-hold
		observed <- snapshot
	}()
}

func RefreshAlertSnapshot(cache *store.SnapshotCache, key string, values map[string]*float64) error {
	snapshot, err := model.BuildAlertSnapshot(values)
	if err != nil {
		return err
	}
	return cache.Put(key, snapshot)
}

func AggregateStream(ctx context.Context, values []float64, failAt int) (float64, error) {
	items, errs := store.StreamValues(values, failAt)
	var total float64
	for items != nil || errs != nil {
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		case value, ok := <-items:
			if !ok {
				items = nil
				continue
			}
			total += value
		case err, ok := <-errs:
			if !ok {
				errs = nil
				continue
			}
			if err != nil {
				return 0, err
			}
		}
	}
	return total, nil
}

func SumStableSeries(registry *store.LabelRegistry, key string, updated []float64) float64 {
	snapshot := append([]float64(nil), registry.Snapshot(key)...)
	registry.Replace(key, updated)
	var total float64
	for _, value := range snapshot {
		total += value
	}
	return total
}

func ParseAndCacheBatch(parser *model.ReusableParser, cache *store.BatchCache, key, metric string, frames ...string) {
	parsed := parser.Parse(metric, frames...).Clone()
	cache.Put(key, parsed)
}

type OperationRecorder struct {
	mu   sync.Mutex
	keys map[string]struct{}
}

func NewOperationRecorder() *OperationRecorder {
	return &OperationRecorder{keys: make(map[string]struct{})}
}

func (r *OperationRecorder) Record(key string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.keys[key]; exists {
		return false
	}
	r.keys[key] = struct{}{}
	return true
}

func ApplyTaskRetry(tasks *store.TaskStore, recorder *OperationRecorder, id string, stale <-chan struct{}) int {
	key := "export:" + id
	tasks.Apply(model.TaskUpdate{ID: id, Version: 1, Status: "running", OperationKey: key})
	count := 0
	if recorder.Record(key) {
		count++
	}
	tasks.Apply(model.TaskUpdate{ID: id, Version: 2, Status: "done", OperationKey: key})
	<-stale
	tasks.Apply(model.TaskUpdate{ID: id, Version: 1, Status: "running", OperationKey: key})
	return count
}

func LoadRuntimeConfig(configs *store.ConfigStore, validator model.Validator, values map[string]string) (map[string]string, error) {
	merged := configs.Merge(values)
	if validator != nil {
		if err := validator.Validate(merged); err != nil {
			return nil, err
		}
	}
	return merged, nil
}

type EventPublisher interface {
	Publish(string) error
}

func CommitAlertEvent(tx *store.EventTx, publisher EventPublisher, event string) error {
	if err := tx.Commit(); err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return errors.Join(err, rollbackErr)
		}
		return err
	}
	if err := publisher.Publish(event); err != nil {
		return fmt.Errorf("publish committed alert: %w", err)
	}
	return nil
}

type SampleSink interface {
	Send(context.Context, []string) error
}

func DeliverSampleBatch(ctx context.Context, tx *store.BatchTx, sink SampleSink, items []string, attempts int) error {
	var last error
	for attempt := 0; attempt < attempts; attempt++ {
		last = sink.Send(ctx, items)
		if last == nil {
			tx.Commit(items)
			return nil
		}
		if !model.IsTemporarySinkError(last) {
			return last
		}
	}
	return last
}
