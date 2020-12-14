package main

import (
	"fmt"
	"net/http"
)

func handerFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Hello, 这里是goblog</h1>")
}

func main() {
	http.HandleFunc("/", handerFunc)
	http.ListenAndServe(":3000", nil)
}
