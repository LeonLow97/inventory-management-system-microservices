package outbound

import (
	"fmt"

	coreerrors "github.com/LeonLow97/internal/core/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func mapGRPCError(err error) error {
	if err == nil {
		return nil
	}

	grpcStatus, ok := status.FromError(err)
	if !ok {
		return coreerrors.New(coreerrors.KindInternal, "Internal server error", err)
	}

	kind := coreerrors.KindInternal
	message := "Internal server error"
	switch grpcStatus.Code() {
	case codes.InvalidArgument:
		kind, message = coreerrors.KindInvalidArgument, grpcStatus.Message()
	case codes.Unauthenticated:
		kind, message = coreerrors.KindUnauthenticated, "Unauthorized"
	case codes.PermissionDenied:
		kind, message = coreerrors.KindPermissionDenied, grpcStatus.Message()
	case codes.NotFound:
		kind, message = coreerrors.KindNotFound, grpcStatus.Message()
	case codes.AlreadyExists:
		kind, message = coreerrors.KindAlreadyExists, grpcStatus.Message()
	case codes.FailedPrecondition:
		kind, message = coreerrors.KindFailedPrecondition, grpcStatus.Message()
	case codes.OutOfRange:
		kind, message = coreerrors.KindOutOfRange, grpcStatus.Message()
	case codes.Unavailable:
		kind, message = coreerrors.KindUnavailable, "Service temporarily unavailable"
	case codes.DeadlineExceeded:
		kind, message = coreerrors.KindDeadlineExceeded, "Request timed out"
	case codes.ResourceExhausted:
		kind, message = coreerrors.KindResourceExhausted, "Rate limit exceeded"
	case codes.Aborted:
		kind, message = coreerrors.KindConflict, grpcStatus.Message()
	default:
		return coreerrors.New(kind, message, fmt.Errorf("gRPC %s: %w", grpcStatus.Code(), err))
	}

	return coreerrors.New(kind, message, err)
}
