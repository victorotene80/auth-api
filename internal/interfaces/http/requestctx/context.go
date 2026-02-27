package requestctx

import "context"

type key int

const metaKey key = iota

type RequestMeta struct {
    IPAddress string
    UserAgent string
    DeviceID  string
    RequestID string
}

func WithMeta(ctx context.Context, meta RequestMeta) context.Context {
    return context.WithValue(ctx, metaKey, meta)
}

func MetaFrom(ctx context.Context) (RequestMeta, bool) {
    meta, ok := ctx.Value(metaKey).(RequestMeta)
    return meta, ok
}