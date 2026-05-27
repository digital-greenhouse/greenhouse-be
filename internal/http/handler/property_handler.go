package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"digital-greenhouse/greenhouse-be/internal/domain"
	"digital-greenhouse/greenhouse-be/internal/http/dto"
	"digital-greenhouse/greenhouse-be/internal/http/middleware"

	"github.com/go-chi/chi/v5"
)

type PropertyHandler struct {
	service domain.PropertyService
}

func NewPropertyHandler(service domain.PropertyService) *PropertyHandler {
	return &PropertyHandler{service: service}
}

func toPropertyResponse(p *domain.Property) dto.PropertyResponse {
	resp := dto.PropertyResponse{
		ID:                p.ID,
		OwnerID:           p.OwnerID,
		Name:              p.Name,
		Description:       p.Description,
		Address:           p.Address,
		BasePricePerNight: p.BasePricePerNight,
		MaxCapacity:       p.MaxCapacity,
		Status:            string(p.Status),
		CreatedAt:         p.CreatedAt.Format(time.RFC3339),
		UpdatedAt:         p.UpdatedAt.Format(time.RFC3339),
	}

	for _, img := range p.Images {
		resp.Images = append(resp.Images, dto.PropertyImageDTO{
			ID:         img.ID,
			PropertyID: img.PropertyID,
			ImageData:  img.ImageData,
			MimeType:   img.MimeType,
			AltText:    img.AltText,
			IsCover:    img.IsCover,
			SortOrder:  img.SortOrder,
		})
	}

	return resp
}

func (h *PropertyHandler) CreateProperty(w http.ResponseWriter, r *http.Request) {
	// Obtener ID del usuario del token
	ownerID := middleware.GetUserID(r.Context())
	if ownerID == 0 {
		errResponse(w, http.StatusUnauthorized, "se requiere autenticación")
		return
	}

	var req dto.CreatePropertyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errResponse(w, http.StatusBadRequest, "payload inválido")
		return
	}

	property := &domain.Property{
		OwnerID:           ownerID, // Sobrescribir siempre con el ID del token por seguridad
		Name:              req.Name,
		Description:       req.Description,
		Address:           req.Address,
		BasePricePerNight: req.BasePricePerNight,
		MaxCapacity:       req.MaxCapacity,
	}

	for _, img := range req.Images {
		property.Images = append(property.Images, domain.PropertyImage{
			ImageData: img.ImageData,
			MimeType:  img.MimeType,
			AltText:   img.AltText,
			IsCover:   img.IsCover,
			SortOrder: img.SortOrder,
		})
	}

	if err := h.service.CreateProperty(r.Context(), property); err != nil {
		errResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	jsonResponse(w, http.StatusCreated, toPropertyResponse(property))
}

func (h *PropertyHandler) ListProperties(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	filter := domain.PropertyFilter{
		Search:   q.Get("search"),
		Location: q.Get("location"),
	}

	if minPrice, err := strconv.ParseFloat(q.Get("min_price"), 64); err == nil {
		filter.MinPrice = minPrice
	}
	if maxPrice, err := strconv.ParseFloat(q.Get("max_price"), 64); err == nil {
		filter.MaxPrice = maxPrice
	}
	if guests, err := strconv.Atoi(q.Get("guests")); err == nil {
		filter.GuestCount = guests
	}

	// Parsing de fechas
	if checkInStr := q.Get("check_in"); checkInStr != "" {
		if t, err := time.Parse("2006-01-02", checkInStr); err == nil {
			filter.CheckInDate = &t
		}
	}
	if checkOutStr := q.Get("check_out"); checkOutStr != "" {
		if t, err := time.Parse("2006-01-02", checkOutStr); err == nil {
			filter.CheckOutDate = &t
		}
	}

	properties, err := h.service.ListProperties(r.Context(), filter)
	if err != nil {
		errResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := make([]dto.PropertyResponse, len(properties))
	for i := range properties {
		resp[i] = toPropertyResponse(&properties[i])
	}

	jsonResponse(w, http.StatusOK, resp)
}

func (h *PropertyHandler) GetPropertyByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		errResponse(w, http.StatusBadRequest, "ID de propiedad inválido")
		return
	}

	property, err := h.service.GetPropertyByID(r.Context(), uint(id))
	if err != nil {
		errResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	if property == nil {
		errResponse(w, http.StatusNotFound, "propiedad no encontrada")
		return
	}

	jsonResponse(w, http.StatusOK, toPropertyResponse(property))
}

func (h *PropertyHandler) GetPropertiesByOwner(w http.ResponseWriter, r *http.Request) {
	ownerID, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		errResponse(w, http.StatusBadRequest, "ID de dueño inválido")
		return
	}

	properties, err := h.service.GetPropertiesByOwner(r.Context(), uint(ownerID))
	if err != nil {
		errResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := make([]dto.PropertyResponse, len(properties))
	for i := range properties {
		resp[i] = toPropertyResponse(&properties[i])
	}

	jsonResponse(w, http.StatusOK, resp)
}

func (h *PropertyHandler) AddImage(w http.ResponseWriter, r *http.Request) {
	propertyID, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		errResponse(w, http.StatusBadRequest, "ID de propiedad inválido")
		return
	}

	var req dto.PropertyImageDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errResponse(w, http.StatusBadRequest, "payload inválido")
		return
	}

	img := &domain.PropertyImage{
		PropertyID: uint(propertyID),
		ImageData:  req.ImageData,
		MimeType:   req.MimeType,
		AltText:    req.AltText,
		IsCover:    req.IsCover,
		SortOrder:  req.SortOrder,
	}

	if err := h.service.AddPropertyImage(r.Context(), img); err != nil {
		errResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	jsonResponse(w, http.StatusCreated, req)
}

func (h *PropertyHandler) DeleteImage(w http.ResponseWriter, r *http.Request) {
	imageID, err := strconv.ParseUint(chi.URLParam(r, "imageID"), 10, 32)
	if err != nil {
		errResponse(w, http.StatusBadRequest, "ID de imagen inválido")
		return
	}

	if err := h.service.DeletePropertyImage(r.Context(), uint(imageID)); err != nil {
		errResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	jsonResponse(w, http.StatusNoContent, nil)
}

func (h *PropertyHandler) CreatePricingRule(w http.ResponseWriter, r *http.Request) {
	propertyID, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		errResponse(w, http.StatusBadRequest, "ID de propiedad inválido")
		return
	}

	var req dto.CreatePricingRuleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errResponse(w, http.StatusBadRequest, "payload inválido")
		return
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		errResponse(w, http.StatusBadRequest, "fecha de inicio inválida (use AAAA-MM-DD)")
		return
	}
	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		errResponse(w, http.StatusBadRequest, "fecha de fin inválida (use AAAA-MM-DD)")
		return
	}

	rule := &domain.PricingRule{
		PropertyID:    uint(propertyID),
		Name:          req.Name,
		StartDate:     startDate,
		EndDate:       endDate,
		PriceModifier: req.PriceModifier,
		Description:   req.Description,
		IsActive:      true,
	}

	if err := h.service.CreatePricingRule(r.Context(), rule); err != nil {
		errResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	jsonResponse(w, http.StatusCreated, toPricingRuleResponse(rule))
}

func (h *PropertyHandler) GetPricingRules(w http.ResponseWriter, r *http.Request) {
	propertyID, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		errResponse(w, http.StatusBadRequest, "ID de propiedad inválido")
		return
	}

	rules, err := h.service.ListPricingRulesByProperty(r.Context(), uint(propertyID))
	if err != nil {
		errResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := make([]dto.PricingRuleDTO, len(rules))
	for i := range rules {
		resp[i] = toPricingRuleResponse(&rules[i])
	}

	jsonResponse(w, http.StatusOK, resp)
}

func (h *PropertyHandler) DeletePricingRule(w http.ResponseWriter, r *http.Request) {
	ruleID, err := strconv.ParseUint(chi.URLParam(r, "ruleId"), 10, 32)
	if err != nil {
		errResponse(w, http.StatusBadRequest, "ID de regla inválido")
		return
	}

	if err := h.service.DeletePricingRule(r.Context(), uint(ruleID)); err != nil {
		errResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	jsonResponse(w, http.StatusNoContent, nil)
}

func (h *PropertyHandler) AutoGeneratePricingRules(w http.ResponseWriter, r *http.Request) {
	propertyID, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		errResponse(w, http.StatusBadRequest, "ID de propiedad inválido")
		return
	}

	if err := h.service.AutoGenerateHighSeasonRules(r.Context(), uint(propertyID)); err != nil {
		errResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	jsonResponse(w, http.StatusOK, map[string]string{"message": "reglas de temporada alta generadas exitosamente"})
}

func toPricingRuleResponse(r *domain.PricingRule) dto.PricingRuleDTO {
	return dto.PricingRuleDTO{
		ID:            r.ID,
		PropertyID:    r.PropertyID,
		Name:          r.Name,
		StartDate:     r.StartDate.Format("2006-01-02"),
		EndDate:       r.EndDate.Format("2006-01-02"),
		PriceModifier: r.PriceModifier,
		Description:   r.Description,
		IsActive:      r.IsActive,
	}
}
