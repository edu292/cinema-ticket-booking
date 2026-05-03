package utils

import (
	"context"
	"encoding/json"
	"net/http"
)

type ContextKey string

const (
	ContextUserID ContextKey = "UserID"
)

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

type userIDPayload struct {
	UserID string `json:"user_id"`
}

func AuthMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var userIDPayload userIDPayload
		err := json.NewDecoder(r.Body).Decode(&userIDPayload)
		if err != nil {
			WriteJSON(w, http.StatusBadRequest, nil)
		}

		userID := userIDPayload.UserID
		if userID == "" {
			WriteJSON(w, http.StatusUnauthorized, nil)
		}

		ctx := context.WithValue(r.Context(), ContextUserID, userIDPayload.UserID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
