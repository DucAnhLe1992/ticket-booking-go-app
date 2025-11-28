package tickets

import (
	"encoding/json"
	"net/http"

	cmw "github.com/DucAnhLe1992/ticket-booking-go-app/internal/common/middleware"
)

type HTTPHandler struct{ svc *Service }

func NewHTTPHandler(s *Service) *HTTPHandler { return &HTTPHandler{svc: s} }

type createReq struct {
	Title string `json:"title"`
	Price int64  `json:"price"`
}
type updateReq struct {
	Title   string `json:"title"`
	Price   int64  `json:"price"`
	Version int    `json:"version"`
}

func (h *HTTPHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Title == "" || req.Price <= 0 {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	cu := cmw.GetCurrentUser(r.Context())
	if cu == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	t, err := h.svc.Create(r.Context(), req.Title, req.Price, cu.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(t)
}

func (h *HTTPHandler) Update(w http.ResponseWriter, r *http.Request) {
	var req updateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Title == "" || req.Price <= 0 {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}
	cu := cmw.GetCurrentUser(r.Context())
	if cu == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	t, err := h.svc.Update(r.Context(), id, req.Version, req.Title, req.Price, cu.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(t)
}

func (h *HTTPHandler) Show(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}
	t, err := h.svc.Get(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if t == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(t)
}

func (h *HTTPHandler) Index(w http.ResponseWriter, r *http.Request) {
	pageSize := int64(50)
	_ = pageSize // placeholder for future pagination
	list, err := h.svc.List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// ensure consistent numeric encoding
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(list)
}
