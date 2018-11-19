package main

import (
	"net/http"
	"fmt"
	"strings"
	"os"
)

func customHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Sending 'hello'")
	fmt.Fprintf(w,"Hello")
}

func clientHandler(w http.ResponseWriter, r *http.Request) {

	filePath :=  "." + r.URL.Path

	if !strings.HasPrefix(filePath, "./client/") {
		fmt.Println("ERROR: Invalid client path", filePath)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if strings.HasSuffix(filePath, "/") { filePath = filePath + "index.html" }

	info, err := os.Stat(filePath)
	if os.IsNotExist(err) || info.IsDir() {
		fmt.Println("ERROR: File not found", filePath)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	fmt.Println("Serving file", filePath)
	http.ServeFile(w, r, filePath)

}

func main() {

	http.HandleFunc("/client/", clientHandler)
	http.HandleFunc("/", customHandler)

	http.ListenAndServe(":8080", http.DefaultServeMux)
}