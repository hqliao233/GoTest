package main

import (
	"fmt"
	"net/http"
)

func handerFunc(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		fmt.Fprint(w, "<h1>Hello, 这里是goblog</h1>")
	} else if r.URL.Path == "/about" {
		fmt.Fprint(w, "这里是关于blog的内容")
	} else {
		fmt.Fprint(w, "_我不知道你在说什么_")
	}

}

func main() {
	http.HandleFunc("/", handerFunc)
	http.ListenAndServe(":3000", nil)
}
