package middleware

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"rms/database/dbHelper"
	"strings"
	"time"
)

func SessionValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := ""

		// get refresh token
		authHeader := r.Header.Get("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		}

		if token == "" {
			logrus.Error("token is empty")
			http.Error(w, "missing refresh token", http.StatusUnauthorized)
			return
		}

		session, err := dbHelper.GetSessionByToken(token)
		if err != nil {
			logrus.Errorf("session not found: %v", err)
			http.Error(w, "invalid or expired session", http.StatusUnauthorized)
			return
		}

		if session.ExpiresAt.Before(time.Now()) {
			logrus.Error("session expired")
			http.Error(w, "session expired", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
