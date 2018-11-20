package main

import (
	"net/http"

	"github.com/stevebirtles/gophers9a/server/controllers"
	_ "github.com/stevebirtles/gophers9a/server/models"
)

func main() {

	http.HandleFunc("/client/", controllers.StaticContent)

	http.HandleFunc("/message/list", controllers.ListMessages)
	http.HandleFunc("/message/new", controllers.NewMessage)
	http.HandleFunc("/message/delete", controllers.DeleteMessage)
	http.HandleFunc("/message/edit", controllers.EditMessage)

	http.HandleFunc("/user/new", controllers.NewUser)
	http.HandleFunc("/user/login", controllers.AttemptLogin)
	http.HandleFunc("/user/get", controllers.GetUser)

	http.ListenAndServe(":8081", http.DefaultServeMux)

}
