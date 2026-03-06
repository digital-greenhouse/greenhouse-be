package route

import (
	"digital-greenhouse/greenhouse-be/internal/http/handler"

	"github.com/go-chi/chi/v5"
)

func RegisterUserRoutes(r chi.Router, userHandler *handler.UserHandler) {
	r.Route("/users", func(r chi.Router) {
		r.Post("/", userHandler.CreateUser)
		r.Get("/", userHandler.GetUsers)
		r.Get("/{id}", userHandler.GetUser)
		r.Put("/{id}", userHandler.UpdateUser)
		r.Delete("/{id}", userHandler.DeleteUser)
	})
}
