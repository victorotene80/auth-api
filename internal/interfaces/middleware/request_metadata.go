package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"net"
	"net/http"

	chiMw "github.com/go-chi/chi/v5/middleware"

	"github.com/victorotene80/authentication_api/internal/interfaces/http/requestctx"
)

func RequestMetadata(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ua := r.UserAgent()
		deviceID := r.Header.Get("X-Device-ID")
		acceptLang := r.Header.Get("Accept-Language")

		ip := clientIP(r)

		fpRaw := ua + "|" + deviceID + "|" + acceptLang
		fpHash := sha256.Sum256([]byte(fpRaw))
		fingerprint := hex.EncodeToString(fpHash[:])

		deviceName := ua

		reqID := chiMw.GetReqID(r.Context())
		if reqID == "" {
			reqID = r.Header.Get("X-Request-ID")
		}

		meta := requestctx.RequestMeta{
			IPAddress:         ip,
			UserAgent:         ua,
			DeviceID:          deviceID,
			DeviceFingerprint: fingerprint,
			DeviceName:        deviceName,
			RequestID:         reqID,
		}

		ctx := requestctx.WithMeta(r.Context(), meta)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func clientIP(r *http.Request) string {
	// If you're behind a reverse proxy and trust X-Forwarded-For,
	// you can uncomment this logic:
	//
	// if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
	//     parts := strings.Split(xff, ",")
	//     return strings.TrimSpace(parts[0])
	// }
	//change this to X-Real-IP if using that header instead
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		// fallback: r.RemoteAddr as-is
		return r.RemoteAddr
	}
	return host
}
