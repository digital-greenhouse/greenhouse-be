package route

import (
	"digital-greenhouse/greenhouse-be/internal/http/handler"
	"digital-greenhouse/greenhouse-be/internal/http/middleware"

	"github.com/go-chi/chi/v5"
)

func RegisterBookingRoutes(r chi.Router, bookingHandler *handler.BookingHandler) {
	r.Route("/bookings", func(r chi.Router) {
		// Ruta con autenticación opcional (puede ser invitado)
		r.With(middleware.OptionalAuthMiddleware).Post("/quote", bookingHandler.CreateQuote)

		// Grupo con middleware estricto aplicado al resto
		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware)
			
			r.Post("/", bookingHandler.CreateBooking)
			r.Get("/history", bookingHandler.GetMyHistory)
			r.Post("/{id}/cancel", bookingHandler.CancelBooking)
		})
	})
}
