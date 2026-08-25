package service

import (
	"fmt"
	"sort"
	"time"

	"metricscollector/internal/model"
	"metricscollector/pkg/idgen"
)

func (s *Service) CreateExportJob(input model.ExportJob) (*model.ExportJob, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	e := &model.ExportJob{
		ID:        idgen.Hex(),
		Format:    input.Format,
		Filter:    input.Filter,
		Status:    model.ExportStatusPending,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	if err := s.store.CreateExportJob(e); err != nil {
		return nil, err
	}
	return e, nil
}

func (s *Service) ListExportJobs(filter model.ExportJobFilter, page, size int) ([]*model.ExportJob, int, error) {
	all := s.store.ListExportJobs()
	matched := make([]*model.ExportJob, 0, len(all))
	for _, e := range all {
		if filter.Match(e) {
			matched = append(matched, e)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.ExportJob{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) GetExportJob(id string) (*model.ExportJob, error) {
	return s.store.GetExportJob(id)
}

func (s *Service) UpdateExportJob(id string, input model.ExportJob) (*model.ExportJob, error) {
	e, err := s.store.GetExportJob(id)
	if err != nil {
		return nil, err
	}
	e.Format = input.Format
	e.Filter = input.Filter
	e.Status = input.Status
	e.FilePath = input.FilePath
	e.UpdatedAt = time.Now().UTC()
	if err := e.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateExportJob(e); err != nil {
		return nil, err
	}
	return e, nil
}

func (s *Service) DeleteExportJob(id string) error {
	return s.store.DeleteExportJob(id)
}

func (s *Service) ExecuteExportJob(id string) (*model.ExportJob, error) {
	e, err := s.store.GetExportJob(id)
	if err != nil {
		return nil, err
	}
	if e.Status != model.ExportStatusPending {
		return nil, model.NewValidationError("status", "任务不是 pending 状态，无法执行")
	}
	e.Status = model.ExportStatusRunning
	e.UpdatedAt = time.Now().UTC()
	if err := s.store.UpdateExportJob(e); err != nil {
		return nil, err
	}

	e.FilePath = fmt.Sprintf("/tmp/exports/%s_%s.%s", e.ID, time.Now().UTC().Format("20060102T150405"), e.Format)
	e.Status = model.ExportStatusDone
	e.UpdatedAt = time.Now().UTC()
	if err := s.store.UpdateExportJob(e); err != nil {
		return nil, err
	}
	return e, nil
}
