package main

import (
	"fmt"
	"net/http"
	"strings"
)

var chttp = http.NewServeMux()

func main() {

	chttp.Handle("/", http.FileServer(http.Dir("./public")))

	http.HandleFunc("/", HomeHandler) // homepage
	http.ListenAndServe(":8080", nil)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {

	if strings.Contains(r.URL.Path, ".") {
		chttp.ServeHTTP(w, r)
	} else {
		fmt.Fprintf(w, "HomeHandler")
	}
}
