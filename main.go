package main

import (
	"fmt"
	"net/http"
)

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}
func aboutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "About Page")
}
func main() {
	http.HandleFunc("/", defaultHandler)
	http.HandleFunc("/about", aboutHandler)
	http.ListenAndServe(":3000", nil)
}
