package outbound

import (
	"context"
	"testing"

	pb "github.com/LeonLow97/authify/grpc-proto/pb"
	"github.com/LeonLow97/internal/contextstore"
	"google.golang.org/grpc/metadata"
)

func TestCurrentUserForwardsAuthenticatedIdentity(t *testing.T) {
	stub := &authifyServiceClientStub{currentUserResponse: &pb.User{ID: 42, FirstName: "Ada", LastName: "Lovelace", Email: "ada@example.com", Active: true}}
	adapter := NewUserOutbound(stub)
	if _, err := adapter.GetCurrentUser(context.Background()); err == nil || stub.currentUserCalls != 0 {
		t.Fatal("must reject missing identity before calling gRPC")
	}
	ctx := contextstore.WithRequestID(contextstore.WithUserID(context.Background(), 42), "current-user-request")
	user, err := adapter.GetCurrentUser(ctx)
	if err != nil {
		t.Fatal(err)
	}
	md, _ := metadata.FromOutgoingContext(stub.currentUserContext)
	if len(md.Get("user_id")) != 1 || md.Get("user_id")[0] != "42" {
		t.Fatal("missing authenticated identity metadata")
	}
	if len(md.Get("x-request-id")) != 1 || md.Get("x-request-id")[0] != "current-user-request" {
		t.Fatal("missing request ID metadata")
	}
	if user.ID != 42 || user.FirstName != "Ada" || user.LastName != "Lovelace" || user.Email != "ada@example.com" {
		t.Fatal("incorrect profile mapping")
	}
}
