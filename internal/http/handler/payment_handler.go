package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"digital-greenhouse/greenhouse-be/internal/domain"
	"digital-greenhouse/greenhouse-be/internal/http/dto"
	"digital-greenhouse/greenhouse-be/internal/http/middleware"

	"github.com/go-chi/chi/v5"
)

type PaymentHandler struct {
	service domain.PaymentService
}

func NewPaymentHandler(service domain.PaymentService) *PaymentHandler {
	return &PaymentHandler{service: service}
}

func (h *PaymentHandler) UploadProof(w http.ResponseWriter, r *http.Request) {
	var req dto.UploadPaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errResponse(w, http.StatusBadRequest, "payload inválido")
		return
	}

	payment, err := h.service.ProcessPaymentProof(
		r.Context(),
		req.BookingID,
		req.Amount,
		req.PaymentMethod,
		req.ProofData,
		req.ProofMimeType,
	)
	if err != nil {
		errResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	jsonResponse(w, http.StatusCreated, dto.ToPaymentResponse(payment))
}

func (h *PaymentHandler) VerifyPayment(w http.ResponseWriter, r *http.Request) {
	paymentID, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		errResponse(w, http.StatusBadRequest, "ID de pago inválido")
		return
	}

	var req dto.VerifyPaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errResponse(w, http.StatusBadRequest, "payload inválido")
		return
	}

	verifierID := middleware.GetUserID(r.Context())
	if verifierID == 0 {
		errResponse(w, http.StatusUnauthorized, "se requiere autenticación")
		return
	}

	err = h.service.VerifyPayment(
		r.Context(),
		uint(paymentID),
		verifierID,
		req.Status,
		req.RejectionReason,
	)
	if err != nil {
		errResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	jsonResponse(w, http.StatusOK, map[string]string{"message": "pago procesado exitosamente"})
}
