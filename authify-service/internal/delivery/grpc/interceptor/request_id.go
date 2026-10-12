package interceptor

import (
	"context"

	"github.com/LeonLow97/internal/delivery/grpc/contextstore"
	"google.golang.org/grpc"
)

// RequestID copies the request ID from incoming gRPC metadata into the RPC context.
func RequestID() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if requestID, ok := contextstore.RequestIDFromIncomingMetadata(ctx); ok {
			ctx = contextstore.WithRequestID(ctx, requestID)
		}
		return handler(ctx, req)
	}
}
