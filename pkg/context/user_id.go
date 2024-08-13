package context

import "context"

type contextKey int

var userIdKey contextKey

func SetUserID(ctx context.Context, userId int) context.Context {
	return context.WithValue(ctx, userIdKey, userId)
}

func GetUserID(ctx context.Context) int {
	return ctx.Value(userIdKey).(int)
}
