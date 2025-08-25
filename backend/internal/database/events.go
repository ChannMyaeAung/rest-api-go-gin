package database

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"
)

type EventModel struct {
    DB *sql.DB 
}

// NullableTime handles potentially invalid datetime strings from database
type NullableTime struct {
    time.Time
    Valid bool
}

// Scan implements the Scanner interface for database/sql
func (nt *NullableTime) Scan(value interface{}) error {
    if value == nil {
        nt.Valid = false
        return nil
    }

    switch v := value.(type) {
    case time.Time:
        nt.Time = v
        nt.Valid = true
        return nil
    case string:
        // Try to parse various date formats
        formats := []string{
            time.RFC3339,                 // "2006-01-02T15:04:05Z07:00"
            "2006-01-02T15:04:05Z",      // "2025-05-20T00:00:00Z"
            "2006-01-02 15:04:05",       // "2025-01-02 15:04:05"
            "2006-01-02",                // "2025-01-02"
            "02-01-2006",                // "01-09-2025"
        }

        for _, format := range formats {
            if t, err := time.Parse(format, v); err == nil {
                nt.Time = t
                nt.Valid = true
                return nil
            }
        }

        // If all parsing fails, set a default time
        nt.Time = time.Now()
        nt.Valid = true
        return nil
    }

    return fmt.Errorf("cannot scan %T into NullableTime", value)
}

// Value implements the driver Valuer interface
func (nt NullableTime) Value() (driver.Value, error) {
    if !nt.Valid {
        return nil, nil
    }
    return nt.Time, nil
}

// Event represents an event record in the database.
type Event struct{
    Id          int       `json:"id"`
    OwnerId     int       `json:"ownerId"`
    Name        string    `json:"name" binding:"required,min=3"`
    Description string    `json:"description" binding:"required,min=10"`
    Date        time.Time `json:"date" binding:"required"`
    Location    string    `json:"location" binding:"required,min=3"`
}

// Insert creates a new event record in the database.
func (m *EventModel) Insert(event *Event) error{
    ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
    defer cancel()

    query := "INSERT INTO events (owner_id, name, description, date, location) VALUES ($1, $2, $3, $4, $5) RETURNING id"

    return m.DB.QueryRowContext(ctx, query, event.OwnerId, event.Name, event.Description, event.Date, event.Location).Scan(&event.Id)
}

// GetAll retrieves all events from the database with better error handling
func (m *EventModel) GetAll() ([]*Event, error){
    ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
    defer cancel() 

    query := "SELECT id, owner_id, name, description, date, location FROM events"

    rows, err := m.DB.QueryContext(ctx, query)
    if err != nil{
        return nil, err 
    }
    defer rows.Close()

    events := []*Event{}

    for rows.Next(){
        var event Event
        var dateStr string

        // Scan date as string first to handle invalid formats
        err := rows.Scan(&event.Id, &event.OwnerId, &event.Name, &event.Description, &dateStr, &event.Location)
        if err != nil {
            return nil, err 
        }

        // Try to parse the date string
        parsedDate, err := parseFlexibleDate(dateStr)
        if err != nil {
            // If parsing fails, skip this event or use default date
            fmt.Printf("Warning: Invalid date format for event %d: %s\n", event.Id, dateStr)
            continue // Skip this event
            // Or use: event.Date = time.Now() // Use current time as default
        } else {
            event.Date = parsedDate
        }

        events = append(events, &event)
    }

    if err = rows.Err(); err != nil {
        return nil, err 
    }

    return events, nil
}

// parseFlexibleDate tries to parse various date formats
func parseFlexibleDate(dateStr string) (time.Time, error) {
    formats := []string{
        time.RFC3339,                 // "2006-01-02T15:04:05Z07:00"
        "2006-01-02T15:04:05Z",      // "2025-05-20T00:00:00Z"
        "2006-01-02 15:04:05",       // "2025-01-02 15:04:05"
        "2006-01-02",                // "2025-01-02"
        "02-01-2006",                // "01-09-2025"
    }

    for _, format := range formats {
        if t, err := time.Parse(format, dateStr); err == nil {
            return t, nil
        }
    }

    return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
}

// Get retrieves a single event by its ID with better date handling
func (m *EventModel) Get(id int) (*Event, error){
    ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
    defer cancel() 

    query := "SELECT id, owner_id, name, description, date, location FROM events WHERE id = $1"

    var event Event
    var dateStr string

    err := m.DB.QueryRowContext(ctx, query, id).Scan(&event.Id, &event.OwnerId, &event.Name, &event.Description, &dateStr, &event.Location)

    if err != nil{
        if err == sql.ErrNoRows{
            return nil, nil
        }
        return nil, err
    }

    // Parse the date
    parsedDate, err := parseFlexibleDate(dateStr)
    if err != nil {
        // Use current time as fallback
        event.Date = time.Now()
    } else {
        event.Date = parsedDate
    }

    return &event, nil 
}

// Update, Delete methods remain the same...
func (m *EventModel) Update(event *Event) error{
    ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
    defer cancel()

    query := "UPDATE events SET owner_id = $1, name = $2, description = $3, date = $4, location = $5 WHERE id = $6"

    _, err := m.DB.ExecContext(ctx, query, event.OwnerId, event.Name, event.Description, event.Date, event.Location, event.Id)

    if err != nil {
        return err 
    }

    return nil
}

func (m *EventModel) Delete(id int) error {
    ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
    defer cancel()

    query := "DELETE FROM events WHERE id = $1"

    _, err := m.DB.ExecContext(ctx, query, id)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil
        }
        return err
    }
    return nil
}