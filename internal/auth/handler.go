package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	cmw "github.com/DucAnhLe1992/ticket-booking-go-app/internal/common/middleware"
)

type HTTPHandler struct {
	svc  *Service
	repo UserRepository
}

func NewHTTPHandler(svc *Service, repo UserRepository) *HTTPHandler {
	return &HTTPHandler{svc: svc, repo: repo}
}

func (h *HTTPHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var in SignupInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	u, token, err := h.svc.Signup(r.Context(), in)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	setJWTCookie(w, token)
	writeCurrentUser(w, u)
}

func (h *HTTPHandler) Signin(w http.ResponseWriter, r *http.Request) {
	var in SigninInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	u, token, err := h.svc.Signin(r.Context(), in)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	setJWTCookie(w, token)
	writeCurrentUser(w, u)
}

func (h *HTTPHandler) Signout(w http.ResponseWriter, r *http.Request) {
	// expire cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"message":"signed out"}`))
}

func (h *HTTPHandler) CurrentUser(w http.ResponseWriter, r *http.Request) {
	cu := cmw.GetCurrentUser(r.Context())
	if cu == nil {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"currentUser":null}`))
		return
	}

	// hydrate latest email if desired
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	u, _ := h.repo.FindByID(ctx, cu.ID)
	if u == nil {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"currentUser":null}`))
		return
	}
	resp := struct {
		CurrentUser any `json:"currentUser"`
	}{CurrentUser: struct {
		ID    string `json:"id"`
		Email string `json:"email"`
	}{ID: u.ID, Email: u.Email}}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func setJWTCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(24 * time.Hour),
	})
}

func writeCurrentUser(w http.ResponseWriter, u *User) {
	resp := struct {
		CurrentUser any `json:"currentUser"`
	}{CurrentUser: struct {
		ID    string `json:"id"`
		Email string `json:"email"`
	}{ID: u.ID, Email: u.Email}}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
