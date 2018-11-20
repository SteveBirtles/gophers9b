package models

import "fmt"

type Message struct {
	Id       int    `json:"id"`
	Text     string `json:"text"`
	PostDate string `json:"postDate"`
	Author   string `json:"author"`
}

func SelectAllMessages() []Message {

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

func InsertMessage(message Message) {

	_, err := database.Exec(fmt.Sprintf("insert into Messages (Id, Text, PostDate, Author) values (%d, '%s', '%s', '%s')",
		message.Id, message.Text, message.PostDate, message.Author))

	if err != nil {
		fmt.Println("Database insert error:", err)
	}

}

func UpdateMessage(message Message) {

	_, err := database.Exec(fmt.Sprintf("update Messages set Text = '%s', PostDate = '%s', Author = '%s' where Id = %d",
		message.Text, message.PostDate, message.Author, message.Id))

	if err != nil {
		fmt.Println("Database update error:", err)
	}

}

func DeleteMessage(id int) {

	_, err := database.Exec(fmt.Sprintf("delete from Messages where Id = %d", id))

	if err != nil {
		fmt.Println("Database delete error:", err)
	}
}
