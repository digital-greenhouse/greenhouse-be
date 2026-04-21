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

			// Reglas de Precios
			r.Post("/{id}/pricing-rules", propertyHandler.CreatePricingRule)
			r.Get("/{id}/pricing-rules", propertyHandler.GetPricingRules)
			r.Delete("/pricing-rules/{ruleId}", propertyHandler.DeletePricingRule)
			r.Post("/{id}/pricing-rules/auto-generate", propertyHandler.AutoGeneratePricingRules)
		})
	})
}
