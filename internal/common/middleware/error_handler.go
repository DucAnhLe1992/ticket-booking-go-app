package middleware

import (
	"log"
	"net/http"

	apperr "github.com/DucAnhLe1992/ticket-booking-go-app/internal/common/errors"
)

// RecoverAndJSON converts panics or returned errors into JSON responses.
func RecoverAndJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("panic: %v", rec)
				apperr.WriteJSONError(w, apperr.NewBadRequest("unexpected error"))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// JSONError writes an application error in JSON form.
func JSONError(w http.ResponseWriter, err error) {
	apperr.WriteJSONError(w, err)
}
