package auth

import "context"

type contextKey string

const userIDKey contextKey = "user_id"

// WithUserID сохраняет идентификатор пользователя в контексте
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// UserIDFromContext возвращает идентификатор пользователя из контекста
func UserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userIDKey).(string)
	return userID, ok
}
