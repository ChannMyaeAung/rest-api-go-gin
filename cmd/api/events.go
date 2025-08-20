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
		return 
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
		return 
	}

	// Fetch event from database by ID
	event, err := app.models.Events.Get(id)

	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return 
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve event"})
		return 
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
		return
	}

	// Verify the event exists before attempting update
	existingEvent, err := app.models.Events.Get(id)
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve event"})
		return 
	}

	if existingEvent == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return 
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
		return 
	}

	if err := app.models.Events.Delete(id); err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete event"})
		return 
	}

	c.JSON(http.StatusNoContent, nil)
}

// addAttendeeToEvent handles POST /events/:id/attendees/:userId requests to register a user as an attendee for a specific event.
// It validates that both the event and user exist, checks for duplicate registrations,
// and creates a new attendee relationship in the database.
//
// URL Parameters:
//   - id: The unique identifier of the event
//   - userId: The unique identifier of the user to register as attendee
//
// Returns:
//   - 201 Created with attendee data on successful registration
//   - 400 Bad Request for invalid event or user ID format
//   - 404 Not Found if event or user doesn't exist
//   - 409 Conflict if user is already registered for the event
//   - 500 Internal Server Error on database failure
func (app *application) addAttendeeToEvent(c *gin.Context){
	eventId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return 
	}

	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Verify the target event exists before proceeding
	event, err := app.models.Events.Get(eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve event"})
		return 
	}
	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return 
	}

	// Verify the user to be added exists before proceeding
	userToAdd, err := app.models.Users.GetUserByID(userId)
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		return
	}
	if userToAdd == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Check if the user is already registered as an attendee for this event
    // This prevents duplicate registrations
	existingAttendee, err := app.models.Attendees.GetByEventAndAttendee(event.Id, userToAdd.Id)
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve attendee"})
		return
	}
	if existingAttendee != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User is already an attendee"})
		return
	}

	// Create new attendee relationship record
	attendee := database.Attendee{
		EventId: event.Id,
		UserId: userToAdd.Id,
	}

	// Insert the new attendee relationship into the database
	_, err = app.models.Attendees.Insert(&attendee)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add attendee"})
		return
	}

	c.JSON(http.StatusCreated, attendee)
}

// getAttendeesForEvent handles GET /events/:id/attendees requests to retrieve all users registered for a specific event.
// Returns a list of user objects representing the attendees for the given event.
//
// URL Parameters:
//   - id: The unique identifier of the event
//
// Returns:
//   - 200 OK with array of attendee user objects
//   - 400 Bad Request for invalid event ID format
//   - 500 Internal Server Error on database failure
func (app *application) getAttendeesForEvent(c *gin.Context){
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	// Fetch all users who are registered as attendees for this event
	users, err := app.models.Attendees.GetAttendeesByEvent(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve attendees"})
		return
	}

	// Return the list of attendees
	c.JSON(http.StatusOK, users)
}

// deleteAttendeeFromEvent handles DELETE /events/:id/attendees/:userId requests to remove a user's registration from an event.
// This allows users to unregister from events or administrators to remove attendees.
//
// URL Parameters:
//   - id: The unique identifier of the event
//   - userId: The unique identifier of the user to remove from attendees
//
// Returns:
//   - 204 No Content on successful removal
//   - 400 Bad Request for invalid event or user ID format
//   - 500 Internal Server Error on database failure
//
// Note: This implementation doesn't verify if the attendee relationship exists before deletion.
// The operation succeeds silently even if the user wasn't registered for the event.
func (app *application) deleteAttendeeFromEvent(c *gin.Context){
	eventId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return 
	}

	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Remove the attendee relationship from the database
	if err := app.models.Attendees.DeleteByEventAndUser(userId, eventId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete attendee"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// getEventsByAttendee handles GET /events/:id/attendees/:userId requests to retrieve all events that a specific user is registered for.
// This allows users to view their event registrations or administrators to see a user's event participation.
//
// URL Parameters:
//   - userId: The unique identifier of the user/attendee
//
// Returns:
//   - 200 OK with array of event objects the user is registered for
//   - 400 Bad Request for invalid user ID format
//   - 500 Internal Server Error on database failure
func (app *application) getEventsByAttendee(c *gin.Context){
	id, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid attendee ID"})
		return
	}

	// Fetch all events that this user is registered to attend
	events, err := app.models.Attendees.GetEventsByAttendee(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve events"})
		return
	}

	// Return the list of events the user is attending
	c.JSON(http.StatusOK, events)
}