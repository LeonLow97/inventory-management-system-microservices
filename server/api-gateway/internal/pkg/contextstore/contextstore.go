package contextstore

import (
	"context"
	"log"
)

type contextKeyInt int

const (
	contextKeyUserID contextKeyInt = iota
)

func InjectUserIDIntoContext(ctx context.Context, userID int) context.Context {
	return context.WithValue(ctx, contextKeyUserID, userID)
}

func UserIDFromContext(ctx context.Context) (int, error) {
	v, found := ctx.Value(contextKeyUserID).(int)
	if !found {
		log.Println("User ID not found in context")
		return 0, ErrUserIDNotInContext
	}
	return v, nil
}
