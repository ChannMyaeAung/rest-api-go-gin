package database

import (
	"context"
	"database/sql"
<<<<<<< HEAD
=======
	"errors"
>>>>>>> b2b83c2 (Added add-attendee page, menus for profile and settings)
	"time"
)

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

<<<<<<< HEAD
func (m *UserModel) Insert(user *User) error{
	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
=======
var defaultTimeout = 3 * time.Second

func (m *UserModel) Insert(user *User) error{
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
>>>>>>> b2b83c2 (Added add-attendee page, menus for profile and settings)
	defer cancel()

	query := "INSERT INTO users (email, password, name) VALUES ($1, $2, $3) RETURNING id"

	return m.DB.QueryRowContext(ctx, query, user.Email, user.Password, user.Name).Scan(&user.Id)
}

func (m *UserModel) getUser(query string, args ...interface{}) (*User, error){
<<<<<<< HEAD
	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
=======
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
>>>>>>> b2b83c2 (Added add-attendee page, menus for profile and settings)
	defer cancel()

	var user User 
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.Id, &user.Email, &user.Name, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (m *UserModel) GetUserByID(id int) (*User, error){
	query := "SELECT * FROM users WHERE id = $1"
	return m.getUser(query, id)
}

func (m *UserModel) GetByEmail(email string) (*User, error){
	query := "SELECT * FROM users WHERE email = $1"
	return m.getUser(query, email)
<<<<<<< HEAD
=======
}

func (m *UserModel) Delete(ctx context.Context, id int) error{
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	tx, err := m.DB.BeginTx(ctx, nil)
	if err != nil {
		return err 
	}

	if _, err := tx.ExecContext(ctx, `DELETE FROM attendees WHERE event_id IN (SELECT id FROM events WHERE owner_id = ?)`, id); err != nil {
		tx.Rollback()
		return err 
	}

	if _, err := tx.ExecContext(ctx, `DELETE FROM events WHERE owner_id = ?`, id); err != nil {
		tx.Rollback()
		return err
	}

	result, err := tx.ExecContext(ctx, `DELETE FROM users WHERE id = ?`, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err 
	}

	if rowsAffected == 0 {
		tx.Rollback()
		return errors.New("not found")
	}
	return tx.Commit()
>>>>>>> b2b83c2 (Added add-attendee page, menus for profile and settings)
}