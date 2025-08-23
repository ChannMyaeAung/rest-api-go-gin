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
// @Summary Create a new event
// @Description Create a new event with the provided details. Requires authentication.
// @Tags Events
// @Accept json
// @Produce json
// @Param event body database.Event true "Event data"
// @Success 201 {object} gin.H "Event created successfully"
// @Failure 400 {object} gin.H "Invalid request body"
// @Failure 401 {object} gin.H "Unauthorized"
// @Failure 500 {object} gin.H "Internal server error"
// @Security BearerAuth
// @Router /api/v1/events [post]
func (app *application) createEvent(c *gin.Context){
	var event database.Event

	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return 
	}

	user := app.getUserFromContext(c)
	event.OwnerId = user.Id // Set the event owner to the authenticated user

	// Insert the event into the database
	err := app.models.Events.Insert(&event)

	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create event"})
		return 
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Event created successfully", "event": event})
}

// getEvents return all events
//
// @Summary Returns all events
// @Description Returns all events
// @Tags Events
// @Accept json 
// @Product json 
// @Success 200 {object} []database.Event 
// @Router /api/v1/events [get]
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
// @Summary Get event by ID
// @Description Retrieve a specific event by its ID
// @Tags Events
// @Accept json
// @Produce json
// @Param id path int true "Event ID"
// @Success 200 {object} database.Event "Event details"
// @Failure 400 {object} gin.H "Invalid event ID"
// @Failure 404 {object} gin.H "Event not found"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /api/v1/events/{id} [get]
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
// It first verifies the event exists, then checks if the authenticated user owns the event,
// and finally binds the JSON request body to update the event.
// Only the event owner can update their events.
//
// @Summary Update an event
// @Description Update an existing event. Only the event owner can update their event.
// @Tags Events
// @Accept json
// @Produce json
// @Param id path int true "Event ID"
// @Param event body database.Event true "Updated event data"
// @Success 200 {object} database.Event "Updated event"
// @Failure 400 {object} gin.H "Invalid request"
// @Failure 401 {object} gin.H "Unauthorized"
// @Failure 403 {object} gin.H "Forbidden - not event owner"
// @Failure 404 {object} gin.H "Event not found"
// @Failure 500 {object} gin.H "Internal server error"
// @Security BearerAuth
// @Router /api/v1/events/{id} [put]
func (app *application) updateEvent(c *gin.Context){
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	user := app.getUserFromContext(c)

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

	if existingEvent.OwnerId != user.Id{
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to update this event"})
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
// Only the event owner can delete their events. The event ID is extracted from the URL path parameter.
//
// @Summary Delete an event
// @Description Delete an event from the database. Only the event owner can delete their event.
// @Tags Events
// @Accept json
// @Produce json
// @Param id path int true "Event ID"
// @Success 204 "Event deleted successfully"
// @Failure 400 {object} gin.H "Invalid event ID"
// @Failure 401 {object} gin.H "Unauthorized"
// @Failure 403 {object} gin.H "Forbidden - not event owner"
// @Failure 404 {object} gin.H "Event not found"
// @Failure 500 {object} gin.H "Internal server error"
// @Security BearerAuth
// @Router /api/v1/events/{id} [delete]
func (app *application) deleteEvent(c *gin.Context){
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return 
	}

	user := app.getUserFromContext(c)

	existingEvent, err := app.models.Events.Get(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve event"})
		return 
	}

	if existingEvent == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return 
	}

	if existingEvent.OwnerId != user.Id{
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to delete this event"})
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
// Only the event owner can add attendees to their events.
//
// @Summary Add attendee to event
// @Description Register a user as an attendee for a specific event. Only the event owner can add attendees.
// @Tags Events
// @Accept json
// @Produce json
// @Param id path int true "Event ID"
// @Param userId path int true "User ID to add as attendee"
// @Success 201 {object} database.Attendee "Attendee added successfully"
// @Failure 400 {object} gin.H "Invalid event or user ID"
// @Failure 401 {object} gin.H "Unauthorized"
// @Failure 403 {object} gin.H "Forbidden - not event owner"
// @Failure 404 {object} gin.H "Event or user not found"
// @Failure 409 {object} gin.H "User already an attendee"
// @Failure 500 {object} gin.H "Internal server error"
// @Security BearerAuth
// @Router /api/v1/events/{id}/attendees/{userId} [post]
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

	user := app.getUserFromContext(c)

	if event.OwnerId != user.Id{
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to add attendees to this event"})
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
// @Summary Get attendees for event
// @Description Retrieve all users registered as attendees for a specific event
// @Tags Events
// @Accept json
// @Produce json
// @Param id path int true "Event ID"
// @Success 200 {array} database.User "List of attendees"
// @Failure 400 {object} gin.H "Invalid event ID"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /api/v1/events/{id}/attendees [get]
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
// This allows event owners to remove attendees from their events.
// Only the event owner can remove attendees from their events.
//
// @Summary Remove attendee from event
// @Description Remove a user's registration from an event. Only the event owner can remove attendees.
// @Tags Events
// @Accept json
// @Produce json
// @Param id path int true "Event ID"
// @Param userId path int true "User ID to remove from attendees"
// @Success 204 "Attendee removed successfully"
// @Failure 400 {object} gin.H "Invalid event or user ID"
// @Failure 401 {object} gin.H "Unauthorized"
// @Failure 403 {object} gin.H "Forbidden - not event owner"
// @Failure 404 {object} gin.H "Event not found"
// @Failure 500 {object} gin.H "Internal server error"
// @Security BearerAuth
// @Router /api/v1/events/{id}/attendees/{userId} [delete]
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

	event, err := app.models.Events.Get(eventId)
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve event"})
		return
	}

	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	user := app.getUserFromContext(c)
	if event.OwnerId != user.Id{
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to delete an attendee from an event"})
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
// @Summary Get events by attendee
// @Description Retrieve all events that a specific user is registered for as an attendee
// @Tags Events
// @Accept json
// @Produce json
// @Param id path int true "Event ID"
// @Param userId path int true "User/Attendee ID"
// @Success 200 {array} database.Event "List of events user is attending"
// @Failure 400 {object} gin.H "Invalid attendee ID"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /api/v1/events/{id}/attendees/{userId} [get]
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