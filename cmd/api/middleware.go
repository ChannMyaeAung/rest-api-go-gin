package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// AuthMiddleware creates a Gin middleware function for JWT-based authentication.
// This middleware validates JWT tokens in the Authorization header and ensures only authenticated users
// can access protected endpoints. Upon successful validation, it loads the user data and makes it
// available to subsequent handlers via the Gin context.
//
// The middleware expects the Authorization header in the format: "Bearer <jwt_token>"
// and validates the token against the application's JWT secret.
// Authentication flow:
//   1. Extracts Authorization header from the request
//   2. Validates Bearer token format
//   3. Parses and validates the JWT token signature and expiration
//   4. Extracts user ID from token claims
//   5. Loads user data from database to verify user still exists
//   6. Stores user object in Gin context for use by protected handlers
//
// On authentication failure:
//   - Returns 401 Unauthorized with error message
//   - Calls c.Abort() to stop the request chain
//
// On authentication success:
//   - Sets "user" key in Gin context with database.User object
//   - Calls c.Next() to continue to the protected handler
func (app *application) AuthMiddleware() gin.HandlerFunc{
	return func(c *gin.Context){
		// Extract the Authorization header from the HTTP request
        // Standard format: "Authorization: Bearer <token>"
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort() // Stop request processing chain
			return
		}

		// Extract the token part by removing the "Bearer " prefix
        // TrimPrefix only removes the prefix if it exists, otherwise returns original string
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Check if the prefix was actually removed (indicates proper Bearer format)
        // If tokenString equals authHeader, then "Bearer " prefix was not found
		if tokenString == authHeader{
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token is required"})
			c.Abort()
			return
		}

		// Parse and validate the JWT token
        // The key function validates the signing method and returns the secret for verification
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error){
			// Ensure the token uses HMAC signing method (HS256, HS384, HS512)
            // This prevents algorithm substitution attacks
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok{
				return nil, jwt.ErrSignatureInvalid
			}
			// Return the application's JWT secret for signature verification
			return []byte(app.jwtSecret), nil
		})

		// Check for parsing errors or invalid token (expired, malformed signature, etc.)
		if err != nil || !token.Valid{
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Extract claims from the validated token
        // Claims contain the payload data (user ID, expiration, etc.)
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Extract user ID from token claims
        // JWT stores numbers as float64, so we need to cast appropriately
        // This assumes the token was created with "userId" claim during login
		userId := claims["userId"].(float64)

		// Load the user from database to ensure they still exist and are active
        // This catches cases where user accounts have been deleted/disabled after token issuance
		user, err := app.models.Users.GetUserByID(int(userId))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
			c.Abort()
			return
		}

		// Additional check: ensure user was actually found in database
        // GetUserByID might return nil user with no error for non-existent users
        if user == nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
            c.Abort()
            return
        }

		// Store the authenticated user in the Gin context
        // This makes the user object available to all subsequent handlers in the request chain
        // Handlers can retrieve it with: user := c.MustGet("user").(*database.User)
		c.Set("user", user)

		// Continue to the next handler in the middleware chain
        // This allows the protected route handler to execute
		c.Next()
	}
}