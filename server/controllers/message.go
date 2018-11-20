package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/stevebirtles/gophers9a/server/models"
)

func ListMessages(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	fmt.Println("/message/list")

	messages := models.SelectAllMessages()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)

}

func NewMessage(w http.ResponseWriter, r *http.Request) {

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

	messages := models.SelectAllMessages()

	id := 0
	for _, m := range messages {
		if m.Id > id {
			id = m.Id
		}
	}
	id++

	date := time.Now().Format("02-03-2006 15:04:05")

	fmt.Println("/message/new", id, messageText, username, date)

	models.InsertMessage(models.Message{
		Id:       id,
		Text:     messageText,
		PostDate: date,
		Author:   username,
	})

	fmt.Fprint(w, "OK")

}

func EditMessage(w http.ResponseWriter, r *http.Request) {

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

	messages := models.SelectAllMessages()

	for i := range messages {
		if messages[i].Id == id {
			if messages[i].Author != username {
				fmt.Fprint(w, "Error: You didn't post that originally")
				return
			}
			messages[i].Text = newText
			models.UpdateMessage(messages[i])
		}
	}

	fmt.Println("/message/edit", username, id, newText)
	fmt.Fprint(w, "OK")
}

func DeleteMessage(w http.ResponseWriter, r *http.Request) {

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

	messages := models.SelectAllMessages()

	for i := range messages {
		if messages[i].Id == id {
			if messages[i].Author != username {
				fmt.Println(username, messages[i].Author)
				fmt.Fprint(w, "Error: You didn't post that originally")
				return
			}
			models.DeleteMessage(id)
			fmt.Fprint(w, "OK")
			return
		}
	}

	fmt.Fprint(w, "Error: Can't find post with id ", id)

}
