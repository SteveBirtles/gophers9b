package main

import (
	"net/http"
	"fmt"
	"strings"
	"os"
	"encoding/json"
	"github.com/satori/go.uuid"
)

type Message struct {
	Id int `json:"id"`
	Text string `json:"text"`
	PostDate string `json:"postDate"`
	Author string `json:"author"`
}

type User struct {
	Id int `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	SessionToken string `json:"sessionToken"`
}

var (
	messages = []Message{{1, "Testing", "Today", "Steve"}}
	users = []User{{1, "Steve", "beans", ""}}
)

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

func listMessagesHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	fmt.Println("/message/list")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)

}

func newMessageHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	fmt.Println("/message/new")
}

func deleteMessageHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	fmt.Println("/message/delete")
}

func editMessageHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	fmt.Println("/message/edit")
}

func newUserHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	fmt.Println("/user/new")
}

func loginUserHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	fmt.Println("/user/login")

	username := r.FormValue("username")
	password := r.FormValue("password")

	for i := range users {
		if users[i].Username == username {
			if  users[i].Password != password {
				fmt.Fprintln(w, "Error: Incorrect password")
				return
			}
			token := uuid.Must(uuid.NewV4()).String()
			users[i].SessionToken = token
			fmt.Fprintln(w, token)
			return
		}
	}
	fmt.Fprintln(w, "Error: Can't find user account.")

}

func getUserHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	fmt.Println("/user/get")

	var sessionToken string
	cookie, err := r.Cookie("sessionToken")
	if err != nil {
		sessionToken = cookie.Value
	}

	for i := range users {
		if users[i].SessionToken == sessionToken {
			fmt.Fprintln(w, users[i].Username)
			return
		}
	}

	fmt.Fprintln(w, "")
}

func main() {

	http.HandleFunc("/client/", clientHandler)

	http.HandleFunc("/message/list", listMessagesHandler)
	http.HandleFunc("/message/new", newMessageHandler)
	http.HandleFunc("/message/delete", deleteMessageHandler)
	http.HandleFunc("/message/edit", editMessageHandler)

	http.HandleFunc("/user/new", newUserHandler)
	http.HandleFunc("/user/login", loginUserHandler)
	http.HandleFunc("/user/get", getUserHandler)

	http.ListenAndServe(":8081", http.DefaultServeMux)
}