package utils

import "testing"

func TestPointerHelpers(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		want := "value"
		got := FromPointer(ToPointer(want))
		if got != want {
			t.Errorf("FromPointer(ToPointer(%q)) = %q, want %q", want, got, want)
		}
	})

	t.Run("integer", func(t *testing.T) {
		want := 42
		got := FromPointer(ToPointer(want))
		if got != want {
			t.Errorf("FromPointer(ToPointer(%d)) = %d, want %d", want, got, want)
		}
	})
}
