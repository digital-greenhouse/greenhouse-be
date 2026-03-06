package service

import (
	"context"
	"errors"
	"fmt"

	"digital-greenhouse/greenhouse-be/internal/domain"

	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repo domain.UserRepository
}

func NewUserService(repo domain.UserRepository) domain.UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(ctx context.Context, user *domain.User) error {
	if user.Email == "" || user.PasswordHash == "" || user.Name == "" {
		return errors.New("name, email y password son requeridos")
	}

	existing, _ := s.repo.GetByEmail(ctx, user.Email)
	if existing != nil {
		return errors.New("el email ya está registrado")
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error al hashear la contraseña: %w", err)
	}
	user.PasswordHash = string(hashed)

	if user.Role == "" {
		user.Role = domain.RoleClient
	}

	return s.repo.Create(ctx, user)
}

func (s *userService) GetUserByID(ctx context.Context, id uint) (*domain.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *userService) GetAllUsers(ctx context.Context) ([]domain.User, error) {
	return s.repo.GetAll(ctx)
}

func (s *userService) UpdateUser(ctx context.Context, user *domain.User) error {
	if user.ID == 0 {
		return errors.New("el ID del usuario es requerido para actualizar")
	}
	return s.repo.Update(ctx, user)
}

func (s *userService) DeleteUser(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
