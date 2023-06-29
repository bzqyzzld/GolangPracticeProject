package main

import (
	"fmt"
	"gee"
	"html/template"
	"net/http"
	"time"
)

// 模板，类似与django的模板引擎类似的效果

func main() {
	r := gee.NewEngine()
	r.Static("/assets", "./static")
	r.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})

	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./static")

	stu1 := &Student{Name: "Jack", Age: 22}
	stu2 := &Student{Name: "Micheal", Age: 33}

	r.GET("/", func(context *gee.Context) {
		context.HTML(http.StatusOK, "css.tmpl", nil)
	})

	r.GET("/students", func(context *gee.Context) {
		context.HTML(http.StatusOK, "arr.tmpl", gee.H{
			"title":  "gee",
			"stuArr": [2]*Student{stu1, stu2},
		})
	})

	r.GET("/date", func(context *gee.Context) {
		context.HTML(http.StatusOK, "custom_func.tmpl", gee.H{
			"title": "gee",
			"now":   time.Date(2019, 8, 7, 0, 0, 0, 0, time.UTC),
		})
	})

	r.Run(":9999")
}

func FormatAsDate(t time.Time) string {
	y, m, d := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", y, m, d)
}

type Student struct {
	Name string
	Age  int
}
