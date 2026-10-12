package outbound

import (
	"context"
	"reflect"
	"testing"

	pb "github.com/LeonLow97/authify/grpc-proto/pb"
	"github.com/LeonLow97/internal/contextstore"
	"github.com/LeonLow97/internal/core/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

type authifyServiceClientStub struct {
	loginRequest     *pb.LoginRequest
	loginResponse    *pb.LoginResponse
	loginErr         error
	signUpRequest    *pb.SignUpRequest
	signUpErr        error
	getUsersRequest  *pb.GetUsersRequest
	getUsersContext  context.Context
	getUsersResponse *pb.GetUsersResponse
	getUsersErr      error
	updateRequest    *pb.UpdateUserRequest
	updateContext    context.Context
	updateErr        error
}

func (s *authifyServiceClientStub) Login(_ context.Context, req *pb.LoginRequest, _ ...grpc.CallOption) (*pb.LoginResponse, error) {
	s.loginRequest = req
	return s.loginResponse, s.loginErr
}

func (s *authifyServiceClientStub) SignUp(_ context.Context, req *pb.SignUpRequest, _ ...grpc.CallOption) (*emptypb.Empty, error) {
	s.signUpRequest = req
	return &emptypb.Empty{}, s.signUpErr
}

func (s *authifyServiceClientStub) GetUsers(ctx context.Context, req *pb.GetUsersRequest, _ ...grpc.CallOption) (*pb.GetUsersResponse, error) {
	s.getUsersContext = ctx
	s.getUsersRequest = req
	return s.getUsersResponse, s.getUsersErr
}

func (s *authifyServiceClientStub) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest, _ ...grpc.CallOption) (*emptypb.Empty, error) {
	s.updateContext = ctx
	s.updateRequest = req
	return &emptypb.Empty{}, s.updateErr
}

func TestUserOutboundLogin(t *testing.T) {
	stub := &authifyServiceClientStub{loginResponse: &pb.LoginResponse{
		FirstName: "Ada", LastName: "Lovelace", Email: "ada@example.com", Active: true, Admin: true, Token: "signed-token",
	}}
	adapter := NewUserOutbound(stub)

	gotUser, gotToken, err := adapter.Login(context.Background(), domain.User{Email: "ada@example.com", Password: "Secret1!"})
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}
	if want := (&pb.LoginRequest{Email: "ada@example.com", Password: "Secret1!"}); !reflect.DeepEqual(stub.loginRequest, want) {
		t.Errorf("gRPC login request = %#v, want %#v", stub.loginRequest, want)
	}
	wantUser := &domain.User{FirstName: "Ada", LastName: "Lovelace", Email: "ada@example.com", Active: true, Admin: true}
	if !reflect.DeepEqual(gotUser, wantUser) {
		t.Errorf("Login() user = %#v, want %#v", gotUser, wantUser)
	}
	if gotToken != "signed-token" {
		t.Errorf("Login() token = %q, want signed-token", gotToken)
	}
}

func TestUserOutboundSignUp(t *testing.T) {
	stub := &authifyServiceClientStub{}
	adapter := NewUserOutbound(stub)
	request := domain.User{FirstName: "Ada", LastName: "Lovelace", Email: "ada@example.com", Password: "Secret1!"}

	if err := adapter.SignUp(context.Background(), request); err != nil {
		t.Fatalf("SignUp() error = %v", err)
	}
	want := &pb.SignUpRequest{FirstName: request.FirstName, LastName: request.LastName, Email: request.Email, Password: request.Password}
	if !reflect.DeepEqual(stub.signUpRequest, want) {
		t.Errorf("gRPC signup request = %#v, want %#v", stub.signUpRequest, want)
	}
}

func TestUserOutboundGetUsersAddsIdentityMetadata(t *testing.T) {
	stub := &authifyServiceClientStub{getUsersResponse: &pb.GetUsersResponse{
		Users:      []*pb.User{{ID: 17, FirstName: "Ada", LastName: "Lovelace", Email: "ada@example.com", Active: true, Admin: true}},
		NextCursor: "next-page",
	}}
	adapter := NewUserOutbound(stub)
	ctx := contextstore.WithUserID(context.Background(), 42)

	users, nextCursor, err := adapter.GetUsers(ctx, 25, "current-page")
	if err != nil {
		t.Fatalf("GetUsers() error = %v", err)
	}
	if want := (&pb.GetUsersRequest{Limit: 25, Cursor: "current-page"}); !reflect.DeepEqual(stub.getUsersRequest, want) {
		t.Errorf("gRPC get-users request = %#v, want %#v", stub.getUsersRequest, want)
	}
	wantUsers := []domain.User{{ID: 17, FirstName: "Ada", LastName: "Lovelace", Email: "ada@example.com", Active: true, Admin: true}}
	if !reflect.DeepEqual(users, wantUsers) {
		t.Errorf("GetUsers() users = %#v, want %#v", users, wantUsers)
	}
	if nextCursor != "next-page" {
		t.Errorf("GetUsers() cursor = %q, want next-page", nextCursor)
	}
	md, ok := metadata.FromOutgoingContext(stub.getUsersContext)
	if !ok || !reflect.DeepEqual(md.Get("user_id"), []string{"42"}) {
		t.Errorf("outgoing user_id metadata = %v, want [42]", md.Get("user_id"))
	}
}

func TestUserOutboundUpdateUser(t *testing.T) {
	stub := &authifyServiceClientStub{}
	adapter := NewUserOutbound(stub)
	request := domain.User{FirstName: "Ada", LastName: "Lovelace", Password: "Secret1!"}
	ctx := contextstore.WithUserID(context.Background(), 9)

	if err := adapter.UpdateUser(ctx, request); err != nil {
		t.Fatalf("UpdateUser() error = %v", err)
	}
	want := &pb.UpdateUserRequest{FirstName: request.FirstName, LastName: request.LastName, Password: request.Password}
	if !reflect.DeepEqual(stub.updateRequest, want) {
		t.Errorf("gRPC update request = %#v, want %#v", stub.updateRequest, want)
	}
	md, ok := metadata.FromOutgoingContext(stub.updateContext)
	if !ok || !reflect.DeepEqual(md.Get("user_id"), []string{"9"}) {
		t.Errorf("outgoing user_id metadata = %v, want [9]", md.Get("user_id"))
	}
}

func TestWithUserIDMetadataLeavesContextWithoutIdentityUnchanged(t *testing.T) {
	ctx := context.Background()
	got := withUserIDMetadata(ctx)
	if got != ctx {
		t.Error("context without user ID was replaced")
	}
	if _, ok := metadata.FromOutgoingContext(got); ok {
		t.Error("context without user ID unexpectedly has outgoing metadata")
	}
}
