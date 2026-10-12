package outbound

import (
	"errors"
	"testing"

	coreerrors "github.com/LeonLow97/internal/core/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestMapGRPCError(t *testing.T) {
	testCases := []struct {
		name        string
		err         error
		wantKind    coreerrors.Kind
		wantMessage string
	}{
		{name: "nil error", err: nil},
		{name: "invalid argument", err: status.Error(codes.InvalidArgument, "invalid input"), wantKind: coreerrors.KindInvalidArgument, wantMessage: "invalid input"},
		{name: "unauthenticated", err: status.Error(codes.Unauthenticated, "credentials invalid"), wantKind: coreerrors.KindUnauthenticated, wantMessage: "Unauthorized"},
		{name: "permission denied", err: status.Error(codes.PermissionDenied, "forbidden"), wantKind: coreerrors.KindPermissionDenied, wantMessage: "forbidden"},
		{name: "not found", err: status.Error(codes.NotFound, "missing"), wantKind: coreerrors.KindNotFound, wantMessage: "missing"},
		{name: "already exists", err: status.Error(codes.AlreadyExists, "duplicate"), wantKind: coreerrors.KindAlreadyExists, wantMessage: "duplicate"},
		{name: "failed precondition", err: status.Error(codes.FailedPrecondition, "precondition"), wantKind: coreerrors.KindFailedPrecondition, wantMessage: "precondition"},
		{name: "out of range", err: status.Error(codes.OutOfRange, "range"), wantKind: coreerrors.KindOutOfRange, wantMessage: "range"},
		{name: "unavailable", err: status.Error(codes.Unavailable, "internal detail"), wantKind: coreerrors.KindUnavailable, wantMessage: "Service temporarily unavailable"},
		{name: "deadline exceeded", err: status.Error(codes.DeadlineExceeded, "internal detail"), wantKind: coreerrors.KindDeadlineExceeded, wantMessage: "Request timed out"},
		{name: "resource exhausted", err: status.Error(codes.ResourceExhausted, "internal detail"), wantKind: coreerrors.KindResourceExhausted, wantMessage: "Rate limit exceeded"},
		{name: "aborted maps to conflict", err: status.Error(codes.Aborted, "conflict"), wantKind: coreerrors.KindConflict, wantMessage: "conflict"},
		{name: "unknown status is internal", err: status.Error(codes.Internal, "private detail"), wantKind: coreerrors.KindInternal, wantMessage: "Internal server error"},
		{name: "non status error is internal", err: errors.New("transport failure"), wantKind: coreerrors.KindInternal, wantMessage: "Internal server error"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := mapGRPCError(testCase.err)
			if testCase.err == nil {
				if got != nil {
					t.Fatalf("mapGRPCError(nil) = %v, want nil", got)
				}
				return
			}

			var appErr *coreerrors.Error
			if !errors.As(got, &appErr) {
				t.Fatalf("mapGRPCError() error type = %T, want *coreerrors.Error", got)
			}
			if appErr.Kind != testCase.wantKind {
				t.Errorf("Kind = %q, want %q", appErr.Kind, testCase.wantKind)
			}
			if appErr.Message != testCase.wantMessage {
				t.Errorf("Message = %q, want %q", appErr.Message, testCase.wantMessage)
			}
			if !errors.Is(got, testCase.err) {
				t.Errorf("mapped error does not wrap original error")
			}
		})
	}
}
