package main

import (
	"fmt"
	"net/http"
)

func defaultHander(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if r.URL.Path == "/" {
		fmt.Fprint(w, "<h1>Hello, 这里是goblog</h1>")
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "_我不知道你在说什么_")
	}
}

func aboutHander(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "这里是关于blog的内容12123123"+
		"<a href=\"mailto:summer@example.com\">summer@example.com</a>")
}

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/", defaultHander)
	router.HandleFunc("/about", aboutHander)
	http.ListenAndServe(":3000", router)
}
