package main

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// getUserByID handles GET /users/:id requests to retrieve a specific user by ID.
// Returns user information excluding the password field for security.
//
// @Summary Get user by ID
// @Description Retrieve a specific user's information by their ID
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} database.User "User information"
// @Failure 400 {object} gin.H "Invalid user ID"
// @Failure 404 {object} gin.H "User not found"
// @Router /api/v1/users/{id} [get]
func (app *application) getUserByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	user, err := app.models.Users.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (app *application) getCurrentUser(c *gin.Context) {
	user := app.getUserFromContext(c)
	if user == nil || user.Id == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (app *application) deleteCurrentUser(c *gin.Context) {
	user := app.getUserFromContext(c)
	if user == nil || user.Id == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := app.models.Users.Delete(c.Request.Context(), user.Id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.Status(http.StatusNoContent)
}
