package middleware

type contextKey string

const (
	UserIDKey contextKey = "userID"
	RolesKey  contextKey = "roles"
)
