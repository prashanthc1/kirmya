package common

import "context"

type contextKey string

const authUserContextKey contextKey = "authUser"

type AuthUser struct {
	ID    string
	Email string
	Role  string
}

func ContextWithAuthUser(ctx context.Context, user AuthUser) context.Context {
	return context.WithValue(ctx, authUserContextKey, user)
}

func ContextWithUserID(ctx context.Context, userID string) context.Context {
	return ContextWithAuthUser(ctx, AuthUser{ID: userID})
}

func AuthUserFromContext(ctx context.Context) AuthUser {
	value, ok := ctx.Value(authUserContextKey).(AuthUser)
	if !ok {
		return AuthUser{}
	}
	return value
}

func UserIDFromContext(ctx context.Context) string {
	return AuthUserFromContext(ctx).ID
}

func UserEmailFromContext(ctx context.Context) string {
	return AuthUserFromContext(ctx).Email
}

func UserRoleFromContext(ctx context.Context) string {
	return AuthUserFromContext(ctx).Role
}
