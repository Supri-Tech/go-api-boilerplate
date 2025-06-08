package infrastructure

import (
	"go-crud-api/m/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func SetupRoutes(r chi.Router, c *Container) {
	r.Route("/api", func(r chi.Router) {
		r.Post("/login", c.UserHandler.Login)

		r.Group(func(r chi.Router) {
			r.Use(middleware.JWTAuth)
			r.Use(middleware.AdminOnly)
			r.Post("/register", c.UserHandler.Register)
		})

		r.Group(func(r chi.Router) {
			r.Use(middleware.JWTAuth)
			r.Get("/profile", c.UserHandler.Profile)

			r.Route("/products", func(r chi.Router) {
				r.Get("/", c.ProductHandler.GetProduct)
				r.Get("/{id}", c.ProductHandler.GetProductByID)
				r.Post("/", c.ProductHandler.CreateProduct)
				r.Put("/{id}", c.ProductHandler.UpdateProduct)
				r.Delete("/{id}", c.ProductHandler.DeleteProduct)
			})
		})
	})
}
