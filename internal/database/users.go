package database

import "database/sql"

type UserModel struct {
	DB *sql.DB 
}

type User struct{
	Id int `json:"id"`
	Email string `json:"email"`
	Name string `json:"name"`

	// tells JSON library to always ignore this field when converting the struct back into JSON
	// to prevent from accidentally sending a user's password hash back to the client.
	Password string `json:"-"`
}