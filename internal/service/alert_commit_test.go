package service_test

import (
	"errors"
	"testing"

	"metricscollector/internal/service"
	"metricscollector/internal/store"
)

type countingPublisher struct{ count int }

func (p *countingPublisher) Publish(string) error {
	p.count++
	return nil
}

func TestAlertSuccessPublishesOnlyAfterCommit(t *testing.T) {
	commitErr := errors.New("alert transaction commit failed")
	rollbackErr := errors.New("alert transaction rollback failed")
	tx := &store.EventTx{CommitErr: commitErr, RollbackErr: rollbackErr}
	publisher := &countingPublisher{}
	err := service.CommitAlertEvent(tx, publisher, "cpu-critical")
	if publisher.count != 0 {
		t.Fatalf("failed alert transaction still published a success event: count=%d", publisher.count)
	}
	if !errors.Is(err, commitErr) || !errors.Is(err, rollbackErr) {
		t.Fatalf("alert failure lost commit or rollback details: %v", err)
	}
}
