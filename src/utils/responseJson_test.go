package utils

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_ResponseJson(t *testing.T) {
	responseJsonTests := []struct {
		testName string
		code int
		message string
		expectedJsonResponse string
	} {
		{"Successful Login", http.StatusOK, "Successfully Logged In!", `{"code":200,"message":"Successfully Logged In!"}`},
		{"Duplicate username", http.StatusBadRequest, "Username already exists. Please try again.", `{"code":400,"message":"Username already exists. Please try again."}`},
	}

	for _, e := range responseJsonTests {
		// Creating a mock ResponseWriter
		w := httptest.NewRecorder()

		ResponseJson(w, e.code, e.message)

		// Read the response body as a string
		body, _ := io.ReadAll(w.Result().Body)
		actual := string(body)

		expected := e.expectedJsonResponse
		if actual != expected {
			t.Errorf("%s: expected %s but got %s", e.testName, e.expectedJsonResponse, actual)
		}
	}
}