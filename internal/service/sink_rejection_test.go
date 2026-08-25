package service_test

import (
	"context"
	"errors"
	"testing"

	"metricscollector/internal/model"
	"metricscollector/internal/service"
	"metricscollector/internal/store"
)

type rejectingSink struct{ calls int }

func (s *rejectingSink) Send(context.Context, []string) error {
	s.calls++
	return model.NewRejectedSinkError(errors.New("metric schema denied"))
}

func TestRejectedBatchIsNotRetriedOrCommitted(t *testing.T) {
	tx := &store.BatchTx{}
	sink := &rejectingSink{}
	err := service.DeliverSampleBatch(context.Background(), tx, sink, []string{"cpu=91"}, 3)
	if !errors.Is(err, model.ErrSinkRejected) || sink.calls != 1 || len(tx.Items()) != 0 {
		t.Fatalf("rejected sample batch was retried or committed: calls=%d committed=%v err=%v", sink.calls, tx.Items(), err)
	}
}
