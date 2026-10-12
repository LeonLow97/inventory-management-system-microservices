package contextstore

import (
	"context"
	"fmt"
	"strconv"

	"google.golang.org/grpc/metadata"
)

const requestIDMetadataKey = "x-request-id"

type requestIDContextKey struct{}

func UserIDFromContext(ctx context.Context) (int64, error) {
	// Read the authenticated user ID from incoming gRPC metadata.
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, fmt.Errorf("grpc metadata not found in context")
	}
	userIDValues := md.Get("user_id")
	if len(userIDValues) == 0 {
		return 0, fmt.Errorf("user_id metadata not found in context")
	}
	userID, err := strconv.ParseInt(userIDValues[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to convert user ID in grpc metadata to int: %w", err)
	}
	return userID, nil
}

func RequestIDFromIncomingMetadata(ctx context.Context) (string, bool) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", false
	}

	requestIDs := md.Get(requestIDMetadataKey)
	if len(requestIDs) == 0 || requestIDs[0] == "" {
		return "", false
	}
	return requestIDs[0], true
}

func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDContextKey{}, requestID)
}

func RequestIDFromContext(ctx context.Context) (string, bool) {
	requestID, ok := ctx.Value(requestIDContextKey{}).(string)
	return requestID, ok
}
