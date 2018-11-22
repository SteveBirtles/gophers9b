package controllers

import (
	"fmt"
	"github.com/satori/go.uuid"
	"net/http"

	"github.com/stevebirtles/gophers9b/server/models"
)

func validateSessionToken(r *http.Request) string {

	var sessionToken string
	cookie, err := r.Cookie("sessionToken")
	if err == nil {
		sessionToken = cookie.Value
	}

	users := models.SelectAllUsers()

	for i := range users {
		if users[i].SessionToken == sessionToken {
			return users[i].Username
		}
	}

	return ""

}

func GetUser(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	username := validateSessionToken(r)
	fmt.Println("/user/get", username)
	fmt.Fprint(w, username)

}

func NewUser(w http.ResponseWriter, r *http.Request) {

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

	users := models.SelectAllUsers()

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

	models.InsertUser(models.User{Id: id, Username: username, Password: password1, SessionToken: token})

	fmt.Println("/user/new", username, password1)
	fmt.Fprint(w, token)
}

func AttemptLogin(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	fmt.Println("/user/login", username, password)

	users := models.SelectAllUsers()

	for i := range users {
		if users[i].Username == username {
			if users[i].Password != password {
				fmt.Fprint(w, "Error: Incorrect password")
				return
			}
			token := uuid.Must(uuid.NewV4()).String()
			users[i].SessionToken = token
			models.UpdateUser(users[i])
			fmt.Fprint(w, token)
			return
		}
	}
	fmt.Fprint(w, "Error: Can't find user account.")

}
