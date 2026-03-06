package repository

import (
	"context"
	"time"

	"digital-greenhouse/greenhouse-be/internal/domain"

	"gorm.io/gorm"
)

type UserDBModel struct {
	ID           uint      `gorm:"primaryKey;autoIncrement"`
	Name         string    `gorm:"type:varchar(100);not null"`
	Email        string    `gorm:"type:varchar(150);uniqueIndex;not null"`
	PasswordHash string    `gorm:"type:varchar(255);not null"`
	Role         string    `gorm:"type:enum('SUPERADMIN','OWNER','CLIENT');default:'CLIENT'"`
	Phone        *string   `gorm:"type:varchar(20)"`
	IsActive     bool      `gorm:"default:true"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}

func (UserDBModel) TableName() string {
	return "users"
}

func toDomainUser(dbUser UserDBModel) *domain.User {
	return &domain.User{
		ID:           dbUser.ID,
		Name:         dbUser.Name,
		Email:        dbUser.Email,
		PasswordHash: dbUser.PasswordHash,
		Role:         domain.Role(dbUser.Role),
		Phone:        dbUser.Phone,
		IsActive:     dbUser.IsActive,
		CreatedAt:    dbUser.CreatedAt,
		UpdatedAt:    dbUser.UpdatedAt,
	}
}

func fromDomainUser(domainUser *domain.User) UserDBModel {
	return UserDBModel{
		ID:           domainUser.ID,
		Name:         domainUser.Name,
		Email:        domainUser.Email,
		PasswordHash: domainUser.PasswordHash,
		Role:         string(domainUser.Role),
		Phone:        domainUser.Phone,
		IsActive:     domainUser.IsActive,
		CreatedAt:    domainUser.CreatedAt,
		UpdatedAt:    domainUser.UpdatedAt,
	}
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	dbModel := fromDomainUser(user)
	err := r.db.WithContext(ctx).Create(&dbModel).Error
	if err == nil {
		user.ID = dbModel.ID
		user.CreatedAt = dbModel.CreatedAt
		user.UpdatedAt = dbModel.UpdatedAt
	}
	return err
}

func (r *userRepository) GetByID(ctx context.Context, id uint) (*domain.User, error) {
	var dbModel UserDBModel
	err := r.db.WithContext(ctx).First(&dbModel, id).Error
	if err != nil {
		return nil, err
	}
	return toDomainUser(dbModel), nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var dbModel UserDBModel
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&dbModel).Error
	if err != nil {
		return nil, err
	}
	return toDomainUser(dbModel), nil
}

func (r *userRepository) GetAll(ctx context.Context) ([]domain.User, error) {
	var dbModels []UserDBModel
	err := r.db.WithContext(ctx).Find(&dbModels).Error
	if err != nil {
		return nil, err
	}

	var users []domain.User
	for _, dbModel := range dbModels {
		users = append(users, *toDomainUser(dbModel))
	}
	return users, nil
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	dbModel := fromDomainUser(user)
	err := r.db.WithContext(ctx).Save(&dbModel).Error
	if err == nil {
		user.UpdatedAt = dbModel.UpdatedAt
	}
	return err
}

func (r *userRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&UserDBModel{}, id).Error
}
