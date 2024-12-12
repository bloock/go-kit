package gin

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/bloock/go-kit/cache"
	"github.com/gin-gonic/gin"
)

const PREFIX = "gin:cache:"

func md5String(url string) string {
	h := md5.New()
	io.WriteString(h, url)
	return hex.EncodeToString(h.Sum(nil))
}

func Middleware(cache cache.Cache, duration time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		var key string
		uri := c.Request.RequestURI
		rawQuery := c.Request.URL.RawQuery
		key = PREFIX + md5String(fmt.Sprintf("%s%s", uri, rawQuery))
		CacheRequest(c, key, cache, duration)
	}
}

type InvalidateCacheFunc func(uri string) gin.HandlerFunc

func InvalidateMiddleware(cache cache.Cache) InvalidateCacheFunc {
	return func(uri string) gin.HandlerFunc {
		return func(c *gin.Context) {
			var key, newUri string
			userID := c.GetHeader("X-User-ID")
			newUri = strings.Replace(uri, ":user_id", userID, 1)
			rawQuery := c.Request.URL.RawQuery
			key = PREFIX + md5String(fmt.Sprintf("%s%s", newUri, rawQuery))
			CacheInvalidate(c, key, cache)
		}
	}
}
