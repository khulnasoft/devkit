package winlayers

import "context"

type contextKeyT string

var contextKey = contextKeyT("devkit/winlayers-on")

func UseWindowsLayerMode(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextKey, true)
}

func hasWindowsLayerMode(ctx context.Context) bool {
	v := ctx.Value(contextKey)
	return v != nil
}
