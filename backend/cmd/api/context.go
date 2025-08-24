package main

import (
	"rest-api-in-gin/internal/database"

	"github.com/gin-gonic/gin"
)

// getUserFromContext extracts the authenticated user from the Gin context.
// This is a helper function used by protected route handlers to access the current user's data.
// The user object is stored in the context by the AuthMiddleware after successful JWT validation.
//
// The function performs safe type assertions and provides fallback behavior to prevent panics
// when the user object is missing or has an unexpected type.
func (app *application) getUserFromContext(c *gin.Context) *database.User{
	contextUser, exists := c.Get("user")
	if !exists{
		// Return empty user struct if the key doesn't exist
        // This prevents nil pointer dereferences in calling code
		return &database.User{} 
	}

	// Perform type assertion to convert interface{} to *database.User
	user, ok := contextUser.(*database.User)
	if !ok {
		return &database.User{} // Return empty user if type assertion fails
	}	

	// Return the validated user object
	return user
}