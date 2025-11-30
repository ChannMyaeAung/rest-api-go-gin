package main

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	_ "rest-api-in-gin/docs"
	"rest-api-in-gin/internal/database"
	"rest-api-in-gin/internal/env"

	_ "github.com/joho/godotenv/autoload"
	_ "modernc.org/sqlite"
)

// @title Go Gin Rest API
// @version 1.0
// @description A rest API in Go using Gin framework
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter your bearer token in the format **Bearer &lt;token&gt;**

type application struct {
	port      int
	jwtSecret string
	uploadDir string
	models    database.Models
}

func main() {
	dbPath := env.GetEnvString("DATABASE_PATH", "./data.db")
	if err := ensureDir(filepath.Dir(dbPath)); err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	models := database.NewModels(db)
	uploadDir := env.GetEnvString("UPLOAD_DIR", "./tmp/uploads")
	if err := ensureDir(uploadDir); err != nil {
		log.Fatal(err)
	}
	app := &application{
		port:      env.GetEnvInt("PORT", 8080),
		jwtSecret: env.GetEnvString("JWT_SECRET", "some-secret-123456"),
		uploadDir: uploadDir,
		models:    models,
	}

	if err := app.serve(); err != nil {
		log.Fatal(err)
	}
}

func ensureDir(path string) error {
	if path == "" || path == "." {
		return nil
	}
	return os.MkdirAll(path, 0o755)
}
