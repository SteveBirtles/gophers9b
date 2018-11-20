package main

import (
	"encoding/json"
	"fmt"
	"github.com/satori/go.uuid"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Message struct {
	Id       int    `json:"id"`
	Text     string `json:"text"`
	PostDate string `json:"postDate"`
	Author   string `json:"author"`
}

type User struct {
	Id           int    `json:"id"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	SessionToken string `json:"sessionToken"`
}

var (
	messages = make([]Message, 0)
	users    = make([]User, 0)
)

func clientHandler(w http.ResponseWriter, r *http.Request) {

	filePath := "." + r.URL.Path

	if !strings.HasPrefix(filePath, "./client/") {
		fmt.Println("ERROR: Invalid client path", filePath)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

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

func validateSessionToken(r *http.Request) string {

	var sessionToken string
	cookie, err := r.Cookie("sessionToken")
	if err == nil {
		sessionToken = cookie.Value
	}

	for i := range users {
		if users[i].SessionToken == sessionToken {
			return users[i].Username
		}
	}

	return ""

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

	username := validateSessionToken(r)
	if username == "" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	messageText := r.FormValue("messageText")

	id := 0
	for _, m := range messages {
		if m.Id > id {
			id = m.Id
		}
	}
	id++

	date := time.Now().Format("02-03-2006 15:04:05")

	fmt.Println("/message/new", id, messageText, username, date)

	messages = append(messages, Message{
		Id:       id,
		Text:     messageText,
		PostDate: date,
		Author:   username,
	})
	fmt.Fprint(w, "OK")

}

func deleteMessageHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	username := validateSessionToken(r)
	id, err := strconv.Atoi(r.FormValue("messageId"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Println("/message/delete", username, id)

	for i := 0; i < len(messages); {
		if messages[i].Id == id {
			if messages[i].Author != username {
				fmt.Println(username, messages[i].Author)
				fmt.Fprint(w, "Error: You didn't post that originally")
				return
			}
			messages = append(messages[:i], messages[i+1:]...)
			fmt.Fprint(w, "OK")
			return
		} else {
			i++
		}
	}

	fmt.Fprint(w, "Error: Can't find post with id ", id)

}

func editMessageHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	username := validateSessionToken(r)
	id, err := strconv.Atoi(r.FormValue("messageId"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	newText := r.FormValue("messageText")

	for i := range messages {
		if messages[i].Id == id {
			if messages[i].Author != username {
				fmt.Fprint(w, "Error: You didn't post that originally")
				return
			}
			messages[i].Text = newText
		}
	}

	fmt.Println("/message/edit", username, id, newText)
	fmt.Fprint(w, "OK")
}

func newUserHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	password1 := r.FormValue("password1")
	password2 := r.FormValue("password2")

	if password1 != password2 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Error: Passwords don't match")
		return
	}

	id := 0
	for _, u := range users {
		if u.Id > id {
			id = u.Id
		}
	}
	id++

	token := uuid.Must(uuid.NewV4()).String()
	users = append(users, User{id, username, password1, token})

	fmt.Println("/user/new", username, password1)
	fmt.Fprint(w, token)
}

func loginUserHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	fmt.Println("/user/login", username, password)

	for i := range users {
		if users[i].Username == username {
			if users[i].Password != password {
				fmt.Fprint(w, "Error: Incorrect password")
				return
			}
			token := uuid.Must(uuid.NewV4()).String()
			users[i].SessionToken = token
			fmt.Fprint(w, token)
			return
		}
	}
	fmt.Fprint(w, "Error: Can't find user account.")

}

func getUserHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	username := validateSessionToken(r)
	fmt.Println("/user/get", username)
	fmt.Fprint(w, username)

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
