package controllers

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

func StaticContent(w http.ResponseWriter, r *http.Request) {

	filePath := "./resources" + r.URL.Path

	if strings.HasSuffix(filePath, "/") {
		filePath = filePath + "index.html"
	}

	info, err := os.Stat(filePath)
	if os.IsNotExist(err) || info.IsDir() {
		if !strings.HasSuffix(filePath, ".map") {
			fmt.Println("ERROR: File not found", filePath)
		}
		w.WriteHeader(http.StatusNotFound)
		return
	}

	fmt.Println("Serving file", filePath)
	http.ServeFile(w, r, filePath)

}
