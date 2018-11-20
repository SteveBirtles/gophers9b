package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/satori/go.uuid"
	"net/http"
	"os"
	"os/signal"
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

var database *sql.DB

func selectAllUsers() []User {

	users := make([]User, 0)

	rows, err := database.Query("select Id, Username, Password, SessionToken from Users")
	if err != nil {
		fmt.Println("Database select all error:", err)
		return users
	}
	defer rows.Close()

	for rows.Next() {
		var u User
		err = rows.Scan(&u.Id, &u.Username, &u.Password, &u.SessionToken)
		if err != nil {
			fmt.Println("Database select all error:", err)
			break
		}
		users = append(users, u)
	}

	return users

}

func insertUser(user User) {

	_, err := database.Exec(fmt.Sprintf("insert into Users (Id, Username, Password, SessionToken) values (%d, '%s', '%s', '%s')",
		user.Id, user.Username, user.Password, user.SessionToken))

	if err != nil {
		fmt.Println("Database insert error:", err)
	}

}

func updateUser(user User) {

	_, err := database.Exec(fmt.Sprintf("update Users set Username = '%s', Password = '%s', SessionToken = '%s' where Id = %d",
		user.Username, user.Password, user.SessionToken, user.Id))

	if err != nil {
		fmt.Println("Database update error:", err)
	}

}

func selectAllMessages() []Message {

	messages := make([]Message, 0)

	rows, err := database.Query("select Id, Text, PostDate, Author from Messages")
	if err != nil {
		fmt.Println("Database select all error:", err)
		return messages
	}
	defer rows.Close()

	for rows.Next() {
		var m Message
		err = rows.Scan(&m.Id, &m.Text, &m.PostDate, &m.Author)
		if err != nil {
			fmt.Println("Database select all error:", err)
			break
		}
		messages = append(messages, m)
	}

	return messages
}

func insertMessage(message Message) {

	_, err := database.Exec(fmt.Sprintf("insert into Messages (Id, Text, PostDate, Author) values (%d, '%s', '%s', '%s')",
		message.Id, message.Text, message.PostDate, message.Author))

	if err != nil {
		fmt.Println("Database insert error:", err)
	}

}

func updateMessage(message Message) {

	_, err := database.Exec(fmt.Sprintf("update Messages set Text = '%s', PostDate = '%s', Author = '%s' where Id = %d",
		message.Text, message.PostDate, message.Author, message.Id))

	if err != nil {
		fmt.Println("Database update error:", err)
	}

}

func deleteMessage(id int) {

	_, err := database.Exec(fmt.Sprintf("delete from Messages where Id = %d", id))

	if err != nil {
		fmt.Println("Database delete error:", err)
	}
}

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

	users := selectAllUsers()

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

	messages := selectAllMessages()

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

	messages := selectAllMessages()

	id := 0
	for _, m := range messages {
		if m.Id > id {
			id = m.Id
		}
	}
	id++

	date := time.Now().Format("02-03-2006 15:04:05")

	fmt.Println("/message/new", id, messageText, username, date)

	insertMessage(Message{
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

	messages := selectAllMessages()

	for i := range messages {
		if messages[i].Id == id {
			if messages[i].Author != username {
				fmt.Println(username, messages[i].Author)
				fmt.Fprint(w, "Error: You didn't post that originally")
				return
			}
			deleteMessage(id)
			fmt.Fprint(w, "OK")
			return
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

	messages := selectAllMessages()

	for i := range messages {
		if messages[i].Id == id {
			if messages[i].Author != username {
				fmt.Fprint(w, "Error: You didn't post that originally")
				return
			}
			messages[i].Text = newText
			updateMessage(messages[i])
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
		fmt.Fprint(w, "Error: Passwords don't match")
		return
	}

	users := selectAllUsers()

	id := 0
	for _, u := range users {
		if u.Username == username {
			fmt.Fprint(w, "Error: User already exists")
			return
		}
		if u.Id > id {
			id = u.Id
		}
	}
	id++

	token := uuid.Must(uuid.NewV4()).String()

	insertUser(User{id, username, password1, token})

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

	users := selectAllUsers()

	for i := range users {
		if users[i].Username == username {
			if users[i].Password != password {
				fmt.Fprint(w, "Error: Incorrect password")
				return
			}
			token := uuid.Must(uuid.NewV4()).String()
			users[i].SessionToken = token
			updateUser(users[i])
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

func init() {

	var err error
	fmt.Println("Connecting to database...")
	database, err = sql.Open("sqlite3", "MessageBoard.db")
	if err != nil {
		panic(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go awaitShutdown(c)

}

func awaitShutdown(c chan os.Signal) {
	for range c {
		fmt.Println("Disconnecting from to database...")
		database.Close()
		os.Exit(0)
	}
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
