package handlers

import "context"

type key string

const userKey = key("5276639454010324953")

//NewContext with the user stored
func NewContext(ctx context.Context, u uint) context.Context {
	return context.WithValue(ctx, userKey, u)
}

// FromContext returns the User value stored in ctx, if any.
func FromContext(ctx context.Context) (uint, bool) {
	u, ok := ctx.Value(userKey).(uint)
	return u, ok
}
