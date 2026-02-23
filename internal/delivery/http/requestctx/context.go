package requestctx

import (
	"github.com/gin-gonic/gin"
)

const metaKey = "request_metadata"

type Metadata struct {
	IP        string
	UserAgent string
	DeviceID  string
	RequestID string
}

func WithMetadata(c *gin.Context, meta Metadata) {
	c.Set(metaKey, meta)
}

func GetMetadata(c *gin.Context) Metadata {
	value, exists := c.Get(metaKey)
	if !exists {
		return Metadata{}
	}

	if meta, ok := value.(Metadata); ok {
		return meta
	}

	return Metadata{}
}