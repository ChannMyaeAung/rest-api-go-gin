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
func (app *application) createEvent(c *gin.Context) {
	var event database.Event

	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := app.getUserFromContext(c)
	event.OwnerId = user.Id

	if err := app.models.Events.Insert(&event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create event"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Event created successfully", "event": event})
}

// getAllEvents returns events scoped to the authenticated owner.
//
// @Summary Returns all events for the current user
// @Description Returns all events owned by the authenticated user
// @Tags Events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} []database.Event
// @Router /api/v1/events [get]
func (app *application) getAllEvents(c *gin.Context) {
	user := app.getUserFromContext(c)
	events, err := app.models.Events.GetAllByOwner(user.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve events"})
		return
	}

	c.JSON(http.StatusOK, events)
}

// getEventByID handles GET /events/:id requests.
//
// @Summary Get event by ID
// @Description Retrieve a specific event by its ID
// @Tags Events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Event ID"
// @Success 200 {object} database.Event "Event details"
// @Failure 400 {object} gin.H "Invalid event ID"
// @Failure 404 {object} gin.H "Event not found"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /api/v1/events/{id} [get]
func (app *application) getEventByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	user := app.getUserFromContext(c)
	event, err := app.models.Events.Get(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve event"})
		return
	}
	if event == nil || event.OwnerId != user.Id {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	c.JSON(http.StatusOK, event)
}

// updateEvent handles PUT /events/:id to update an existing event.
func (app *application) updateEvent(c *gin.Context) {
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
	if existingEvent == nil || existingEvent.OwnerId != user.Id {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	updatedEvent := &database.Event{}
	if err := c.ShouldBindJSON(updatedEvent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedEvent.Id = id
	updatedEvent.OwnerId = user.Id

	if err := app.models.Events.Update(updatedEvent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update event"})
		return
	}

	c.JSON(http.StatusOK, updatedEvent)
}

// deleteEvent handles DELETE /events/:id requests.
func (app *application) deleteEvent(c *gin.Context) {
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
	if existingEvent == nil || existingEvent.OwnerId != user.Id {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	if err := app.models.Events.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete event"})
		return
	}

	c.Status(http.StatusNoContent)
}

// addAttendeeToEvent handles POST /events/:id/attendees/:userId.
func (app *application) addAttendeeToEvent(c *gin.Context) {
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
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve event"})
		return
	}
	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	user := app.getUserFromContext(c)
	if event.OwnerId != user.Id {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to add attendees to this event"})
		return
	}

	userToAdd, err := app.models.Users.GetUserByID(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		return
	}
	if userToAdd == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	existingAttendee, err := app.models.Attendees.GetByEventAndAttendee(event.Id, userToAdd.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve attendee"})
		return
	}
	if existingAttendee != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User is already an attendee"})
		return
	}

	attendee := database.Attendee{EventId: event.Id, UserId: userToAdd.Id}
	if _, err := app.models.Attendees.Insert(&attendee); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add attendee"})
		return
	}

	c.JSON(http.StatusCreated, attendee)
}

// getAttendeesForEvent handles GET /events/:id/attendees.
func (app *application) getAttendeesForEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	event, err := app.models.Events.Get(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve event"})
		return
	}
	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	user := app.getUserFromContext(c)
	if event.OwnerId != user.Id {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to view attendees for this event"})
		return
	}

	users, err := app.models.Attendees.GetAttendeesByEvent(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve attendees"})
		return
	}

	c.JSON(http.StatusOK, users)
}

// deleteAttendeeFromEvent handles DELETE /events/:id/attendees/:userId.
func (app *application) deleteAttendeeFromEvent(c *gin.Context) {
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
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve event"})
		return
	}
	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	user := app.getUserFromContext(c)
	if event.OwnerId != user.Id {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to delete an attendee from this event"})
		return
	}

	if err := app.models.Attendees.DeleteByEventAndUser(userId, eventId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete attendee"})
		return
	}

	c.Status(http.StatusNoContent)
}

// getEventsByAttendee handles GET /events/:id/attendees/:userId requests to retrieve all events for a user.
func (app *application) getEventsByAttendee(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid attendee ID"})
		return
	}

	events, err := app.models.Attendees.GetEventsByAttendee(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve events"})
		return
	}

	c.JSON(http.StatusOK, events)
}
