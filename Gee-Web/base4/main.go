package main

import (
	"gee"
	"net/http"
)

func main() {
	r := gee.NewEngine()
	v1 := r.Group("/v1")
	v1.GET("/", func(context *gee.Context) {
		context.HTML(http.StatusOK, "<h1>Hello gee</h1>")
	})

	r.Run(":9999")
}
