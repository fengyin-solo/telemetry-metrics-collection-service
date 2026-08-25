package model

import (
	"time"
)

const (
	ExportFormatJSON = "json"
	ExportFormatCSV  = "csv"

	ExportStatusPending = "pending"
	ExportStatusRunning = "running"
	ExportStatusDone    = "done"
	ExportStatusFailed  = "failed"
)

var exportTransitions = map[string]map[string]bool{
	ExportStatusPending: {ExportStatusRunning: true, ExportStatusDone: true, ExportStatusFailed: true},
	ExportStatusRunning: {ExportStatusDone: true, ExportStatusFailed: true},
}

func ExportCanTransition(from, to string) bool {
	if m, ok := exportTransitions[from]; ok {
		return m[to]
	}
	return false
}

type ExportJob struct {
	ID        string    `json:"id"`
	Format    string    `json:"format"`
	Filter    string    `json:"filter"`
	Status    string    `json:"status"`
	FilePath  string    `json:"file_path"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (e *ExportJob) Validate() error {
	if e.Format == "" {
		e.Format = ExportFormatJSON
	}
	if e.Format != ExportFormatJSON && e.Format != ExportFormatCSV {
		return NewValidationError("format", "导出格式不合法，应为 json/csv")
	}
	if e.Status == "" {
		e.Status = ExportStatusPending
	}
	if e.Status != ExportStatusPending && e.Status != ExportStatusRunning &&
		e.Status != ExportStatusDone && e.Status != ExportStatusFailed {
		return NewValidationError("status", "任务状态不合法")
	}
	return nil
}

type ExportJobFilter struct {
	Format string
	Status string
}

func (f ExportJobFilter) Match(e *ExportJob) bool {
	if f.Format != "" && e.Format != f.Format {
		return false
	}
	if f.Status != "" && e.Status != f.Status {
		return false
	}
	return true
}
