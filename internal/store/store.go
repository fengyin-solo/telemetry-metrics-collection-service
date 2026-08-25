// Package store 定义数据访问接口与内存实现。
package store

import (
	"errors"

	"metricscollector/internal/model"
)

var (
	ErrNotFound = errors.New("记录不存在")
	ErrConflict = errors.New("记录已存在或状态冲突")
)

// Store 聚合全部实体的数据访问方法，便于测试时替换实现。
type Store interface {
	// Metric
	CreateMetric(m *model.Metric) error
	GetMetric(id string) (*model.Metric, error)
	GetMetricByName(name string) (*model.Metric, error)
	ListMetrics() []*model.Metric
	UpdateMetric(m *model.Metric) error
	DeleteMetric(id string) error

	// Sample
	CreateSample(s *model.Sample) error
	GetSample(id string) (*model.Sample, error)
	ListSamples() []*model.Sample
	ListSamplesByMetricID(metricID string) []*model.Sample

	// Aggregation
	CreateAggregation(a *model.Aggregation) error
	GetAggregation(id string) (*model.Aggregation, error)
	ListAggregations() []*model.Aggregation
	ListAggregationsByMetricIDFunc(metricID, fn string) []*model.Aggregation

	// ExportJob
	CreateExportJob(e *model.ExportJob) error
	GetExportJob(id string) (*model.ExportJob, error)
	ListExportJobs() []*model.ExportJob
	UpdateExportJob(e *model.ExportJob) error
	DeleteExportJob(id string) error

	// AlertRule
	CreateAlertRule(a *model.AlertRule) error
	GetAlertRule(id string) (*model.AlertRule, error)
	ListAlertRules() []*model.AlertRule
	UpdateAlertRule(a *model.AlertRule) error
	DeleteAlertRule(id string) error
}
