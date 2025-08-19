package main

import (
	"net/http"
	"rest-api-in-gin/internal/database"
	"strconv"

	"github.com/gin-gonic/gin"
)

// createEvent handles POST /events requests to create a new event.
// It binds JSON request data to an Event struct, validates the input,
// and inserts the event into the database via the Events model.
//
// Request body should contain JSON with event fields (name, description, location, dateTime, userId).
// Returns 201 Created with the created event on success, or appropriate error status codes.
func (app *application) createEvent(c *gin.Context){
	var event database.Event

	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	// Insert the event into the database
	err := app.models.Events.Insert(&event)

	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create event"})
		return 
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Event created successfully", "event": event})
}

// getAllEvents handles GET /events requests to retrieve all events.
// Returns a JSON array of all events in the database.
//
// Returns 200 OK with array of events on success, or 500 Internal Server Error on database failure.
func (app *application) getAllEvents(c *gin.Context){
	// Fetch all events from database
	events, err := app.models.Events.GetAll()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve events"})
		return 
	}	

	// Return events as JSON array
	c.JSON(http.StatusOK, events)
}

// getEventByID handles GET /events/:id requests to retrieve a specific event by ID.
// The event ID is extracted from the URL path parameter.
//
// Returns 200 OK with event data on success, 400 Bad Request for invalid ID format,
// 404 Not Found if event doesn't exist, or 500 Internal Server Error on database failure.
func (app *application) getEventByID(c *gin.Context){
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
	}

	// Fetch event from database by ID
	event, err := app.models.Events.Get(id)

	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve event"})
	}

	c.JSON(http.StatusOK, event)
}

// updateEvent handles PUT /events/:id requests to update an existing event.
// It first verifies the event exists, then binds the JSON request body to update the event.
// The ID from the URL path takes precedence over any ID in the request body.
//
// Returns 200 OK with updated event on success, 400 Bad Request for invalid input,
// 404 Not Found if event doesn't exist, or 500 Internal Server Error on database failure.
func (app *application) updateEvent(c *gin.Context){
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
	}

	// Verify the event exists before attempting update
	existingEvent, err := app.models.Events.Get(id)
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve event"})
		return 
	}

	if existingEvent == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
	}

	// Create new Event struct to hold update data
	updatedEvent := &database.Event{}

	if err := c.ShouldBindJSON(updatedEvent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return 
	}

	// Ensure the ID from URL path is used (not from request body)
	updatedEvent.Id = id

	// Update the event in the database
	if err := app.models.Events.Update(updatedEvent); err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update event"})
		return 
	}

	c.JSON(http.StatusOK, updatedEvent)
}

// deleteEvent handles DELETE /events/:id requests to remove an event from the database.
// The event ID is extracted from the URL path parameter.
//
// Returns 204 No Content on successful deletion, 400 Bad Request for invalid ID format,
// or 500 Internal Server Error on database failure.
// Note: This implementation doesn't check if the event exists before deletion.
func (app *application) deleteEvent(c *gin.Context){
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
	}

	if err := app.models.Events.Delete(id); err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete event"})
	}

	c.JSON(http.StatusNoContent, nil)
}