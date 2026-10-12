package contextstore

import (
	"context"
	"testing"
)

func TestUserIDFromContext(t *testing.T) {
	testCases := []struct {
		name      string
		ctx       context.Context
		wantUser  int64
		wantFound bool
	}{
		{
			name:      "user ID is present",
			ctx:       WithUserID(context.Background(), 42),
			wantUser:  42,
			wantFound: true,
		},
		{
			name:      "user ID is absent",
			ctx:       context.Background(),
			wantFound: false,
		},
		{
			name:      "zero user ID is still present",
			ctx:       WithUserID(context.Background(), 0),
			wantUser:  0,
			wantFound: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			gotUser, gotFound := UserIDFromContext(testCase.ctx)
			if gotUser != testCase.wantUser || gotFound != testCase.wantFound {
				t.Fatalf("UserIDFromContext() = (%d, %t), want (%d, %t)", gotUser, gotFound, testCase.wantUser, testCase.wantFound)
			}
		})
	}
}
