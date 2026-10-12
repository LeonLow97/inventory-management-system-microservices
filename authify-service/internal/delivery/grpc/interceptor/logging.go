package interceptor

import (
	"context"
	"runtime/debug"
	"time"

	"github.com/LeonLow97/internal/core/ports"
	"github.com/LeonLow97/internal/delivery/grpc/contextstore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RequestLogger records a concise summary for each completed unary RPC.
func RequestLogger(logger ports.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		startedAt := time.Now()
		response, err := handler(ctx, req)

		logger.WithContext(ctx).Info("gRPC request completed",
			ports.LogField{Key: "request_id", Value: requestIDFromContext(ctx)},
			ports.LogField{Key: "method", Value: info.FullMethod},
			ports.LogField{Key: "status", Value: status.Code(err).String()},
			ports.LogField{Key: "duration_ms", Value: time.Since(startedAt).Milliseconds()},
		)
		return response, err
	}
}

// Recovery logs a panic and converts it to a generic internal gRPC error.
func Recovery(logger ports.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (response any, err error) {
		defer func() {
			if recovered := recover(); recovered != nil {
				panicLogger := logger.WithContext(ctx)
				if panicErr, ok := recovered.(error); ok {
					panicLogger = panicLogger.WithError(panicErr)
				}
				panicLogger.Error("panic recovered while handling gRPC request",
					ports.LogField{Key: "request_id", Value: requestIDFromContext(ctx)},
					ports.LogField{Key: "method", Value: info.FullMethod},
					ports.LogField{Key: "panic", Value: recovered},
					ports.LogField{Key: "stack", Value: string(debug.Stack())},
				)
				response = nil
				err = status.Error(codes.Internal, "Internal server error")
			}
		}()

		return handler(ctx, req)
	}
}

func requestIDFromContext(ctx context.Context) string {
	requestID, _ := contextstore.RequestIDFromContext(ctx)
	return requestID
}
