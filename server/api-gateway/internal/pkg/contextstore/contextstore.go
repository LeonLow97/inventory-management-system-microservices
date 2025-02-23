package contextstore

import "context"

type contextKeyInt int

const (
	contextKeyUserID contextKeyInt = iota
)

func InjectUserIDIntoContext(ctx context.Context, userID int) context.Context {
	return context.WithValue(ctx, contextKeyUserID, userID)
}

func UserIDFromContext(ctx context.Context) (int, error) {
	if v, ok := ctx.Value(contextKeyUserID).(int); ok {
		return v, nil
	}
	return 0, ErrUserIDNotInContext
}
