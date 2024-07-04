package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	dbUrl := os.Getenv("DB_URL")//dsn
	if port == "" {
		log.Fatal("Port not found in environment")
	}
	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to database")
	}
	fmt.Println("Connected to database")

	// Auto Migrate the User model
	err = db.AutoMigrate(&User{})
	if err != nil {
		log.Fatal("Error migrating User model")
	}
	fmt.Println("Migration Completed")

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	r.Use(middleware.Logger)

	app := &App{DB: db}

	v1Router := chi.NewRouter()
	v1Router.Get("/", app.handlerReadiness)
	v1Router.Get("/error", app.handlerError)
	v1Router.Post("/users/create", app.handlerCreateUser)
	r.Mount("/v1", v1Router)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	fmt.Println("Server running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
