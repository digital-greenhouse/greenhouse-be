package domain

import (
	"context"
	"time"
)

// Role define los roles posibles de un usuario
type Role string

const (
	RoleSuperAdmin Role = "SUPERADMIN"
	RoleOwner      Role = "OWNER"
	RoleClient     Role = "CLIENT"
)

type User struct {
	ID           uint
	Name         string
	Email        string
	PasswordHash string
	Role         Role
	Phone        *string
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id uint) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetAll(ctx context.Context) ([]User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uint) error
}

type UserService interface {
	CreateUser(ctx context.Context, user *User) error
	GetUserByID(ctx context.Context, id uint) (*User, error)
	GetAllUsers(ctx context.Context) ([]User, error)
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, id uint) error
	Login(ctx context.Context, email, password string) (string, *User, error)
}
