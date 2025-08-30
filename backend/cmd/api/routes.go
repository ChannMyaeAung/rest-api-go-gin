package main

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (app *application) routes() http.Handler {
    g := gin.Default()

    g.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }))

    // Health + auth (public)
    v1 := g.Group("/api/v1")
    {
        v1.GET("/health", func(c *gin.Context) {
            c.JSON(http.StatusOK, gin.H{"status": "ok"})
        })
        v1.POST("/auth/register", app.registerUser)
        v1.POST("/auth/login", app.login)
    }

    // Protected group (requires JWT)
    auth := v1.Group("/")
    auth.Use(app.AuthMiddleware())
    {
        auth.GET("/users/:id", app.getUserByID)
        auth.GET("/events", app.getAllEvents)
        auth.GET("/events/:id", app.getEventByID)

        auth.GET("/events/:id/attendees", app.getAttendeesForEvent)
        auth.GET("/events/:id/attendees/:userId", app.getEventsByAttendee)
        auth.POST("/events/:id/attendees/:userId", app.addAttendeeToEvent)
        auth.DELETE("/events/:id/attendees/:userId", app.deleteAttendeeFromEvent)

        // Event mutations
        auth.POST("/events", app.createEvent)
        auth.PUT("/events/:id", app.updateEvent)
        auth.DELETE("/events/:id", app.deleteEvent)
    }

    // Swagger documentation
    g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    return g
}