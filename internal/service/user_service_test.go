package service

import (
	"context"
	"errors"
	"strings"
	"testing"

	"digital-greenhouse/greenhouse-be/internal/domain"

	"golang.org/x/crypto/bcrypt"
)

// MockUserRepository implements domain.UserRepository
type MockUserRepository struct {
	CreateFunc     func(ctx context.Context, user *domain.User) error
	GetByIDFunc    func(ctx context.Context, id uint) (*domain.User, error)
	GetByEmailFunc func(ctx context.Context, email string) (*domain.User, error)
	GetAllFunc     func(ctx context.Context) ([]domain.User, error)
	UpdateFunc     func(ctx context.Context, user *domain.User) error
	DeleteFunc     func(ctx context.Context, id uint) error
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, user)
	}
	return nil
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uint) (*domain.User, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	if m.GetByEmailFunc != nil {
		return m.GetByEmailFunc(ctx, email)
	}
	return nil, nil
}

func (m *MockUserRepository) GetAll(ctx context.Context) ([]domain.User, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc(ctx)
	}
	return nil, nil
}

func (m *MockUserRepository) Update(ctx context.Context, user *domain.User) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, user)
	}
	return nil
}

func (m *MockUserRepository) Delete(ctx context.Context, id uint) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name        string
		user        *domain.User
		mockSetup   func(mock *MockUserRepository)
		expectedErr string
		checkRole   domain.Role
	}{
		{
			name:        "error - missing name",
			user:        &domain.User{Email: "test@example.com", PasswordHash: "password"},
			mockSetup:   func(mock *MockUserRepository) {},
			expectedErr: "name, email y password son requeridos",
		},
		{
			name:        "error - missing email",
			user:        &domain.User{Name: "Test User", PasswordHash: "password"},
			mockSetup:   func(mock *MockUserRepository) {},
			expectedErr: "name, email y password son requeridos",
		},
		{
			name:        "error - missing password",
			user:        &domain.User{Name: "Test User", Email: "test@example.com"},
			mockSetup:   func(mock *MockUserRepository) {},
			expectedErr: "name, email y password son requeridos",
		},
		{
			name: "error - duplicate email",
			user: &domain.User{Name: "Test User", Email: "test@example.com", PasswordHash: "password"},
			mockSetup: func(mock *MockUserRepository) {
				mock.GetByEmailFunc = func(ctx context.Context, email string) (*domain.User, error) {
					return &domain.User{ID: 1, Email: email}, nil
				}
			},
			expectedErr: "el email ya está registrado",
		},
		{
			name: "error - bcrypt password too long",
			user: &domain.User{
				Name:         "Test User",
				Email:        "test@example.com",
				PasswordHash: string(make([]byte, 73)), // exceeds bcrypt length limit (72 bytes)
			},
			mockSetup: func(mock *MockUserRepository) {
				mock.GetByEmailFunc = func(ctx context.Context, email string) (*domain.User, error) {
					return nil, nil // not found
				}
			},
			expectedErr: "error al hashear la contraseña",
		},
		{
			name: "error - repository create error",
			user: &domain.User{Name: "Test User", Email: "test@example.com", PasswordHash: "password"},
			mockSetup: func(mock *MockUserRepository) {
				mock.GetByEmailFunc = func(ctx context.Context, email string) (*domain.User, error) {
					return nil, nil
				}
				mock.CreateFunc = func(ctx context.Context, user *domain.User) error {
					return errors.New("db error")
				}
			},
			expectedErr: "db error",
		},
		{
			name: "success - default client role",
			user: &domain.User{Name: "Test User", Email: "test@example.com", PasswordHash: "password"},
			mockSetup: func(mock *MockUserRepository) {
				mock.GetByEmailFunc = func(ctx context.Context, email string) (*domain.User, error) {
					return nil, nil
				}
				mock.CreateFunc = func(ctx context.Context, user *domain.User) error {
					return nil
				}
			},
			expectedErr: "",
			checkRole:   domain.RoleClient,
		},
		{
			name: "success - explicit superadmin role",
			user: &domain.User{Name: "Test User", Email: "test@example.com", PasswordHash: "password", Role: domain.RoleSuperAdmin},
			mockSetup: func(mock *MockUserRepository) {
				mock.GetByEmailFunc = func(ctx context.Context, email string) (*domain.User, error) {
					return nil, nil
				}
				mock.CreateFunc = func(ctx context.Context, user *domain.User) error {
					return nil
				}
			},
			expectedErr: "",
			checkRole:   domain.RoleSuperAdmin,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockUserRepository{}
			tt.mockSetup(mockRepo)
			s := NewUserService(mockRepo)

			err := s.CreateUser(context.Background(), tt.user)
			if tt.expectedErr != "" {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tt.expectedErr)
				}
				if !strings.Contains(err.Error(), tt.expectedErr) {
					t.Errorf("expected error %q to contain %q", err.Error(), tt.expectedErr)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if tt.user.Role != tt.checkRole {
					t.Errorf("expected role %s, got %s", tt.checkRole, tt.user.Role)
				}
				if tt.user.PasswordHash == "password" {
					t.Errorf("password was not hashed")
				}
			}
		})
	}
}

func TestGetUserByID(t *testing.T) {
	tests := []struct {
		name         string
		id           uint
		mockSetup    func(mock *MockUserRepository)
		expectedUser *domain.User
		expectedErr  string
	}{
		{
			name: "error - repository error",
			id:   1,
			mockSetup: func(mock *MockUserRepository) {
				mock.GetByIDFunc = func(ctx context.Context, id uint) (*domain.User, error) {
					return nil, errors.New("db error")
				}
			},
			expectedUser: nil,
			expectedErr:  "db error",
		},
		{
			name: "success",
			id:   1,
			mockSetup: func(mock *MockUserRepository) {
				mock.GetByIDFunc = func(ctx context.Context, id uint) (*domain.User, error) {
					return &domain.User{ID: id, Name: "Test User"}, nil
				}
			},
			expectedUser: &domain.User{ID: 1, Name: "Test User"},
			expectedErr:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockUserRepository{}
			tt.mockSetup(mockRepo)
			s := NewUserService(mockRepo)

			user, err := s.GetUserByID(context.Background(), tt.id)
			if tt.expectedErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.expectedErr) {
					t.Fatalf("expected error containing %q, got %v", tt.expectedErr, err)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if user == nil || user.ID != tt.expectedUser.ID || user.Name != tt.expectedUser.Name {
					t.Errorf("expected user %+v, got %+v", tt.expectedUser, user)
				}
			}
		})
	}
}

func TestGetAllUsers(t *testing.T) {
	tests := []struct {
		name         string
		mockSetup    func(mock *MockUserRepository)
		expectedList []domain.User
		expectedErr  string
	}{
		{
			name: "error - repository error",
			mockSetup: func(mock *MockUserRepository) {
				mock.GetAllFunc = func(ctx context.Context) ([]domain.User, error) {
					return nil, errors.New("db error")
				}
			},
			expectedList: nil,
			expectedErr:  "db error",
		},
		{
			name: "success",
			mockSetup: func(mock *MockUserRepository) {
				mock.GetAllFunc = func(ctx context.Context) ([]domain.User, error) {
					return []domain.User{
						{ID: 1, Name: "User 1"},
						{ID: 2, Name: "User 2"},
					}, nil
				}
			},
			expectedList: []domain.User{
				{ID: 1, Name: "User 1"},
				{ID: 2, Name: "User 2"},
			},
			expectedErr:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockUserRepository{}
			tt.mockSetup(mockRepo)
			s := NewUserService(mockRepo)

			users, err := s.GetAllUsers(context.Background())
			if tt.expectedErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.expectedErr) {
					t.Fatalf("expected error containing %q, got %v", tt.expectedErr, err)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if len(users) != len(tt.expectedList) {
					t.Fatalf("expected slice of len %d, got %d", len(tt.expectedList), len(users))
				}
				for i := range users {
					if users[i].ID != tt.expectedList[i].ID || users[i].Name != tt.expectedList[i].Name {
						t.Errorf("expected user %+v at index %d, got %+v", tt.expectedList[i], i, users[i])
					}
				}
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	tests := []struct {
		name        string
		user        *domain.User
		mockSetup   func(mock *MockUserRepository)
		expectedErr string
	}{
		{
			name:        "error - ID is zero",
			user:        &domain.User{ID: 0, Name: "No ID"},
			mockSetup:   func(mock *MockUserRepository) {},
			expectedErr: "el ID del usuario es requerido para actualizar",
		},
		{
			name: "error - repository error",
			user: &domain.User{ID: 1, Name: "Test User"},
			mockSetup: func(mock *MockUserRepository) {
				mock.UpdateFunc = func(ctx context.Context, user *domain.User) error {
					return errors.New("db error")
				}
			},
			expectedErr: "db error",
		},
		{
			name: "success",
			user: &domain.User{ID: 1, Name: "Test User"},
			mockSetup: func(mock *MockUserRepository) {
				mock.UpdateFunc = func(ctx context.Context, user *domain.User) error {
					return nil
				}
			},
			expectedErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockUserRepository{}
			tt.mockSetup(mockRepo)
			s := NewUserService(mockRepo)

			err := s.UpdateUser(context.Background(), tt.user)
			if tt.expectedErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.expectedErr) {
					t.Fatalf("expected error containing %q, got %v", tt.expectedErr, err)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {
	tests := []struct {
		name        string
		id          uint
		mockSetup   func(mock *MockUserRepository)
		expectedErr string
	}{
		{
			name: "error - repository error",
			id:   1,
			mockSetup: func(mock *MockUserRepository) {
				mock.DeleteFunc = func(ctx context.Context, id uint) error {
					return errors.New("db error")
				}
			},
			expectedErr: "db error",
		},
		{
			name: "success",
			id:   1,
			mockSetup: func(mock *MockUserRepository) {
				mock.DeleteFunc = func(ctx context.Context, id uint) error {
					return nil
				}
			},
			expectedErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockUserRepository{}
			tt.mockSetup(mockRepo)
			s := NewUserService(mockRepo)

			err := s.DeleteUser(context.Background(), tt.id)
			if tt.expectedErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.expectedErr) {
					t.Fatalf("expected error containing %q, got %v", tt.expectedErr, err)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
			}
		})
	}
}

func TestLogin(t *testing.T) {
	password := "correct_password"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash password for test: %v", err)
	}

	tests := []struct {
		name          string
		email         string
		password      string
		mockSetup     func(mock *MockUserRepository)
		expectedToken bool
		expectedUser  *domain.User
		expectedErr   string
	}{
		{
			name:     "error - user not found/db error in GetByEmail",
			email:    "notfound@example.com",
			password: "password",
			mockSetup: func(mock *MockUserRepository) {
				mock.GetByEmailFunc = func(ctx context.Context, email string) (*domain.User, error) {
					return nil, errors.New("user not found")
				}
			},
			expectedToken: false,
			expectedUser:  nil,
			expectedErr:   "credenciales inválidas",
		},
		{
			name:     "error - bcrypt compare failure (wrong password)",
			email:    "test@example.com",
			password: "wrong_password",
			mockSetup: func(mock *MockUserRepository) {
				mock.GetByEmailFunc = func(ctx context.Context, email string) (*domain.User, error) {
					return &domain.User{
						ID:           1,
						Email:        email,
						PasswordHash: string(hash),
						Name:         "Test User",
						Role:         domain.RoleClient,
					}, nil
				}
			},
			expectedToken: false,
			expectedUser:  nil,
			expectedErr:   "credenciales inválidas",
		},
		{
			name:     "success",
			email:    "test@example.com",
			password: password,
			mockSetup: func(mock *MockUserRepository) {
				mock.GetByEmailFunc = func(ctx context.Context, email string) (*domain.User, error) {
					return &domain.User{
						ID:           1,
						Email:        email,
						PasswordHash: string(hash),
						Name:         "Test User",
						Role:         domain.RoleClient,
					}, nil
				}
			},
			expectedToken: true,
			expectedUser: &domain.User{
				ID:    1,
				Email: "test@example.com",
				Name:  "Test User",
				Role:  domain.RoleClient,
			},
			expectedErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockUserRepository{}
			tt.mockSetup(mockRepo)
			s := NewUserService(mockRepo)

			token, user, err := s.Login(context.Background(), tt.email, tt.password)
			if tt.expectedErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.expectedErr) {
					t.Fatalf("expected error containing %q, got %v", tt.expectedErr, err)
				}
				if token != "" {
					t.Errorf("expected empty token, got %q", token)
				}
				if user != nil {
					t.Errorf("expected nil user, got %+v", user)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if tt.expectedToken && token == "" {
					t.Errorf("expected a token, got empty string")
				}
				if user == nil || user.ID != tt.expectedUser.ID || user.Email != tt.expectedUser.Email || user.Name != tt.expectedUser.Name || user.Role != tt.expectedUser.Role {
					t.Errorf("expected user %+v, got %+v", tt.expectedUser, user)
				}
			}
		})
	}
}
