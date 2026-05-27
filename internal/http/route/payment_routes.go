package route

import (
	"digital-greenhouse/greenhouse-be/internal/http/handler"
	"digital-greenhouse/greenhouse-be/internal/http/middleware"

	"github.com/go-chi/chi/v5"
)

func RegisterPaymentRoutes(r chi.Router, paymentHandler *handler.PaymentHandler) {
	r.Route("/payments", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)

		// Cliente sube comprobante
		r.Post("/upload", paymentHandler.UploadProof)

		// Admin/Owner verifica pago
		r.Post("/{id}/verify", paymentHandler.VerifyPayment)

		// Descarga de comprobante
		r.Get("/{id}/proof", paymentHandler.DownloadProof)
	})
}
