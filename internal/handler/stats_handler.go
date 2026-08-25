package handler

import (
	"net/http"
	"strconv"

	"metricscollector/pkg/httpx"
)

func (s *Server) registerStatsRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/stats/overview", s.statsOverview)
	mux.HandleFunc("GET /api/stats/metric-detail/{id}", s.metricDetail)
	mux.HandleFunc("GET /api/stats/metric-history/{id}", s.metricHistory)
	mux.HandleFunc("GET /api/stats/metric-aggregation-summary/{id}", s.metricAggregationSummary)
}

func (s *Server) statsOverview(w http.ResponseWriter, r *http.Request) {
	overview, err := s.svc.GetStatsOverview()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, overview)
}

func (s *Server) metricDetail(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	detail, err := s.svc.GetMetricDetail(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, detail)
}

func (s *Server) metricHistory(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	limit := 100
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			limit = n
		}
	}
	history, err := s.svc.GetMetricHistory(id, limit)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, history)
}

func (s *Server) metricAggregationSummary(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	summary, err := s.svc.GetMetricAggregationSummary(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, summary)
}
