package middleware

import (
	"bytes"
	"github.com/chloyka/ginche/cache"
	"github.com/gin-gonic/gin"
)

type writer struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *writer) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func Serve(storage *cache.Cache, options *Options) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var key string
		if options != nil {
			if options.KeyFunc != nil {
				key = options.KeyFunc(ctx)
				if key == "" {
					ctx.Next()
					return
				}
			}
			if options.ExcludeMethods != nil && sliceContainsString(options.ExcludeMethods, ctx.Request.Method) {
				ctx.Next()
				return
			}
		} else {
			key = ctx.Request.URL.Path
		}

		if data, ok := storage.Get(key); ok {
			entry := data.(*httpCacheItem)
			ctx.JSON(entry.Status, entry.Data)
			return
		}
		w := &writer{body: &bytes.Buffer{}, ResponseWriter: ctx.Writer}
		ctx.Writer = w
		ctx.Next()

		if options != nil && sliceContainsInt(options.ExcludeStatuses, ctx.Writer.Status()) {
			return
		}
		storage.Set(&key, &httpCacheItem{Status: ctx.Writer.Status(), Data: w.body.String()})
	}
}

type Options struct {
	KeyFunc         func(c *gin.Context) string
	ExcludeStatuses []int
	ExcludeMethods  []string
}

type httpCacheItem struct {
	Status int
	Data   interface{}
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
