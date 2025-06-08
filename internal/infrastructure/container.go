package infrastructure

import (
	"database/sql"
	"go-crud-api/m/internal/user"
)

type Container struct {
	UserHandler *user.Handler
}

func NewContainer(db *sql.DB) *Container {
	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

	return &Container{
		UserHandler: userHandler,
	}
}
