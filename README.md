![Code coverage](https://img.shields.io/badge/Coverage-100%25-green)
# Ginche 
<img src="https://github.com/chloyka/ginche-logo/raw/main/ginche.jpg" width="250px">

## What is Ginche?
Ginche (gin(ca)che) is a simple gin-gonic middleware for caching HTTP responses in memory.
It allows developers to cache Gin HTTP responses in memory, reducing the load on the server and improving response times. 

## Ideology of Ginche
Ginche was created for make this possible:

1. Create fast caches with minimum overhead
2. Define your own cache key generation rules
3. Minimum effort for use in already finished projects
4. Define your own exclusion rules

## Benchmarks
#### Cache implementation (No middleware)
```
$ go test -bench . -benchmem

goos: darwin
goarch: arm64
pkg: github.com/chloyka/ginche

BenchmarkCache_Set-8     1318003               785.2 ns/op           261 B/op         10 allocs/op
BenchmarkCache_Get-8     1000000              1020 ns/op             314 B/op         12 allocs/op
```
#### Middleware
```
$ go test -bench . -benchmem

goos: darwin
goarch: arm64
pkg: github.com/chloyka/ginche

BenchmarkMiddleware-8            1000000              1725 ns/op            1584 B/op         16 allocs/op
```

## Basic usage
```go
package main

import (
    "github.com/chloyka/ginche"
    "github.com/gin-gonic/gin"
)

func main() {
    store := ginche.NewCache()
    r := gin.New()
    r.Use(ginche.Middleware(store, nil))
    r.GET("/ping", func(c *gin.Context) {
        c.String(200, "pong")
    })
}
```

## Examples
See [Full Examples](https://github.com/chloyka/ginche/blob/master/examples)

## FAQ

#### Can i cache responses based on its request payload? 
Yes! You can define your own cache key generation algos. Here is an Example
```go
router.Use(ginche.Middleware(store, &ginche.Options{
    KeyFunc: func(ctx *gin.Context) string {
        if ctx.Request.Method == "POST" {
            // Use the full URL and the name field from the body as the cache key
            var body map[string]interface{}
            err := ctx.BindJSON(&body)
            if err != nil {
                // If there is an error, skip caching
                return ginche.SkipCacheKeyValue
            }
            return ctx.Request.URL.String() + body["name"].(string)
        } else {
            // Otherwise, skip caching
            return ginche.SkipCacheKeyValue
        }
    },
}))
```

## TODO:
1. Implement adapter interface for external storages
2. Implement Redis storage
3. Implement Memcached storage


Feel free to Open issues, requesting features or contributing