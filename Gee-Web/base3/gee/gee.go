package gee

import (
	"net/http"
)

type Engine struct {
	router *router
}

func NewEngine() *Engine {
	return &Engine{
		router: newRouter(),
	}
}

func (engine *Engine) AddRoute(method, pattern string, handler HandlerFunc) {
	engine.router.addRoute(method, pattern, handler)
}

func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.AddRoute("GET", pattern, handler)
}

func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.AddRoute("POST", pattern, handler)
}

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	context := NewContext(w, req)
	engine.router.handle(context)
}
