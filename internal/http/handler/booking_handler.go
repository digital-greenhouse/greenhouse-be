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

type BookingHandler struct {
	service domain.BookingService
}

func NewBookingHandler(service domain.BookingService) *BookingHandler {
	return &BookingHandler{service: service}
}

func (h *BookingHandler) CreateQuote(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateQuoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errResponse(w, http.StatusBadRequest, "payload inválido")
		return
	}

	// El clientID es opcional para cotizaciones (puede ser invitado)
	var clientID *uint
	uID := middleware.GetUserID(r.Context())
	if uID != 0 {
		clientID = &uID
	}

	quote, err := h.service.CalculateQuote(r.Context(), req.PropertyID, clientID, req.CheckInDate, req.CheckOutDate, req.GuestCount)
	if err != nil {
		errResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	jsonResponse(w, http.StatusCreated, dto.ToQuoteResponse(quote))
}

func (h *BookingHandler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateBookingFromQuoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errResponse(w, http.StatusBadRequest, "payload inválido")
		return
	}

	// Requerimos clientID (extraído del contexto por el middleware)
	clientID := middleware.GetUserID(r.Context())
	if clientID == 0 {
		errResponse(w, http.StatusUnauthorized, "se requiere autenticación")
		return
	}

	booking, err := h.service.CreateBookingFromQuote(r.Context(), req.QuoteID, clientID, req.SpecialRequests)
	if err != nil {
		errResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	jsonResponse(w, http.StatusCreated, dto.ToBookingResponse(booking))
}

func (h *BookingHandler) GetMyHistory(w http.ResponseWriter, r *http.Request) {
	clientID := middleware.GetUserID(r.Context())
	if clientID == 0 {
		errResponse(w, http.StatusUnauthorized, "se requiere autenticación")
		return
	}

	bookings, err := h.service.GetClientHistory(r.Context(), clientID)
	if err != nil {
		errResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := make([]dto.BookingResponse, len(bookings))
	for i := range bookings {
		resp[i] = dto.ToBookingResponse(&bookings[i])
	}

	jsonResponse(w, http.StatusOK, resp)
}

func (h *BookingHandler) CancelBooking(w http.ResponseWriter, r *http.Request) {
	bookingID, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		errResponse(w, http.StatusBadRequest, "ID de reserva inválido")
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	if err := h.service.CancelBooking(r.Context(), uint(bookingID), req.Reason); err != nil {
		errResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	jsonResponse(w, http.StatusOK, map[string]string{"message": "reserva cancelada exitosamente"})
}
