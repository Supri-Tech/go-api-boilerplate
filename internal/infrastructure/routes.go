package infrastructure

import (
	"go-crud-api/m/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func SetupRoutes(r chi.Router, c *Container) {
	r.Route("/api", func(r chi.Router) {
		r.Post("/login", c.UserHandler.Login)
		r.Post("/register", c.UserHandler.Register)

		r.Group(func(r chi.Router) {
			r.Use(middleware.JWTAuth)
			r.Get("/profile", c.UserHandler.Profile)
		})
	})
}
