package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/TechBowl-japan/go-stations/model"
	"github.com/TechBowl-japan/go-stations/service"
)

// A TODOHandler implements handling REST endpoints.
type TODOHandler struct {
	svc *service.TODOService
}

// NewTODOHandler returns TODOHandler based http.Handler.
func NewTODOHandler(svc *service.TODOService) *TODOHandler {
	return &TODOHandler{
		svc: svc,
	}
}

func (h *TODOHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// POST
	if r.Method == http.MethodPost {
		var req model.CreateTODORequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		// validate subject
		if req.Subject == "" {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		result, err := h.Create(r.Context(), &req)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		res := model.CreateTODOResponse{TODO: result.TODO}
		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	} else if r.Method == http.MethodPut {
		var req model.UpdateTODORequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		// validate id
		if req.ID == 0 {
			http.Error(w, "id must not be null or empty", http.StatusBadRequest)
			return
		}
		// validate subject
		if req.Subject == "" {
			http.Error(w, "subject must not be null or empty", http.StatusBadRequest)
			return
		}

		result, err := h.Update(r.Context(), &req)
		if err != nil {
			switch err.(type){
			case model.ErrNotFound:
				http.Error(w, "Not Found", http.StatusNotFound)
				return
			default:
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}
		res := model.UpdateTODOResponse{TODO: result.TODO}
		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}

// Create handles the endpoint that creates the TODO.
func (h *TODOHandler) Create(ctx context.Context, req *model.CreateTODORequest) (*model.CreateTODOResponse, error) {
	res, err := h.svc.CreateTODO(ctx, req.Subject, req.Description)
	if err != nil {
		return nil, err
	}
	return &model.CreateTODOResponse{
		TODO: model.TODO{
			ID:          res.ID,
			Subject:     res.Subject,
			Description: res.Description,
			CreatedAt:   res.CreatedAt,
			UpdatedAt:   res.UpdatedAt,
		},
	}, nil
}

// Read handles the endpoint that reads the TODOs.
func (h *TODOHandler) Read(ctx context.Context, req *model.ReadTODORequest) (*model.ReadTODOResponse, error) {
	_, _ = h.svc.ReadTODO(ctx, 0, 0)
	return &model.ReadTODOResponse{}, nil
}

// Update handles the endpoint that updates the TODO.
func (h *TODOHandler) Update(ctx context.Context, req *model.UpdateTODORequest) (*model.UpdateTODOResponse, error) {
	res, err := h.svc.UpdateTODO(ctx, int64(req.ID), req.Subject, req.Description)
	if err != nil {
		return nil, err
	}

	return &model.UpdateTODOResponse{
		TODO: model.TODO{
			ID:          res.ID,
			Subject:     res.Subject,
			Description: req.Description,
			CreatedAt:   res.CreatedAt,
			UpdatedAt:   res.UpdatedAt,
		},
	}, nil
}

// Delete handles the endpoint that deletes the TODOs.
func (h *TODOHandler) Delete(ctx context.Context, req *model.DeleteTODORequest) (*model.DeleteTODOResponse, error) {
	_ = h.svc.DeleteTODO(ctx, nil)
	return &model.DeleteTODOResponse{}, nil
}
