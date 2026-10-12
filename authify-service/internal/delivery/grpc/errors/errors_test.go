package grpcerror

import (
	"errors"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCustomErrorGRPCError(t *testing.T) {
	cause := errors.New("private cause")
	testCases := []struct {
		name    string
		custom  *CustomError
		want    codes.Code
		message string
	}{
		{name: "duplicate email", custom: ECEmailAlreadyExists, want: codes.AlreadyExists, message: "Email already exists"},
		{name: "unauthorized", custom: ECUnauthorized, want: codes.Unauthenticated, message: "Unauthenticated"},
		{name: "internal error", custom: ECInternalServerError, want: codes.Internal, message: "Internal Server Error"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := testCase.custom.GRPCError(cause)
			if status.Code(got) != testCase.want {
				t.Errorf("gRPC code = %s, want %s", status.Code(got), testCase.want)
			}
			if gotMessage := status.Convert(got).Message(); gotMessage != testCase.message {
				t.Errorf("gRPC message = %q, want %q", gotMessage, testCase.message)
			}
			if errors.Is(got, cause) {
				t.Error("gRPC status error unexpectedly exposes the original cause")
			}
		})
	}
}
