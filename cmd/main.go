package main

import (
	"go-crud-api/m/internal/container"
	"go-crud-api/m/internal/db"
	"go-crud-api/m/internal/routes"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	port := os.Getenv("APP_PORT")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	database := db.NewMySQL()
	container := container.NewContainer(database)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "API is working"}`))
	})

	routes.SetupRoutes(r, container)

	log.Println("Server running at :" + port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}
