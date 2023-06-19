package main

import (
	"gee"
	"net/http"
)

func main() {
	r := gee.NewEngine()
	r.GET("/", func(context *gee.Context) {
		context.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	})

	r.GET("/hello", func(context *gee.Context) {
		context.String(http.StatusOK, "hello %s, you're at %s\n", context.Query("name"), context.Path)
	})

	r.GET("/hello/:name", func(context *gee.Context) {
		context.String(http.StatusOK, "hello %s, you're at %s\n", context.Param("name"), context.Path)
	})

	r.GET("/assets/*filepath", func(context *gee.Context) {
		context.Json(http.StatusOK, gee.H{"filepath": context.Param("filepath")})
	})

	r.Run(":9999")
}
