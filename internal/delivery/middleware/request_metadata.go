package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/victorotene80/authentication_api/internal/delivery/http/requestctx"
)

func RequestMetadata() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		ua := c.Request.UserAgent()
		deviceID := c.GetHeader("X-Device-ID")

		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.NewString()
		}

		meta := requestctx.Metadata{
			IP:        ip,
			UserAgent: ua,
			DeviceID:  deviceID,
			RequestID: requestID,
		}

		requestctx.WithMetadata(c, meta)

		c.Next()
	}
}