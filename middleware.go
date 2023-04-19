package ginche

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	// CTXSkipCacheKey is the key used to skip the cache for a request
	// ctx.Set(CTXSkipCacheKey, CTXSkipCacheValue)
	CTXSkipCacheKey = "cache-key"
	// CTXSkipCacheValue is the value used to skip the cache for a request
	// ctx.Set(CTXSkipCacheKey, CTXSkipCacheValue)
	CTXSkipCacheValue = ""
	// SkipCacheKeyValue is the value used to skip the cache for a request
	SkipCacheKeyValue = ""
	// HeaderXCache is the header key used to indicate the cache status
	HeaderXCache = "X-Cache"
	// HeaderXCacheHit is the header value used to indicate the cache status HIT
	HeaderXCacheHit = "HIT"
	// HeaderXCacheSkip is the header value used to indicate the cache status SKIP
	HeaderXCacheSkip = "SKIP"
	// HeaderXCacheMiss is the header value used to indicate the cache status MISS
	HeaderXCacheMiss = "MISS"
)

type writer struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *writer) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// Middleware is the gin-gonic middleware function
// it will cache the response if the request is not excluded
// and if the response is not excluded
func Middleware(storage CacheAdapter, options *Options) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set(HeaderXCache, HeaderXCacheSkip)
		cacheKey := ctx.Request.URL.Path
		if options != nil {
			if options.KeyFunc != nil {
				cacheKey = options.KeyFunc(ctx)
				if cacheKey == CTXSkipCacheValue {
					ctx.Next()
					return
				}
			}
			if options.ExcludeMethods != nil && sliceContainsString(options.ExcludeMethods, ctx.Request.Method) {
				ctx.Next()
				return
			}
			if options.ExcludePaths != nil && sliceContainsString(options.ExcludePaths, strings.TrimRight(ctx.Request.URL.Path, "/")) {
				ctx.Next()
				return
			}
		}

		if data, ok := storage.Get(cacheKey); ok {
			entry := data.(*httpCacheItem)
			for k, h := range entry.Headers {
				for _, v := range h {
					ctx.Writer.Header().Add(k, v)
				}
			}
			ctx.Writer.Header().Set(HeaderXCache, HeaderXCacheHit)
			ctx.String(entry.Status, entry.Data.(string))
			ctx.Abort()
			return
		}
		w := &writer{body: &bytes.Buffer{}, ResponseWriter: ctx.Writer}
		ctx.Writer = w
		ctx.Next()
		if options != nil && sliceContainsInt(options.ExcludeStatuses, ctx.Writer.Status()) {
			ctx.Abort()
			return
		}
		k, skip := ctx.Get(CTXSkipCacheKey)
		if k == CTXSkipCacheValue && skip {
			ctx.Abort()
			return
		}
		ctx.Writer.Header().Set(HeaderXCache, HeaderXCacheMiss)
		storage.Set(&cacheKey, &httpCacheItem{Status: ctx.Writer.Status(), Data: w.body.String(), Headers: w.Header().Clone()})
	}
}

// Options is the options for the middleware
// KeyFunc is the function used to generate the cache key
// ExcludeStatuses is the list of status codes to exclude from the cache
// ExcludeMethods is the list of methods to exclude from the cache
// ExcludePaths is the list of paths to exclude from the cache
type Options struct {
	KeyFunc         func(c *gin.Context) string
	ExcludeStatuses []int
	ExcludeMethods  []string
	ExcludePaths    []string
}

type httpCacheItem struct {
	Status  int
	Headers http.Header
	Data    interface{}
}

func sliceContainsInt(arr []int, ele int) bool {
	for _, e := range arr {
		if e == ele {
			return true
		}
	}
	return false
}

func sliceContainsString(arr []string, ele string) bool {
	for _, e := range arr {
		if e == ele {
			return true
		}
	}
	return false
}
