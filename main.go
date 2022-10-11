package main

import (
	"fmt"
	"net/http"
)

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if r.URL.Path == "/" {
		fmt.Fprint(w, "<h1>Hello, 这里是 goblo</h1>")
		fmt.Fprint(w, "请求的路径为:"+r.URL.Path)
	} else if r.URL.Path == "/about" {
		fmt.Fprint(w, "<h1>Hello, 这里是关于页面</h1>")
		fmt.Fprint(w, "请求的路径为:"+r.URL.Path)
	} else {
		fmt.Fprint(w, "404")
	}
}
func main() {
	http.HandleFunc("/", handlerFunc)
	http.ListenAndServe(":3000", nil)
}
