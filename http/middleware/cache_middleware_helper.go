package middleware

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bloock/go-kit/cache"
	"github.com/gin-gonic/gin"
)

type Cached struct {
	Status int
	Body   []byte
	Header http.Header
}

func CacheInvalidate(c *gin.Context, key string, cache cache.Cache) {
	err := cache.Del(c, key)
	if err != nil {
		log.Printf("Conection err with redis %s", err.Error())
		c.Next()
		return
	}

	log.Printf("Invalidate cache from key %s", key)
}

func CacheRequest(c *gin.Context, key string, cache cache.Cache, duration time.Duration) {

	ca, err := cache.Get(c, key)
	if err != nil {
		c.Next()
		return
	}

	if ca == nil {
		log.Printf("Caching from key %s with duration %s", key, duration.String())
		writer := c.Writer
		rw := cacheWrappedWriter{ResponseWriter: c.Writer}
		c.Writer = &rw
		c.Next()
		c.Writer = writer

		cache.Set(
			c,
			key,
			encodeCache(&Cached{
				Status: rw.Status(),
				Body:   rw.body.Bytes(),
				Header: http.Header(rw.Header()),
			}),
			duration,
		)

	} else {
		log.Printf("Retriving cache from key %s with ttl %s", key, duration.String())
		cach := decodeCache(ca)
		start := time.Now()
		c.Writer.WriteHeader(cach.Status)
		for k, val := range cach.Header {
			for _, v := range val {
				c.Writer.Header().Add(k, v)
			}
		}
		c.Writer.Header().Add("X-Cache", fmt.Sprintf("%f ms", time.Now().Sub(start).Seconds()*1000))
		c.Writer.Write(cach.Body)
		c.Abort()
	}

}

type cacheWrappedWriter struct {
	gin.ResponseWriter
	body bytes.Buffer
}

func (rw *cacheWrappedWriter) Write(body []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(body)
	if err == nil {
		rw.body.Write(body)
	}
	return n, err
}

func encodeCache(cache *Cached) []byte {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	enc.Encode(cache)

	return b.Bytes()
}

func decodeCache(b []byte) Cached {
	var cached *Cached
	buff := bytes.NewBuffer(b)
	dec := gob.NewDecoder(buff)
	dec.Decode(&cached)

	return *cached
}
