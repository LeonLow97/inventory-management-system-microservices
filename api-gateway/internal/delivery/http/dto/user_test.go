package dto

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/LeonLow97/internal/core/domain"
)

func TestFromLoginRequest(t *testing.T) {
	got := FromLoginRequest(&LoginRequest{Email: "jane@example.com", Password: "Secret123!"})
	want := domain.User{Email: "jane@example.com", Password: "Secret123!"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("FromLoginRequest() = %#v, want %#v", got, want)
	}
}

func TestToLoginResponse(t *testing.T) {
	got := ToLoginResponse(&domain.User{
		FirstName: "Jane",
		LastName:  "Doe",
		Email:     "jane@example.com",
		Active:    true,
		Admin:     true,
	})
	want := LoginResponse{
		FirstName: "Jane",
		LastName:  "Doe",
		Email:     "jane@example.com",
		Active:    true,
		Admin:     true,
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ToLoginResponse() = %#v, want %#v", got, want)
	}

	body, err := json.Marshal(got)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}
	if string(body) != `{"first_name":"Jane","last_name":"Doe","email":"jane@example.com","admin":true}` {
		t.Fatalf("login response JSON = %s, want sensitive active field omitted", body)
	}
}

func TestFromSignUpRequest(t *testing.T) {
	got := FromSignUpRequest(&SignUpRequest{
		FirstName: "Jane",
		LastName:  "Doe",
		Password:  "Secret123!",
		Email:     "jane@example.com",
	})
	want := domain.User{
		FirstName: "Jane",
		LastName:  "Doe",
		Password:  "Secret123!",
		Email:     "jane@example.com",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("FromSignUpRequest() = %#v, want %#v", got, want)
	}
}

func TestToGetUsersResponse(t *testing.T) {
	got := ToGetUsersResponse([]domain.User{
		{ID: 1, FirstName: "Jane", LastName: "Doe", Email: "jane@example.com", Active: true, Admin: true},
		{ID: 2, FirstName: "John", LastName: "Smith", Email: "john@example.com"},
	}, "next-page")
	want := GetUsersResponse{
		Users: []User{
			{ID: 1, FirstName: "Jane", LastName: "Doe", Email: "jane@example.com", Active: true, Admin: true},
			{ID: 2, FirstName: "John", LastName: "Smith", Email: "john@example.com"},
		},
		NextCursor: "next-page",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ToGetUsersResponse() = %#v, want %#v", got, want)
	}
}

func TestFromUpdateUserRequest(t *testing.T) {
	got := FromUpdateUserRequest(&UpdateUserRequest{
		FirstName: "Jane",
		LastName:  "Doe",
		Password:  "Secret123!",
	})
	want := domain.User{FirstName: "Jane", LastName: "Doe", Password: "Secret123!"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("FromUpdateUserRequest() = %#v, want %#v", got, want)
	}
}
