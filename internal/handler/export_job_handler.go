package handler

import (
	"net/http"

	"metricscollector/internal/model"
	"metricscollector/pkg/httpx"
)

func (s *Server) registerExportJobRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/export-jobs", s.createExportJob)
	mux.HandleFunc("GET /api/export-jobs", s.listExportJobs)
	mux.HandleFunc("GET /api/export-jobs/{id}", s.getExportJob)
	mux.HandleFunc("PUT /api/export-jobs/{id}", s.updateExportJob)
	mux.HandleFunc("DELETE /api/export-jobs/{id}", s.deleteExportJob)
	mux.HandleFunc("POST /api/export-jobs/{id}/execute", s.executeExportJob)
}

type createExportJobRequest struct {
	Format string `json:"format"`
	Filter string `json:"filter"`
}

func (s *Server) createExportJob(w http.ResponseWriter, r *http.Request) {
	var req createExportJobRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	e, err := s.svc.CreateExportJob(model.ExportJob{
		Format: req.Format,
		Filter: req.Filter,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, e)
}

func (s *Server) listExportJobs(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.ExportJobFilter{
		Format: r.URL.Query().Get("format"),
		Status: r.URL.Query().Get("status"),
	}
	items, total, err := s.svc.ListExportJobs(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getExportJob(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	e, err := s.svc.GetExportJob(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, e)
}

type updateExportJobRequest struct {
	Format   string `json:"format"`
	Filter   string `json:"filter"`
	Status   string `json:"status"`
	FilePath string `json:"file_path"`
}

func (s *Server) updateExportJob(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateExportJobRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	e, err := s.svc.UpdateExportJob(id, model.ExportJob{
		Format:   req.Format,
		Filter:   req.Filter,
		Status:   req.Status,
		FilePath: req.FilePath,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, e)
}

func (s *Server) deleteExportJob(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteExportJob(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) executeExportJob(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	e, err := s.svc.ExecuteExportJob(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, e)
}
