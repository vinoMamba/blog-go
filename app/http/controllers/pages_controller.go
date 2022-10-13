package controllers

import (
	"fmt"
	"net/http"
)

type PagesController struct {
}

func (*PagesController) Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Hello, 欢迎来到 goblog！</h1>")
}

func (*PagesController) About(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>关于我们</h1>")
}

func (*PagesController) NotFound(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>404 页面未找到</h1>")
}
