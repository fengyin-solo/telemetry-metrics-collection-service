package service

import (
	"sort"
	"time"

	"metricscollector/internal/model"
	"metricscollector/pkg/idgen"
)

func (s *Service) CreateAlertRule(input model.AlertRule) (*model.AlertRule, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	a := &model.AlertRule{
		ID:        idgen.Hex(),
		MetricID:  input.MetricID,
		Threshold: input.Threshold,
		Operator:  input.Operator,
		Level:     input.Level,
		Status:    input.Status,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	if err := s.store.CreateAlertRule(a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *Service) ListAlertRules(filter model.AlertRuleFilter, page, size int) ([]*model.AlertRule, int, error) {
	all := s.store.ListAlertRules()
	matched := make([]*model.AlertRule, 0, len(all))
	for _, a := range all {
		if filter.Match(a) {
			matched = append(matched, a)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.AlertRule{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) GetAlertRule(id string) (*model.AlertRule, error) {
	return s.store.GetAlertRule(id)
}

func (s *Service) UpdateAlertRule(id string, input model.AlertRule) (*model.AlertRule, error) {
	a, err := s.store.GetAlertRule(id)
	if err != nil {
		return nil, err
	}
	a.MetricID = input.MetricID
	a.Threshold = input.Threshold
	a.Operator = input.Operator
	a.Level = input.Level
	a.Status = input.Status
	a.UpdatedAt = time.Now().UTC()
	if err := a.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateAlertRule(a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *Service) DeleteAlertRule(id string) error {
	return s.store.DeleteAlertRule(id)
}

func (s *Service) EvaluateAlertRules(metricID string, value float64) ([]*model.AlertRule, error) {
	all := s.store.ListAlertRules()
	var triggered []*model.AlertRule
	for _, a := range all {
		if a.MetricID == metricID && a.Status == model.AlertRuleStatusActive && a.Evaluate(value) {
			triggered = append(triggered, a)
		}
	}
	sort.Slice(triggered, func(i, j int) bool {
		if triggered[i].Level != triggered[j].Level {
			return triggered[i].Level == model.AlertLevelCritical
		}
		return triggered[i].CreatedAt.Before(triggered[j].CreatedAt)
	})
	return triggered, nil
}
