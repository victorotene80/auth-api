package middleware

import (
    "net/http"

    "github.com/victorotene80/authentication_api/internal/interfaces/http/requestctx"
)

func RequestMetadata(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        meta := requestctx.RequestMeta{
            IPAddress: r.RemoteAddr,
            UserAgent: r.UserAgent(),
            DeviceID:  r.Header.Get("X-Device-ID"),
            RequestID: r.Header.Get("X-Request-ID"),
        }
        ctx := requestctx.WithMeta(r.Context(), meta)
        r = r.WithContext(ctx)
        next.ServeHTTP(w, r)
    })
}