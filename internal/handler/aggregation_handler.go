package handler

import (
	"net/http"
	"time"

	"metricscollector/internal/model"
	"metricscollector/pkg/httpx"
)

func (s *Server) registerAggregationRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/aggregations", s.createAggregation)
	mux.HandleFunc("GET /api/aggregations", s.listAggregations)
	mux.HandleFunc("GET /api/aggregations/{id}", s.getAggregation)
	mux.HandleFunc("POST /api/aggregations/compute", s.computeAggregation)
}

type createAggregationRequest struct {
	MetricID  string    `json:"metric_id"`
	WindowSec int       `json:"window_sec"`
	Func      string    `json:"func"`
	Result    float64   `json:"result"`
	Timestamp time.Time `json:"timestamp"`
}

func (s *Server) createAggregation(w http.ResponseWriter, r *http.Request) {
	var req createAggregationRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	a, err := s.svc.CreateAggregation(model.Aggregation{
		MetricID:  req.MetricID,
		WindowSec: req.WindowSec,
		Func:      req.Func,
		Result:    req.Result,
		Timestamp: req.Timestamp,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, a)
}

func (s *Server) listAggregations(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.AggregationFilter{
		MetricID: r.URL.Query().Get("metric_id"),
		Func:     r.URL.Query().Get("func"),
	}
	items, total, err := s.svc.ListAggregations(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getAggregation(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	a, err := s.svc.GetAggregation(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, a)
}

type computeAggregationRequest struct {
	MetricID  string `json:"metric_id"`
	WindowSec int    `json:"window_sec"`
	Func      string `json:"func"`
}

func (s *Server) computeAggregation(w http.ResponseWriter, r *http.Request) {
	var req computeAggregationRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	a, err := s.svc.ComputeAggregation(req.MetricID, req.WindowSec, req.Func)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, a)
}
