package main

import (
	"gee"
	"net/http"
)

func main() {
	r := gee.NewEngine()
	v1 := r.Group("/v1")
	v2 := v1.Group("/test")
	v2.GET("/", func(context *gee.Context) {
		context.HTML(http.StatusOK, "<h1>Hello gee</h1>")
	})

	r.Run(":9999")
}
