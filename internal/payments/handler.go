package payments

import (
	"encoding/json"
	"io"
	"net/http"
)

// HTTPHandler exposes HTTP endpoints for payments.
type HTTPHandler struct {
	svc *Service
}

func NewHTTPHandler(svc *Service) *HTTPHandler { return &HTTPHandler{svc: svc} }

type CreateChargeRequest struct {
	OrderID      string `json:"order_id"`
	Amount       int64  `json:"amount"`
	Currency     string `json:"currency"`
	PaymentMethod string `json:"payment_method"`
}

func (h *HTTPHandler) CreateCharge(w http.ResponseWriter, r *http.Request) {
	var req CreateChargeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	payment, err := h.svc.CreateCharge(r.Context(), req.OrderID, req.Amount, req.Currency)
	if err != nil {
		http.Error(w, "create charge error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(payment)
}

// Webhook reads raw body and forwards to service for verification/processing.
func (h *HTTPHandler) Webhook(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	// In real impl: verify signature header and call ProcessWebhook
	_ = h.svc.ProcessWebhook(r.Context(), body)
	w.WriteHeader(http.StatusOK)
}
