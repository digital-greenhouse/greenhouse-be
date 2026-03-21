package service

import (
	"context"
	"errors"

	"digital-greenhouse/greenhouse-be/internal/domain"
)

type propertyService struct {
	repo     domain.PropertyRepository
	userRepo domain.UserRepository
}

func NewPropertyService(repo domain.PropertyRepository, userRepo domain.UserRepository) domain.PropertyService {
	return &propertyService{
		repo:     repo,
		userRepo: userRepo,
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
