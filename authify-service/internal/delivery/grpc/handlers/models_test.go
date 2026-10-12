package handlers

import (
	"reflect"
	"testing"

	pb "github.com/LeonLow97/authify/grpc-proto/pb"
	"github.com/LeonLow97/internal/core/domain"
)

func TestToLoginInput(t *testing.T) {
	got := toLoginInput(&pb.LoginRequest{Email: "  ada@example.com  ", Password: " Secret123! "})
	want := domain.LoginInput{Email: "ada@example.com", Password: "Secret123!"}
	if got != want {
		t.Errorf("toLoginInput() = %#v, want %#v", got, want)
	}
}

func TestToLoginResponse(t *testing.T) {
	user := &domain.User{
		Email:     "ada@example.com",
		FirstName: stringPointerForHandlerTest("Ada"),
		LastName:  stringPointerForHandlerTest("Lovelace"),
		Active:    true,
		Admin:     true,
	}
	want := &pb.LoginResponse{
		Email:     "ada@example.com",
		FirstName: "Ada",
		LastName:  "Lovelace",
		Active:    true,
		Admin:     true,
		Token:     "signed-token",
	}
	if got := toLoginResponse(user, "signed-token"); !reflect.DeepEqual(got, want) {
		t.Errorf("toLoginResponse() = %#v, want %#v", got, want)
	}
}

func TestToSignUpInput(t *testing.T) {
	got := toSignUpInput(&pb.SignUpRequest{
		Email:     " ada@example.com ",
		Password:  " Secret123! ",
		FirstName: " Ada ",
		LastName:  " Lovelace ",
	})
	want := domain.SignUpInput{
		Email:     "ada@example.com",
		Password:  "Secret123!",
		FirstName: "Ada",
		LastName:  "Lovelace",
	}
	if got != want {
		t.Errorf("toSignUpInput() = %#v, want %#v", got, want)
	}
}

func TestToGetUsersResponse(t *testing.T) {
	users := []domain.User{{
		ID:        17,
		Email:     "ada@example.com",
		FirstName: stringPointerForHandlerTest("Ada"),
		LastName:  stringPointerForHandlerTest("Lovelace"),
		Active:    true,
		Admin:     true,
	}}
	want := &pb.GetUsersResponse{
		Users: []*pb.User{{
			ID:        17,
			Email:     "ada@example.com",
			FirstName: "Ada",
			LastName:  "Lovelace",
			Active:    true,
		}},
		NextCursor: "next-page",
	}
	if got := toGetUsersResponse(users, "next-page"); !reflect.DeepEqual(got, want) {
		t.Errorf("toGetUsersResponse() = %#v, want %#v", got, want)
	}
}

func TestToUpdateUserInput(t *testing.T) {
	got := ToUpdateUserInput(&pb.UpdateUserRequest{FirstName: "Ada", LastName: "Lovelace", Password: "Secret123!"})
	if got.FirstName == nil || *got.FirstName != "Ada" || got.LastName == nil || *got.LastName != "Lovelace" || got.Password == nil || *got.Password != "Secret123!" {
		t.Errorf("ToUpdateUserInput() = %#v, want all request fields copied", got)
	}
}

func TestSanitizePaginationCursor(t *testing.T) {
	testCases := []struct {
		name  string
		limit int64
		want  int64
	}{
		{name: "negative limit defaults", limit: -1, want: 10},
		{name: "zero limit defaults", limit: 0, want: 10},
		{name: "positive limit is preserved", limit: 20, want: 20},
		{name: "maximum limit is preserved", limit: 50, want: 50},
		{name: "large limit is capped", limit: 100, want: 50},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if got := SanitizePaginationCursor(testCase.limit); got != testCase.want {
				t.Errorf("SanitizePaginationCursor(%d) = %d, want %d", testCase.limit, got, testCase.want)
			}
		})
	}
}

func stringPointerForHandlerTest(value string) *string {
	return &value
}
