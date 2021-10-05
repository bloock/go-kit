package cache

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"time"

	"github.com/gin-gonic/gin"
)

const PREFIX = "gin:cache:"

func md5String(url string) string {
	h := md5.New()
	io.WriteString(h, url)
	return hex.EncodeToString(h.Sum(nil))
}

func Middleware(cache Cache, duration time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		var key string
		uri := c.Request.RequestURI
		rawQuery := c.Request.URL.RawQuery
		key = PREFIX + md5String(fmt.Sprintf("%s%s", uri, rawQuery))
		CacheRequest(c, key, cache, duration)
	}
}
