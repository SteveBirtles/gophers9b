package models

import "fmt"

type User struct {
	Id           int    `json:"id"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	SessionToken string `json:"sessionToken"`
}

func SelectAllUsers() []User {

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

func InsertUser(user User) {

	_, err := database.Exec(fmt.Sprintf("insert into Users (Id, Username, Password, SessionToken) values (%d, '%s', '%s', '%s')",
		user.Id, user.Username, user.Password, user.SessionToken))

	if err != nil {
		fmt.Println("Database insert error:", err)
	}

}

func UpdateUser(user User) {

	_, err := database.Exec(fmt.Sprintf("update Users set Username = '%s', Password = '%s', SessionToken = '%s' where Id = %d",
		user.Username, user.Password, user.SessionToken, user.Id))

	if err != nil {
		fmt.Println("Database update error:", err)
	}

}
