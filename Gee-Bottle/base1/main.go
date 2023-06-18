package main

import (
	"fmt"
	"net/http"
)

type Engine struct{}

func (engine *Engine) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/":
		fmt.Fprintf(writer, "this is index url\n")
	case "/hello":
		fmt.Fprintf(writer, "this is index url")
	default:
		fmt.Fprintf(writer, "404 url not found\n")
	}

}

func main() {
	engine := new(Engine)
	http.ListenAndServe(":9999", engine)
}

func indexHandler(writer http.ResponseWriter, request *http.Request) {
	// 写入的时候不可以使用log.Fatal不然会报错
	fmt.Fprintf(writer, "URL.Path = %q\n", request.URL.Path)
}

func helloHandler(writer http.ResponseWriter, request *http.Request) {
	for k, v := range request.Header {
		fmt.Fprintf(writer, "Header[%q] = %q\n", k, v)
	}
}
