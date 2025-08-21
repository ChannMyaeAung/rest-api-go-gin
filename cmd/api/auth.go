package main

import (
	"net/http"
	"rest-api-in-gin/internal/database"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type registerRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required,min=2"`
}

type loginRequest struct{
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// loginResponse defines the JSON structure returned upon successful authentication.
// Contains the JWT token that clients should include in subsequent authenticated requests.
type loginResponse struct{
	Token string `json:"token"`
}

// login handles POST /auth/login requests for user authentication.
// It validates user credentials against the database and returns a JWT token upon success.
// The returned token should be included in the Authorization header for protected endpoints.
//
// Request body should contain JSON with email and password fields.
// Password is compared against the stored bcrypt hash for security.
//
// Returns:
//   - 200 OK with JWT token on successful authentication
//   - 400 Bad Request for invalid JSON structure or validation failures
//   - 401 Unauthorized for invalid credentials (wrong email/password)
//   - 500 Internal Server Error on database or token generation failures
//
// Security considerations:
//   - Uses bcrypt for secure password comparison
//   - Returns generic error message to prevent email enumeration attacks
//   - JWT token expires after 72 hours
func (app *application) login(c *gin.Context){
	var auth loginRequest
	if err := c.ShouldBindJSON(&auth); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existingUser, err := app.models.Users.GetByEmail(auth.Email)
	if existingUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return 
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		return
	}

	// Compare provided password with stored bcrypt hash
    // bcrypt.CompareHashAndPassword is constant-time to prevent timing attacks
	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(auth.Password))	
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return 
	}
	
	// Create JWT token with user ID and expiration claims
    // Token expires in 72 hours for reasonable session duration
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": existingUser.Id,
		"exp":    time.Now().Add(time.Hour * 72).Unix(),
	})
	
	// Sign the token with the application's secret key
	tokenString, err := token.SignedString([]byte(app.jwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return 
	}

	// Return the signed JWT token for client-side storage and future requests
	c.JSON(http.StatusOK, loginResponse{Token: tokenString})
}

// registerUser handles POST /auth/register requests for new user account creation.
// It validates input data, hashes the password securely, and creates a new user record.
// Passwords are hashed using bcrypt before storage for security.
//
// Request body should contain JSON with email, password, and name fields.
// Email must be unique (enforced by database constraints).
//
// Returns:
//   - 201 Created with success message and user data (password excluded from response)
//   - 400 Bad Request for invalid JSON structure or validation failures
//   - 500 Internal Server Error on password hashing or database failures
//
// Security considerations:
//   - Passwords are hashed with bcrypt using default cost (currently 10)
//   - User password is excluded from JSON response via struct tags
//   - Email uniqueness should be enforced by database constraints
func (app *application) registerUser(c *gin.Context) {
	var register registerRequest
	if err := c.ShouldBindJSON(&register); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return 
	}

	// Hash the plaintext password using bcrypt with default cost factor (10)
    // bcrypt is adaptive and includes salt generation automatically
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(register.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	register.Password = string(hashPassword)
	user := database.User{
		Email: register.Email,
		Password: register.Password,
		Name: register.Name,
	}

	err = app.models.Users.Insert(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return 
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "user": user})
}