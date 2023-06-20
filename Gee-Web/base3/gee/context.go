package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	Writer     http.ResponseWriter
	Req        *http.Request
	Path       string
	Method     string
	Params     map[string]string
	StatusCode int
}

func (context *Context) Param(key string) string {
	value, _ := context.Params[key]
	return value
}

func NewContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
	}
}

func (context *Context) PostForm(key string) string {
	return context.Req.FormValue(key)
}

func (context *Context) Query(key string) string {
	return context.Req.URL.Query().Get(key)
}

func (context *Context) Status(code int) {
	context.StatusCode = code
	context.Writer.WriteHeader(code)
}

func (context *Context) SetHeader(key string, value string) {
	context.Writer.Header().Set(key, value)
}

func (context *Context) String(code int, format string, values ...interface{}) {
	context.SetHeader("Content-Type", "text/plain")
	context.Status(code)
	context.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (context *Context) Json(code int, obj interface{}) {
	context.SetHeader("Content-Type", "application/json")
	context.Status(200)
	encoder := json.NewEncoder(context.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(context.Writer, err.Error(), 500)
	}
}

func (context *Context) Data(code int, data []byte) {
	context.Status(code)
	context.Writer.Write(data)
}

func (context *Context) HTML(code int, html string) {
	context.SetHeader("Content-Type", "text/html")
	context.Status(code)
	context.Writer.Write([]byte(html))
}
