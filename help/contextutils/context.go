package contextutils

import "context"

func GetStringValue(ctx context.Context, key string) (string, bool) {
	val := ctx.Value(key)
	if val == nil {
		return "", false
	}
	v, ok := val.(string)
	return v, ok
}
