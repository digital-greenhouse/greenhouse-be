package route

import (
	"digital-greenhouse/greenhouse-be/internal/http/handler"

	"github.com/go-chi/chi/v5"
)

func RegisterPropertyRoutes(r chi.Router, h *handler.PropertyHandler) {
	r.Route("/properties", func(r chi.Router) {
		r.Post("/", h.CreateProperty)
		r.Get("/owner/{ownerID}", h.GetPropertiesByOwner)
		r.Post("/{id}/images", h.AddImage)
		r.Delete("/images/{imageID}", h.DeleteImage)
	})
}
