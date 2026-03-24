package route

import (
	"digital-greenhouse/greenhouse-be/internal/http/handler"
	"digital-greenhouse/greenhouse-be/internal/http/middleware"

	"github.com/go-chi/chi/v5"
)

func RegisterPropertyRoutes(r chi.Router, propertyHandler *handler.PropertyHandler) {
	r.Route("/properties", func(r chi.Router) {
		// Rutas públicas
		r.Get("/", propertyHandler.ListProperties)
		r.Get("/{id}", propertyHandler.GetPropertyByID)
		r.Get("/owner/{id}", propertyHandler.GetPropertiesByOwner)

		// Rutas protegidas
		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware)
			
			r.Post("/", propertyHandler.CreateProperty)
			r.Post("/{id}/images", propertyHandler.AddImage)
			r.Delete("/images/{imageID}", propertyHandler.DeleteImage)
		})
	})
}
