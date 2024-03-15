package server

import (
	"net/http"
)

const (
	XUserID         = "X-User-Id"
	contentType     = "Content-Type"
	applicationJSON = "application/json"
)

// AddUserIDToContext adds the X-User-Id header's value to the context
func AddUserIDToContext(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			userID := r.Header.Get(XUserID)
			if userID != "" {
				ctx := r.Context()
				ctx = setUserIDInContext(ctx, userID)
				r = r.WithContext(ctx)
			}
			next.ServeHTTP(w, r)
		},
	)
}

// AddContentTypeToResponse adds Content-Type application/json to the response
func AddContentTypeToResponse(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set(contentType, applicationJSON)
			next.ServeHTTP(w, r)
		},
	)
}
