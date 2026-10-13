package handlers

import (
	"context"
	"errors"
	"testing"

	"github.com/LeonLow97/internal/core/domain"
	coreerrors "github.com/LeonLow97/internal/core/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestGetCurrentUserUsesAuthenticatedMetadata(t *testing.T) {
	cases := []struct {
		name, id   string
		serviceErr error
		code       codes.Code
		calls      int
	}{
		{"profile", "42", nil, codes.OK, 1},
		{"no metadata", "", nil, codes.Unauthenticated, 0},
		{"invalid metadata", "invalid", nil, codes.Unauthenticated, 0},
		{"nonpositive identity", "0", nil, codes.Unauthenticated, 0},
		{"deleted user", "42", coreerrors.ErrUserNotFound, codes.Unauthenticated, 1},
		{"inactive user", "42", coreerrors.ErrInactiveUser, codes.Unauthenticated, 1},
		{"database error", "42", errors.New("private database error"), codes.Internal, 1},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			service := &userServiceStub{loginUser: &domain.User{ID: 42, Email: "ada@example.com", FirstName: stringPointerForHandlerTest("Ada"), LastName: stringPointerForHandlerTest("Lovelace"), Active: true}, loginErr: tc.serviceErr}
			ctx := context.Background()
			if tc.id != "" {
				ctx = metadata.NewIncomingContext(ctx, metadata.Pairs("user_id", tc.id))
			}
			got, err := NewUserHandler(service).GetCurrentUser(ctx, &emptypb.Empty{})
			if status.Code(err) != tc.code {
				t.Fatalf("error = %v, want code %v", err, tc.code)
			}
			if service.currentUserCalls != tc.calls {
				t.Fatalf("calls = %d, want %d", service.currentUserCalls, tc.calls)
			}
			if tc.calls > 0 && service.currentUserID != 42 {
				t.Fatal("incorrect authenticated identity")
			}
			if tc.code == codes.OK && (got == nil || got.ID != 42 || got.FirstName != "Ada" || got.LastName != "Lovelace") {
				t.Fatal("incorrect profile mapping")
			}
			if tc.code != codes.OK && got != nil {
				t.Fatal("must not return a profile on error")
			}
		})
	}
}
