package database

import (
	"context"
	"database/sql"
	"time"
)

type EventModel struct {
	DB *sql.DB 
}

// Event represents an event record in the database.
// Binding is the process of automatically parsing the incoming request data and mapping it into a Go struct.
type Event struct{
	Id int `json:"id"`
	OwnerId int `json:"ownerId"`
	Name string `json:"name" binding:"required,min=3"`
	Description string `json:"description" binding:"required,min=10"`
	Date string `json:"date" binding:"required"`
	Location string `json:"location" binding:"required,min=3"`
}

// Insert creates a new event record in the database.
// It automatically sets the event ID from the database's auto-increment value.
// 
// The function uses a prepared statement to prevent SQL injection and includes
// a 3-second timeout to prevent the query from hanging indefinitely.
//
// Parameters:
//   - event: Pointer to Event struct containing the event data to insert
//
// Returns:
//   - error: nil on success, database error on failure
func (m *EventModel) Insert(event *Event) error{
	// Prevents the query from hanging forever.
	// If the DB doesn't respond within 3 seconds, it cancels automatically.
	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	query := "INSERT INTO events (owner_id, name, description, date, location) VALUES ($1, $2, $3, $4, $5) RETURNING id"

	return m.DB.QueryRowContext(ctx, query, event.OwnerId, event.Name, event.Description, event.Date, event.Location).Scan(&event.Id)
}


// GetAll retrieves all events from the database.
// Returns a slice of pointers to Event structs containing all event records.
//
// Returns:
//   - []*Event: Slice of event pointers, empty slice if no events found
//   - error: nil on success, database error on failure
func (m *EventModel) GetAll() ([]*Event, error){
	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel() 

	query := "SELECT * FROM events"

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil{
		return nil, err 
	}

	defer rows.Close()

	// Initialize slice to hold event pointers
	events := []*Event{}

	// Iterate through all rows in the result set
	for rows.Next(){
		var event Event 

		// Scan each row's columns into the event struct fields
        // Column order must match the SELECT statement and database schema
		err := rows.Scan(&event.Id, &event.OwnerId, &event.Name, &event.Description, &event.Date, &event.Location)

		if err != nil {
			return nil, err 
		}

		// Append pointer to event to the slice
		events = append(events, &event)
	}

	// Check for any errors that occurred during row iteration
	if err = rows.Err(); err != nil {
		return nil, err 
	}

	return events, nil
}

// Get retrieves a single event by its ID.
// Returns nil if no event is found with the given ID.
//
// Parameters:
//   - id: The unique identifier of the event to retrieve
//
// Returns:
//   - *Event: Pointer to Event struct if found, nil if not found
//   - error: nil on success or when no rows found, database error on failure
func (m *EventModel) Get(id int) (*Event, error){
	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel() 

	query := "SELECT * FROM events WHERE id = $1"

	var event Event

	// QueryRowContext expects exactly one row; use for single record queries
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&event.Id, &event.OwnerId, &event.Name, &event.Description, &event.Date, &event.Location)

	if err != nil{
		if err == sql.ErrNoRows{
			return nil, nil // No event found with the given ID
		}
		return nil, err // Return actual database error
	}

	return &event, nil 

}

// Update modifies an existing event record in the database.
// Updates all fields of the event based on the provided Event struct.
// The event ID is used to identify which record to update.
//
// Parameters:
//   - event: Pointer to Event struct containing updated data and the ID of the record to modify
//
// Returns:
//   - error: nil on success, database error on failure
func (m *EventModel) Update(event *Event) error{
	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	// Update all event fields where ID matches
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
			return nil // No event found with the given ID
		}
		return err
	}
	return nil
}