package gee

import (
	"log"
	"net/http"
)

type Router struct {
	handlers map[string]HandlerFunc
}

func NewRouter() *Router {
	return &Router{
		handlers: make(map[string]HandlerFunc),
	}
}

func (router *Router) addRoute(method, pattern string, handler HandlerFunc) {
	log.Printf("Route %4s - %s\n", method, pattern)
	key := method + "-" + pattern
	router.handlers[key] = handler
}

func (router *Router) Handle(context *Context) {
	key := context.Method + "-" + context.Path
	if handler, ok := router.handlers[key]; ok {
		handler(context)
	} else {
		context.String(http.StatusNotFound, "404 not found: %s\n", context.Path)
	}
}
