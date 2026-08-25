package handler

import (
	"net/http"

	"metricscollector/internal/model"
	"metricscollector/pkg/httpx"
)

func (s *Server) registerMetricRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/metrics", s.createMetric)
	mux.HandleFunc("GET /api/metrics", s.listMetrics)
	mux.HandleFunc("GET /api/metrics/{id}", s.getMetric)
	mux.HandleFunc("PUT /api/metrics/{id}", s.updateMetric)
	mux.HandleFunc("DELETE /api/metrics/{id}", s.deleteMetric)
}

type createMetricRequest struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Unit        string `json:"unit"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

func (s *Server) createMetric(w http.ResponseWriter, r *http.Request) {
	var req createMetricRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	m, err := s.svc.CreateMetric(model.Metric{
		Name:        req.Name,
		Type:        req.Type,
		Unit:        req.Unit,
		Description: req.Description,
		Status:      req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, m)
}

func (s *Server) listMetrics(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.MetricFilter{
		Type:    r.URL.Query().Get("type"),
		Status:  r.URL.Query().Get("status"),
		Keyword: r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListMetrics(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getMetric(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	m, err := s.svc.GetMetric(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, m)
}

type updateMetricRequest struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Unit        string `json:"unit"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

func (s *Server) updateMetric(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateMetricRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	m, err := s.svc.UpdateMetric(id, model.Metric{
		Name:        req.Name,
		Type:        req.Type,
		Unit:        req.Unit,
		Description: req.Description,
		Status:      req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, m)
}

func (s *Server) deleteMetric(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteMetric(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
