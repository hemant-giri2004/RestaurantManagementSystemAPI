package middleware

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func RequireRolesMiddleware(allowedRoles ...string) func(http.Handler) http.Handler {
	// Create a set of allowed roles for efficient lookup
	roleSet := make(map[string]struct{})
	for _, r := range allowedRoles {
		roleSet[strings.ToLower(r)] = struct{}{}
	}

	// Return a middleware function
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// Extract roles from context
			rawRoles := ctx.Value(RolesKey)
			userRoles, ok := rawRoles.([]string)
			if !ok || len(userRoles) == 0 {
				logrus.Info("No roles in context")
				http.Error(w, "Forbidden: no roles in token", http.StatusForbidden)
				return
			}

			// Check if user has any of the allowed roles
			for _, userRole := range userRoles {
				if _, found := roleSet[strings.ToLower(userRole)]; found {
					next.ServeHTTP(w, r) // Allowed, proceed
					return
				}
			}

			// No match found, deny access
			logrus.Info("No roles in context")
			http.Error(w, "Forbidden: insufficient permissions", http.StatusForbidden)
		})
	}
}
