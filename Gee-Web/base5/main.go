package main

import (
	"gee"
	"log"
	"net/http"
	"time"
)

// 中间件

func main() {
	r := gee.NewEngine()
	v1 := r.Group("/v1")
	v1.Use(gee.Logger())
	v1.GET("/", func(context *gee.Context) {
		context.HTML(http.StatusOK, "<h1>hello this is v1</h1>")
	})

	v2 := r.Group("/v2")
	v2.Use(ForV2())
	v2.GET("/hello/:name", func(context *gee.Context) {
		context.String(http.StatusOK, "hello %s, you're at %s\n", context.Param("name"), context.Path)
	})

	r.Run(":9999")
}

func ForV2() gee.HandlerFunc {
	return func(context *gee.Context) {
		t := time.Now()
		context.Fail(500, "Internal Server Error")
		log.Printf("[%d] %s in %v for group v2", context.StatusCode, context.Req.RequestURI, time.Since(t))
	}
}
