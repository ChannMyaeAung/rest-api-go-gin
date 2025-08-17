package database

import (
	"database/sql"
	"time"
)

type EventModel struct {
	DB *sql.DB 
}

// Binding is the process of automatically parsing the incoming request data and mapping it into a Go struct.
type Event struct{
	Id int `json:"id"`
	OwnerId int `json:"ownerId" binding:"required"`
	Name string `json:"name" binding:"required,min=3"`
	Description string `json:"description" binding:"required,min=10"`
	Date time.Time `json:"date" binding:"required, datetime=2006-01-02"`
	Location string `json:"location" binding:"required, min=3"`
}