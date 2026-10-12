package contextstore

import (
	"context"
	"testing"

	"google.golang.org/grpc/metadata"
)

func TestUserIDFromContext(t *testing.T) {
	testCases := []struct {
		name     string
		metadata metadata.MD
		wantID   int64
		wantErr  bool
	}{
		{name: "reads user ID", metadata: metadata.Pairs("user_id", "42"), wantID: 42},
		{name: "reports missing metadata", wantErr: true},
		{name: "reports empty user ID", metadata: metadata.Pairs("user_id", ""), wantErr: true},
		{name: "reports malformed user ID", metadata: metadata.Pairs("user_id", "not-a-number"), wantErr: true},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctx := context.Background()
			if testCase.metadata != nil {
				ctx = metadata.NewIncomingContext(ctx, testCase.metadata)
			}

			gotID, err := UserIDFromContext(ctx)
			if (err != nil) != testCase.wantErr {
				t.Fatalf("UserIDFromContext() error = %v, wantErr %t", err, testCase.wantErr)
			}
			if err == nil && gotID != testCase.wantID {
				t.Errorf("UserIDFromContext() = %d, want %d", gotID, testCase.wantID)
			}
		})
	}
}
