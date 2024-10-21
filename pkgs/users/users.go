package users

import (
	"database/sql"
	"fmt"
)

/*
CREATE TABLE users (

	id INTEGER PRIMARY KEY AUTOINCREMENT,
	username TEXT NOT NULL UNIQUE,
	email TEXT NOT NULL UNIQUE,
	password TEXT NOT NULL

);
*/
type User struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func New(username, email, password string) *User {
	return &User{
		Id:       -1,
		Username: username,
		Email:    email,
		Password: password,
	}
}

func (u *User) AddToDb(db *sql.DB) error {
	//Create statement
	prompt := "INSERT INTO users (username, email, password) VALUES (?, ?, ?)"
	statement, err := db.Prepare(prompt)
	if err != nil {
		return err
	}
	defer statement.Close()

	//Execute statement
	result, err := statement.Exec(u.Username, u.Email, u.Password)
	if err != nil {
		return err
	}

	//Assign the user their id
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	u.Id = id

	return nil
}

func (u *User) Update(db *sql.DB) error {
	prompt := "UPDATE users SET username = ?, email = ?, password = ? WHERE id = ?"
	statement, err := db.Prepare(prompt)
	if err != nil {
		fmt.Println("Error at statement")
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(u.Username, u.Email, u.Password, u.Id)

	if err != nil {
		fmt.Println("Error at execution")
		return err
	}

	return nil
}

func FindById(id int64, db *sql.DB) (*User, error) {
	prompt := "SELECT * FROM users WHERE id=?"
	row := db.QueryRow(prompt, id)
	return rowToUser(row)
}

func FindByEmail(email string, db *sql.DB) (*User, error) {
	prompt := "SELECT * FROM users WHERE email=?"
	row := db.QueryRow(prompt, email)
	return rowToUser(row)

}

func FindByUsername(username string, db *sql.DB) (*User, error) {
	prompt := "SELECT * FROM users WHERE username=?"
	row := db.QueryRow(prompt, username)
	return rowToUser(row)
}

func rowToUser(row *sql.Row) (*User, error) {
	user := User{
		Id: -1,
	}
	err := row.Scan(&user.Id, &user.Username, &user.Email, &user.Password)
	return &user, err
}
