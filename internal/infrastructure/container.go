package infrastructure

import (
	"database/sql"
	"go-crud-api/m/internal/product"
	"go-crud-api/m/internal/user"
)

type Container struct {
	UserHandler    *user.Handler
	ProductHandler *product.Handler
}

func NewContainer(db *sql.DB) *Container {
	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

	productRepo := product.NewRepository(db)
	productService := product.NewService(productRepo)
	productHandler := product.NewHandler(productService)

	return &Container{
		UserHandler:    userHandler,
		ProductHandler: productHandler,
	}
}
