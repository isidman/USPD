// Package api translates HTTP requests into lending.Service calls and
// domain errors into HTTP status codes. It's the only package that knows
// HTTP exists — lending and domain are transport-agnostic, so a CLI or a
// gRPC service could reuse them without changes.
package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"toolshed/internal/domain"
	"toolshed/internal/lending"
)

type Handler struct {
	service *lending.Service
}

func NewHandler(service *lending.Service) *Handler {
	return &Handler{service: service}
}

// Routes wires this handler's endpoints onto a fresh ServeMux using Go's
// built-in method+path pattern matching (stdlib since 1.22) — deliberately
// no router dependency (CLAUDE.md rule 6: no extraneous components).
func (h *Handler) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /resources", h.createResource)
	mux.HandleFunc("GET /resources", h.listResources)
	mux.HandleFunc("POST /resources/{id}/checkout", h.checkOut)
	mux.HandleFunc("POST /loans/{id}/return", h.returnLoan)
	return mux
}

type resourceView struct {
	domain.Resource
	Available bool `json:"available"`
}

func (h *Handler) createResource(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Kind        domain.Kind       `json:"kind"`
		Name        string            `json:"name"`
		Description string            `json:"description"`
		Metadata    map[string]string `json:"metadata"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, domain.ErrInvalidInput("malformed JSON body"))
		return
	}
	created, err := h.service.AddResource(r.Context(), domain.Resource{
		Kind:        in.Kind,
		Name:        in.Name,
		Description: in.Description,
		Metadata:    in.Metadata,
	})
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

func (h *Handler) listResources(w http.ResponseWriter, r *http.Request) {
	resources, err := h.service.ListResources(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	views := make([]resourceView, 0, len(resources))
	for _, res := range resources {
		available, err := h.service.Available(r.Context(), res.ID)
		if err != nil {
			writeError(w, err)
			return
		}
		views = append(views, resourceView{Resource: res, Available: available})
	}
	writeJSON(w, http.StatusOK, views)
}

func (h *Handler) checkOut(w http.ResponseWriter, r *http.Request) {
	resourceID := r.PathValue("id")
	var in struct {
		BorrowerID  string `json:"borrower_id"`
		DurationHrs int    `json:"duration_hours"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, domain.ErrInvalidInput("malformed JSON body"))
		return
	}
	if in.DurationHrs <= 0 {
		in.DurationHrs = 24 * 7 // default: one week
	}
	loan, err := h.service.CheckOut(r.Context(), resourceID, in.BorrowerID, time.Duration(in.DurationHrs)*time.Hour)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, loan)
}

func (h *Handler) returnLoan(w http.ResponseWriter, r *http.Request) {
	loanID := r.PathValue("id")
	loan, err := h.service.Return(r.Context(), loanID)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, loan)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	var invalid domain.ErrInvalidInput
	var notFound domain.ErrNotFound
	var conflict domain.ErrConflict
	switch {
	case errors.As(err, &invalid):
		status = http.StatusBadRequest
	case errors.As(err, &notFound):
		status = http.StatusNotFound
	case errors.As(err, &conflict):
		status = http.StatusConflict
	}
	writeJSON(w, status, map[string]string{"error": err.Error()})
}
