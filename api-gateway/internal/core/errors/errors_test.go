package coreerrors

import (
	"errors"
	"testing"
)

func TestError_Error(t *testing.T) {
	cause := errors.New("database unavailable")
	testCases := []struct {
		name string
		err  *Error
		want string
	}{
		{
			name: "message takes precedence over cause",
			err:  New(KindInternal, "request failed", cause),
			want: "request failed",
		},
		{
			name: "cause is used when message is empty",
			err:  New(KindInternal, "", cause),
			want: "database unavailable",
		},
		{
			name: "kind is used when message and cause are empty",
			err:  New(KindNotFound, "", nil),
			want: string(KindNotFound),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if got := testCase.err.Error(); got != testCase.want {
				t.Fatalf("Error() = %q, want %q", got, testCase.want)
			}
		})
	}
}

func TestError_Unwrap(t *testing.T) {
	cause := errors.New("original failure")
	appErr := New(KindInternal, "operation failed", cause)

	if got := errors.Unwrap(appErr); got != cause {
		t.Fatalf("Unwrap() = %v, want original cause", got)
	}
}
