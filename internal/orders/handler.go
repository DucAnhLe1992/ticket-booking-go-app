package orders

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	cmw "github.com/DucAnhLe1992/ticket-booking-go-app/internal/common/middleware"
)

type HTTPHandler struct {
	svc *Service
}

func NewHTTPHandler(s *Service) *HTTPHandler {
	return &HTTPHandler{svc: s}
}

type createOrderReq struct {
	TicketID string `json:"ticketId"`
}

// Create creates a new order for the current user.
func (h *HTTPHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createOrderReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.TicketID == "" {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	cu := cmw.GetCurrentUser(r.Context())
	if cu == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	order, err := h.svc.CreateOrder(r.Context(), cu.ID, req.TicketID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(order)
}

// Show retrieves a single order by ID.
func (h *HTTPHandler) Show(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "orderId")
	if orderID == "" {
		http.Error(w, "missing orderId", http.StatusBadRequest)
		return
	}

	cu := cmw.GetCurrentUser(r.Context())
	if cu == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	order, err := h.svc.GetOrder(r.Context(), orderID, cu.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(order)
}

// Index retrieves all orders for the current user.
func (h *HTTPHandler) Index(w http.ResponseWriter, r *http.Request) {
	cu := cmw.GetCurrentUser(r.Context())
	if cu == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	orders, err := h.svc.ListOrders(r.Context(), cu.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(orders)
}

// Delete cancels an order.
func (h *HTTPHandler) Delete(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "orderId")
	if orderID == "" {
		http.Error(w, "missing orderId", http.StatusBadRequest)
		return
	}

	cu := cmw.GetCurrentUser(r.Context())
	if cu == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if err := h.svc.CancelOrder(r.Context(), orderID, cu.ID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
