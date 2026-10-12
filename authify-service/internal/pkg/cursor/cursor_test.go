package cursor

import (
	"encoding/base64"
	"testing"

	"github.com/LeonLow97/internal/core/domain"
)

func TestCursorRoundTrip(t *testing.T) {
	testCases := []struct {
		name   string
		cursor domain.UserCursor
	}{
		{name: "zero cursor", cursor: domain.UserCursor{}},
		{name: "user cursor", cursor: domain.UserCursor{ID: 52}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			encoded, err := EncodeCursor(testCase.cursor)
			if err != nil {
				t.Fatalf("EncodeCursor() error = %v", err)
			}
			var got domain.UserCursor
			if err := DecodeCursor(encoded, &got); err != nil {
				t.Fatalf("DecodeCursor() error = %v", err)
			}
			if got != testCase.cursor {
				t.Errorf("decoded cursor = %#v, want %#v", got, testCase.cursor)
			}
		})
	}
}

func TestDecodeCursor(t *testing.T) {
	invalidJSON := base64.StdEncoding.EncodeToString([]byte("{"))
	testCases := []struct {
		name      string
		value     string
		wantID    int64
		wantError bool
	}{
		{name: "empty cursor leaves zero value", wantID: 0},
		{name: "valid cursor", value: base64.StdEncoding.EncodeToString([]byte(`{"id":18}`)), wantID: 18},
		{name: "invalid base64", value: "%", wantError: true},
		{name: "invalid json", value: invalidJSON, wantError: true},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var got domain.UserCursor
			err := DecodeCursor(testCase.value, &got)
			if (err != nil) != testCase.wantError {
				t.Fatalf("DecodeCursor() error = %v, wantErr %t", err, testCase.wantError)
			}
			if err == nil && got.ID != testCase.wantID {
				t.Errorf("decoded cursor ID = %d, want %d", got.ID, testCase.wantID)
			}
		})
	}
}

func TestEncodeCursorRejectsUnmarshalableValue(t *testing.T) {
	if _, err := EncodeCursor(make(chan int)); err == nil {
		t.Fatal("EncodeCursor() error = nil, want error for an unsupported JSON value")
	}
}
