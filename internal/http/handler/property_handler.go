package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"digital-greenhouse/greenhouse-be/internal/domain"
	"digital-greenhouse/greenhouse-be/internal/http/dto"

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
	var req dto.CreatePropertyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errResponse(w, http.StatusBadRequest, "payload inválido")
		return
	}

	property := &domain.Property{
		OwnerID:           req.OwnerID,
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

func (h *PropertyHandler) GetPropertiesByOwner(w http.ResponseWriter, r *http.Request) {
	ownerID, err := strconv.ParseUint(chi.URLParam(r, "ownerID"), 10, 32)
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
