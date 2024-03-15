package server

import "context"

// type contextKey follows Go's suggestion to use custom types where
// setting context keys in order to avoid key collisions.
type contextKey int

const (
	contextUserID contextKey = iota
)

func setUserIDInContext(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, contextUserID, userID)
}

func getUserIdFromContext(ctx context.Context) string {
	return ctx.Value(contextUserID).(string)
}
