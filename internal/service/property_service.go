package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"digital-greenhouse/greenhouse-be/internal/domain"
)

type propertyService struct {
	repo        domain.PropertyRepository
	userRepo    domain.UserRepository
	bookingRepo domain.BookingRepository
}

func NewPropertyService(repo domain.PropertyRepository, userRepo domain.UserRepository, bookingRepo domain.BookingRepository) domain.PropertyService {
	return &propertyService{
		repo:        repo,
		userRepo:    userRepo,
		bookingRepo: bookingRepo,
	}
}

func (s *propertyService) CreateProperty(ctx context.Context, p *domain.Property) error {
	if p.Name == "" || p.OwnerID == 0 {
		return errors.New("nombre y dueño son requeridos")
	}

	// 1. Crear la propiedad
	if p.Status == "" {
		p.Status = domain.PropertyActive
	}

	if err := s.repo.Create(ctx, p); err != nil {
		return err
	}

	// 2. Lógica de cambio de rol: CLIENT -> OWNER
	user, err := s.userRepo.GetByID(ctx, p.OwnerID)
	if err != nil {
		return err
	}

	if user.Role == domain.RoleClient {
		user.Role = domain.RoleOwner
		if err := s.userRepo.Update(ctx, user); err != nil {
			// Nota: En una implementación real, esto debería estar dentro de una transacción.
			// Para este proyecto, seguiremos la estructura simple.
			return err
		}
	}

	// 3. Si vienen imágenes iniciales, guardarlas
	for i := range p.Images {
		p.Images[i].PropertyID = p.ID
		if err := s.repo.AddImage(ctx, &p.Images[i]); err != nil {
			return err
		}
	}

	return nil
}

func (s *propertyService) ListProperties(ctx context.Context, filter domain.PropertyFilter) ([]domain.Property, error) {
	return s.repo.GetAll(ctx, filter)
}

func (s *propertyService) GetPropertyByID(ctx context.Context, id uint) (*domain.Property, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *propertyService) GetPropertiesByOwner(ctx context.Context, ownerID uint) ([]domain.Property, error) {
	return s.repo.GetByOwnerID(ctx, ownerID)
}

func (s *propertyService) UpdateProperty(ctx context.Context, p *domain.Property) error {
	if p.ID == 0 {
		return errors.New("ID de propiedad requerido")
	}
	return s.repo.Update(ctx, p)
}

func (s *propertyService) DeleteProperty(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

func (s *propertyService) AddPropertyImage(ctx context.Context, img *domain.PropertyImage) error {
	if img.PropertyID == 0 || img.ImageData == "" {
		return errors.New("ID de propiedad y datos de imagen requeridos")
	}
	return s.repo.AddImage(ctx, img)
}

func (s *propertyService) UpdatePropertyImage(ctx context.Context, img *domain.PropertyImage) error {
	if img.ID == 0 {
		return errors.New("ID de imagen requerido")
	}
	return s.repo.UpdateImage(ctx, img)
}

func (s *propertyService) DeletePropertyImage(ctx context.Context, imageID uint) error {
	return s.repo.DeleteImage(ctx, imageID)
}

func (s *propertyService) CreatePricingRule(ctx context.Context, rule *domain.PricingRule) error {
	if rule.PropertyID == 0 || rule.StartDate.IsZero() || rule.EndDate.IsZero() {
		return errors.New("propiedad, fecha de inicio y fin son requeridas")
	}
	if rule.PriceModifier <= 0 {
		return errors.New("el modificador de precio debe ser mayor a 0")
	}
	return s.bookingRepo.CreatePricingRule(ctx, rule)
}

func (s *propertyService) ListPricingRulesByProperty(ctx context.Context, propertyID uint) ([]domain.PricingRule, error) {
	return s.bookingRepo.GetAllPricingRulesByPropertyID(ctx, propertyID)
}

func (s *propertyService) DeletePricingRule(ctx context.Context, id uint) error {
	return s.bookingRepo.DeletePricingRule(ctx, id)
}

func (s *propertyService) AutoGenerateHighSeasonRules(ctx context.Context, propertyID uint) error {
	currentYear := time.Now().Year()
	// Generar para el año actual y el siguiente
	years := []int{currentYear, currentYear + 1}
	highSeasonMonths := []struct {
		name  string
		month time.Month
		start int
		end   int
	}{
		{"Junio Alta", time.June, 1, 30},
		{"Julio Alta", time.July, 1, 31},
		{"Diciembre Alta", time.December, 1, 31},
		{"Enero Alta", time.January, 1, 31},
	}

	for _, year := range years {
		for _, m := range highSeasonMonths {
			rule := &domain.PricingRule{
				PropertyID:    propertyID,
				Name:          fmt.Sprintf("%s %d", m.name, year),
				StartDate:     time.Date(year, m.month, m.start, 0, 0, 0, 0, time.UTC),
				EndDate:       time.Date(year, m.month, m.end, 23, 59, 59, 0, time.UTC),
				PriceModifier: 1.10, // 10% de aumento
				Description:   "Generado automáticamente: Temporada Alta",
				IsActive:      true,
			}
			if err := s.bookingRepo.CreatePricingRule(ctx, rule); err != nil {
				return err
			}
		}
	}

	return nil
}
