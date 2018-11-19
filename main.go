package main

import (
	"net/http"
	"fmt"
	"strings"
	"os"
)

func customHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("200: Sending 'hello'")
	fmt.Fprintf(w,"Hello")
}

func clientHandler(w http.ResponseWriter, r *http.Request) {

	filePath :=  "." + r.URL.Path

	if !strings.HasPrefix(filePath, "./client/") {
		fmt.Println("500: Internal server error", filePath)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if strings.HasSuffix(filePath, "/") { filePath = filePath + "index.html" }

	info, err := os.Stat(filePath)
	if os.IsNotExist(err) || info.IsDir() {
		fmt.Println("404: File not found", filePath)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	fmt.Println("200: Serving file", filePath)
	http.ServeFile(w, r, filePath)

}

func main() {

	http.HandleFunc("/client/", clientHandler)

	http.HandleFunc("/", customHandler)


	http.ListenAndServe(":8080", nil)
}