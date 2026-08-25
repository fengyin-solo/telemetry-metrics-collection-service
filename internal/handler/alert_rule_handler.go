package handler

import (
	"net/http"

	"metricscollector/internal/model"
	"metricscollector/pkg/httpx"
)

func (s *Server) registerAlertRuleRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/alert-rules", s.createAlertRule)
	mux.HandleFunc("GET /api/alert-rules", s.listAlertRules)
	mux.HandleFunc("GET /api/alert-rules/{id}", s.getAlertRule)
	mux.HandleFunc("PUT /api/alert-rules/{id}", s.updateAlertRule)
	mux.HandleFunc("DELETE /api/alert-rules/{id}", s.deleteAlertRule)
	mux.HandleFunc("POST /api/alert-rules/evaluate", s.evaluateAlertRules)
}

type createAlertRuleRequest struct {
	MetricID  string  `json:"metric_id"`
	Threshold float64 `json:"threshold"`
	Operator  string  `json:"operator"`
	Level     string  `json:"level"`
	Status    string  `json:"status"`
}

func (s *Server) createAlertRule(w http.ResponseWriter, r *http.Request) {
	var req createAlertRuleRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	a, err := s.svc.CreateAlertRule(model.AlertRule{
		MetricID:  req.MetricID,
		Threshold: req.Threshold,
		Operator:  req.Operator,
		Level:     req.Level,
		Status:    req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, a)
}

func (s *Server) listAlertRules(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.AlertRuleFilter{
		MetricID: r.URL.Query().Get("metric_id"),
		Level:    r.URL.Query().Get("level"),
		Status:   r.URL.Query().Get("status"),
	}
	items, total, err := s.svc.ListAlertRules(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getAlertRule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	a, err := s.svc.GetAlertRule(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, a)
}

type updateAlertRuleRequest struct {
	MetricID  string  `json:"metric_id"`
	Threshold float64 `json:"threshold"`
	Operator  string  `json:"operator"`
	Level     string  `json:"level"`
	Status    string  `json:"status"`
}

func (s *Server) updateAlertRule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateAlertRuleRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	a, err := s.svc.UpdateAlertRule(id, model.AlertRule{
		MetricID:  req.MetricID,
		Threshold: req.Threshold,
		Operator:  req.Operator,
		Level:     req.Level,
		Status:    req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, a)
}

func (s *Server) deleteAlertRule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteAlertRule(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

type evaluateAlertRulesRequest struct {
	MetricID string  `json:"metric_id"`
	Value    float64 `json:"value"`
}

func (s *Server) evaluateAlertRules(w http.ResponseWriter, r *http.Request) {
	var req evaluateAlertRulesRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	triggered, err := s.svc.EvaluateAlertRules(req.MetricID, req.Value)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]interface{}{"triggered": triggered})
}
