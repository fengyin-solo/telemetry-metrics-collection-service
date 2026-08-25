package service_test

import (
	"testing"

	"metricscollector/internal/service"
	"metricscollector/internal/store"
)

func TestRetryKeepsTerminalStateAndSingleSideEffect(t *testing.T) {
	tasks := store.NewTaskStore()
	recorder := service.NewOperationRecorder()
	stale := make(chan struct{})
	close(stale)
	count := service.ApplyTaskRetry(tasks, recorder, "export-42", stale)
	state, ok := tasks.Get("export-42")
	if count != 1 {
		t.Fatalf("export retry executed the external side effect more than once: count=%d", count)
	}
	if !ok || state.Status != "done" || state.Version != 2 {
		t.Fatalf("late first-attempt callback replaced the completed export state: exists=%v state=%+v", ok, state)
	}
}
