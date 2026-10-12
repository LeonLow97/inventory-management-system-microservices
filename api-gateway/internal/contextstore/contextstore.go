package contextstore

import "context"

type userIDContextKey struct{}
type requestIDContextKey struct{}

func WithUserID(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, userIDContextKey{}, userID)
}

func UserIDFromContext(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(userIDContextKey{}).(int64)
	return userID, ok
}

func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDContextKey{}, requestID)
}

func RequestIDFromContext(ctx context.Context) (string, bool) {
	requestID, ok := ctx.Value(requestIDContextKey{}).(string)
	return requestID, ok
}
