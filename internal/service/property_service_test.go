package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"digital-greenhouse/greenhouse-be/internal/domain"
)

type propServiceMockPropertyRepo struct {
	CreateFunc       func(ctx context.Context, property *domain.Property) error
	GetAllFunc       func(ctx context.Context, filter domain.PropertyFilter) ([]domain.Property, error)
	GetByIDFunc      func(ctx context.Context, id uint) (*domain.Property, error)
	GetByOwnerIDFunc func(ctx context.Context, ownerID uint) ([]domain.Property, error)
	UpdateFunc       func(ctx context.Context, property *domain.Property) error
	DeleteFunc       func(ctx context.Context, id uint) error
	AddImageFunc     func(ctx context.Context, image *domain.PropertyImage) error
	GetImageByIDFunc func(ctx context.Context, id uint) (*domain.PropertyImage, error)
	UpdateImageFunc  func(ctx context.Context, image *domain.PropertyImage) error
	DeleteImageFunc  func(ctx context.Context, id uint) error
}

func (m *propServiceMockPropertyRepo) Create(ctx context.Context, property *domain.Property) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, property)
	}
	return nil
}

func (m *propServiceMockPropertyRepo) GetAll(ctx context.Context, filter domain.PropertyFilter) ([]domain.Property, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc(ctx, filter)
	}
	return nil, nil
}

func (m *propServiceMockPropertyRepo) GetByID(ctx context.Context, id uint) (*domain.Property, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *propServiceMockPropertyRepo) GetByOwnerID(ctx context.Context, ownerID uint) ([]domain.Property, error) {
	if m.GetByOwnerIDFunc != nil {
		return m.GetByOwnerIDFunc(ctx, ownerID)
	}
	return nil, nil
}

func (m *propServiceMockPropertyRepo) Update(ctx context.Context, property *domain.Property) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, property)
	}
	return nil
}

func (m *propServiceMockPropertyRepo) Delete(ctx context.Context, id uint) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func (m *propServiceMockPropertyRepo) AddImage(ctx context.Context, image *domain.PropertyImage) error {
	if m.AddImageFunc != nil {
		return m.AddImageFunc(ctx, image)
	}
	return nil
}

func (m *propServiceMockPropertyRepo) GetImageByID(ctx context.Context, id uint) (*domain.PropertyImage, error) {
	if m.GetImageByIDFunc != nil {
		return m.GetImageByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *propServiceMockPropertyRepo) UpdateImage(ctx context.Context, image *domain.PropertyImage) error {
	if m.UpdateImageFunc != nil {
		return m.UpdateImageFunc(ctx, image)
	}
	return nil
}

func (m *propServiceMockPropertyRepo) DeleteImage(ctx context.Context, id uint) error {
	if m.DeleteImageFunc != nil {
		return m.DeleteImageFunc(ctx, id)
	}
	return nil
}

type propServiceMockUserRepo struct {
	CreateFunc     func(ctx context.Context, user *domain.User) error
	GetByIDFunc    func(ctx context.Context, id uint) (*domain.User, error)
	GetByEmailFunc func(ctx context.Context, email string) (*domain.User, error)
	GetAllFunc     func(ctx context.Context) ([]domain.User, error)
	UpdateFunc     func(ctx context.Context, user *domain.User) error
	DeleteFunc     func(ctx context.Context, id uint) error
}

func (m *propServiceMockUserRepo) Create(ctx context.Context, user *domain.User) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, user)
	}
	return nil
}

func (m *propServiceMockUserRepo) GetByID(ctx context.Context, id uint) (*domain.User, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *propServiceMockUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	if m.GetByEmailFunc != nil {
		return m.GetByEmailFunc(ctx, email)
	}
	return nil, nil
}

func (m *propServiceMockUserRepo) GetAll(ctx context.Context) ([]domain.User, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc(ctx)
	}
	return nil, nil
}

func (m *propServiceMockUserRepo) Update(ctx context.Context, user *domain.User) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, user)
	}
	return nil
}

func (m *propServiceMockUserRepo) Delete(ctx context.Context, id uint) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

type propServiceMockBookingRepo struct {
	CreateQuoteFunc                    func(ctx context.Context, quote *domain.Quote) error
	GetQuoteByIDFunc                   func(ctx context.Context, id uint) (*domain.Quote, error)
	GetQuotesByClientIDFunc            func(ctx context.Context, clientID uint) ([]domain.Quote, error)
	UpdateQuoteStatusFunc              func(ctx context.Context, id uint, status domain.QuoteStatus) error
	CreateBookingFunc                  func(ctx context.Context, booking *domain.Booking) error
	GetBookingByIDFunc                 func(ctx context.Context, id uint) (*domain.Booking, error)
	GetBookingsByClientIDFunc          func(ctx context.Context, clientID uint) ([]domain.Booking, error)
	GetBookingsByPropertyIDFunc        func(ctx context.Context, propertyID uint) ([]domain.Booking, error)
	GetBookingsByOwnerIDFunc           func(ctx context.Context, ownerID uint) ([]domain.Booking, error)
	GetReservedDatesByPropertyIDFunc   func(ctx context.Context, propertyID uint) ([]domain.Booking, error)
	UpdateBookingStatusFunc            func(ctx context.Context, id uint, status domain.BookingStatus, reason string) error
	CheckAvailabilityFunc              func(ctx context.Context, propertyID uint, checkIn, checkOut time.Time) (bool, error)
	CreatePricingRuleFunc              func(ctx context.Context, rule *domain.PricingRule) error
	GetPricingRulesByPropertyIDFunc    func(ctx context.Context, propertyID uint, start, end time.Time) ([]domain.PricingRule, error)
	GetAllPricingRulesByPropertyIDFunc func(ctx context.Context, propertyID uint) ([]domain.PricingRule, error)
	DeletePricingRuleFunc              func(ctx context.Context, id uint) error
}

func (m *propServiceMockBookingRepo) CreateQuote(ctx context.Context, quote *domain.Quote) error {
	if m.CreateQuoteFunc != nil {
		return m.CreateQuoteFunc(ctx, quote)
	}
	return nil
}

func (m *propServiceMockBookingRepo) GetQuoteByID(ctx context.Context, id uint) (*domain.Quote, error) {
	if m.GetQuoteByIDFunc != nil {
		return m.GetQuoteByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *propServiceMockBookingRepo) GetQuotesByClientID(ctx context.Context, clientID uint) ([]domain.Quote, error) {
	if m.GetQuotesByClientIDFunc != nil {
		return m.GetQuotesByClientIDFunc(ctx, clientID)
	}
	return nil, nil
}

func (m *propServiceMockBookingRepo) UpdateQuoteStatus(ctx context.Context, id uint, status domain.QuoteStatus) error {
	if m.UpdateQuoteStatusFunc != nil {
		return m.UpdateQuoteStatusFunc(ctx, id, status)
	}
	return nil
}

func (m *propServiceMockBookingRepo) CreateBooking(ctx context.Context, booking *domain.Booking) error {
	if m.CreateBookingFunc != nil {
		return m.CreateBookingFunc(ctx, booking)
	}
	return nil
}

func (m *propServiceMockBookingRepo) GetBookingByID(ctx context.Context, id uint) (*domain.Booking, error) {
	if m.GetBookingByIDFunc != nil {
		return m.GetBookingByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *propServiceMockBookingRepo) GetBookingsByClientID(ctx context.Context, clientID uint) ([]domain.Booking, error) {
	if m.GetBookingsByClientIDFunc != nil {
		return m.GetBookingsByClientIDFunc(ctx, clientID)
	}
	return nil, nil
}

func (m *propServiceMockBookingRepo) GetBookingsByPropertyID(ctx context.Context, propertyID uint) ([]domain.Booking, error) {
	if m.GetBookingsByPropertyIDFunc != nil {
		return m.GetBookingsByPropertyIDFunc(ctx, propertyID)
	}
	return nil, nil
}

func (m *propServiceMockBookingRepo) GetBookingsByOwnerID(ctx context.Context, ownerID uint) ([]domain.Booking, error) {
	if m.GetBookingsByOwnerIDFunc != nil {
		return m.GetBookingsByOwnerIDFunc(ctx, ownerID)
	}
	return nil, nil
}

func (m *propServiceMockBookingRepo) GetReservedDatesByPropertyID(ctx context.Context, propertyID uint) ([]domain.Booking, error) {
	if m.GetReservedDatesByPropertyIDFunc != nil {
		return m.GetReservedDatesByPropertyIDFunc(ctx, propertyID)
	}
	return nil, nil
}

func (m *propServiceMockBookingRepo) UpdateBookingStatus(ctx context.Context, id uint, status domain.BookingStatus, reason string) error {
	if m.UpdateBookingStatusFunc != nil {
		return m.UpdateBookingStatusFunc(ctx, id, status, reason)
	}
	return nil
}

func (m *propServiceMockBookingRepo) CheckAvailability(ctx context.Context, propertyID uint, checkIn, checkOut time.Time) (bool, error) {
	if m.CheckAvailabilityFunc != nil {
		return m.CheckAvailabilityFunc(ctx, propertyID, checkIn, checkOut)
	}
	return false, nil
}

func (m *propServiceMockBookingRepo) CreatePricingRule(ctx context.Context, rule *domain.PricingRule) error {
	if m.CreatePricingRuleFunc != nil {
		return m.CreatePricingRuleFunc(ctx, rule)
	}
	return nil
}

func (m *propServiceMockBookingRepo) GetPricingRulesByPropertyID(ctx context.Context, propertyID uint, start, end time.Time) ([]domain.PricingRule, error) {
	if m.GetPricingRulesByPropertyIDFunc != nil {
		return m.GetPricingRulesByPropertyIDFunc(ctx, propertyID, start, end)
	}
	return nil, nil
}

func (m *propServiceMockBookingRepo) GetAllPricingRulesByPropertyID(ctx context.Context, propertyID uint) ([]domain.PricingRule, error) {
	if m.GetAllPricingRulesByPropertyIDFunc != nil {
		return m.GetAllPricingRulesByPropertyIDFunc(ctx, propertyID)
	}
	return nil, nil
}

func (m *propServiceMockBookingRepo) DeletePricingRule(ctx context.Context, id uint) error {
	if m.DeletePricingRuleFunc != nil {
		return m.DeletePricingRuleFunc(ctx, id)
	}
	return nil
}

func TestCreateProperty(t *testing.T) {
	tests := []struct {
		name        string
		property    *domain.Property
		mockProp    func(m *propServiceMockPropertyRepo)
		mockUser    func(m *propServiceMockUserRepo)
		expectedErr string
		verify      func(t *testing.T, p *domain.Property)
	}{
		{
			name: "error - missing name",
			property: &domain.Property{
				Name:    "",
				OwnerID: 1,
			},
			expectedErr: "nombre y dueño son requeridos",
		},
		{
			name: "error - missing owner ID",
			property: &domain.Property{
				Name:    "Beautiful Cabin",
				OwnerID: 0,
			},
			expectedErr: "nombre y dueño son requeridos",
		},
		{
			name: "error - property repo Create fails",
			property: &domain.Property{
				Name:    "Beautiful Cabin",
				OwnerID: 1,
			},
			mockProp: func(m *propServiceMockPropertyRepo) {
				m.CreateFunc = func(ctx context.Context, p *domain.Property) error {
					return errors.New("db create error")
				}
			},
			expectedErr: "db create error",
		},
		{
			name: "error - user repo GetByID fails",
			property: &domain.Property{
				Name:    "Beautiful Cabin",
				OwnerID: 1,
			},
			mockProp: func(m *propServiceMockPropertyRepo) {
				m.CreateFunc = func(ctx context.Context, p *domain.Property) error {
					p.ID = 100
					return nil
				}
			},
			mockUser: func(m *propServiceMockUserRepo) {
				m.GetByIDFunc = func(ctx context.Context, id uint) (*domain.User, error) {
					return nil, errors.New("user not found")
				}
			},
			expectedErr: "user not found",
		},
		{
			name: "success - owner role is already owner",
			property: &domain.Property{
				Name:    "Beautiful Cabin",
				OwnerID: 1,
			},
			mockProp: func(m *propServiceMockPropertyRepo) {
				m.CreateFunc = func(ctx context.Context, p *domain.Property) error {
					p.ID = 100
					return nil
				}
			},
			mockUser: func(m *propServiceMockUserRepo) {
				m.GetByIDFunc = func(ctx context.Context, id uint) (*domain.User, error) {
					return &domain.User{ID: id, Role: domain.RoleOwner}, nil
				}
			},
			verify: func(t *testing.T, p *domain.Property) {
				if p.Status != domain.PropertyActive {
					t.Errorf("expected property status ACTIVE, got %v", p.Status)
				}
			},
		},
		{
			name: "error - user role update fails",
			property: &domain.Property{
				Name:    "Beautiful Cabin",
				OwnerID: 1,
			},
			mockProp: func(m *propServiceMockPropertyRepo) {
				m.CreateFunc = func(ctx context.Context, p *domain.Property) error {
					return nil
				}
			},
			mockUser: func(m *propServiceMockUserRepo) {
				m.GetByIDFunc = func(ctx context.Context, id uint) (*domain.User, error) {
					return &domain.User{ID: id, Role: domain.RoleClient}, nil
				}
				m.UpdateFunc = func(ctx context.Context, u *domain.User) error {
					return errors.New("db update user error")
				}
			},
			expectedErr: "db update user error",
		},
		{
			name: "success - user role update from client to owner succeeds",
			property: &domain.Property{
				Name:    "Beautiful Cabin",
				OwnerID: 1,
				Status:  domain.PropertyMaintenance,
			},
			mockProp: func(m *propServiceMockPropertyRepo) {
				m.CreateFunc = func(ctx context.Context, p *domain.Property) error {
					return nil
				}
			},
			mockUser: func(m *propServiceMockUserRepo) {
				m.GetByIDFunc = func(ctx context.Context, id uint) (*domain.User, error) {
					return &domain.User{ID: id, Role: domain.RoleClient}, nil
				}
				m.UpdateFunc = func(ctx context.Context, u *domain.User) error {
					if u.Role != domain.RoleOwner {
						t.Errorf("expected updated role to be OWNER, got %s", u.Role)
					}
					return nil
				}
			},
			verify: func(t *testing.T, p *domain.Property) {
				if p.Status != domain.PropertyMaintenance {
					t.Errorf("expected property status MAINTENANCE to be preserved, got %v", p.Status)
				}
			},
		},
		{
			name: "error - adding initial image fails",
			property: &domain.Property{
				Name:    "Beautiful Cabin",
				OwnerID: 1,
				Images: []domain.PropertyImage{
					{ImageData: "base64data"},
				},
			},
			mockProp: func(m *propServiceMockPropertyRepo) {
				m.CreateFunc = func(ctx context.Context, p *domain.Property) error {
					p.ID = 100
					return nil
				}
				m.AddImageFunc = func(ctx context.Context, img *domain.PropertyImage) error {
					return errors.New("image save error")
				}
			},
			mockUser: func(m *propServiceMockUserRepo) {
				m.GetByIDFunc = func(ctx context.Context, id uint) (*domain.User, error) {
					return &domain.User{ID: id, Role: domain.RoleOwner}, nil
				}
			},
			expectedErr: "image save error",
		},
		{
			name: "success - adding initial images succeeds",
			property: &domain.Property{
				Name:    "Beautiful Cabin",
				OwnerID: 1,
				Images: []domain.PropertyImage{
					{ImageData: "base64data1"},
					{ImageData: "base64data2"},
				},
			},
			mockProp: func(m *propServiceMockPropertyRepo) {
				m.CreateFunc = func(ctx context.Context, p *domain.Property) error {
					p.ID = 100
					return nil
				}
				m.AddImageFunc = func(ctx context.Context, img *domain.PropertyImage) error {
					if img.PropertyID != 100 {
						t.Errorf("expected PropertyID 100, got %d", img.PropertyID)
					}
					return nil
				}
			},
			mockUser: func(m *propServiceMockUserRepo) {
				m.GetByIDFunc = func(ctx context.Context, id uint) (*domain.User, error) {
					return &domain.User{ID: id, Role: domain.RoleOwner}, nil
				}
			},
			verify: func(t *testing.T, p *domain.Property) {
				for _, img := range p.Images {
					if img.PropertyID != 100 {
						t.Errorf("expected images PropertyID updated to 100, got %d", img.PropertyID)
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mProp := &propServiceMockPropertyRepo{}
			if tt.mockProp != nil {
				tt.mockProp(mProp)
			}
			mUser := &propServiceMockUserRepo{}
			if tt.mockUser != nil {
				tt.mockUser(mUser)
			}
			mBooking := &propServiceMockBookingRepo{}

			s := NewPropertyService(mProp, mUser, mBooking)
			err := s.CreateProperty(context.Background(), tt.property)

			if tt.expectedErr != "" {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tt.expectedErr)
				}
				if err.Error() != tt.expectedErr {
					t.Errorf("expected error message %q, got %q", tt.expectedErr, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if tt.verify != nil {
					tt.verify(t, tt.property)
				}
			}
		})
	}
}

func TestListProperties(t *testing.T) {
	tests := []struct {
		name        string
		filter      domain.PropertyFilter
		mockSetup   func(m *propServiceMockPropertyRepo)
		expectedLen int
		expectedErr string
	}{
		{
			name:   "success",
			filter: domain.PropertyFilter{Search: "cabin"},
			mockSetup: func(m *propServiceMockPropertyRepo) {
				m.GetAllFunc = func(ctx context.Context, filter domain.PropertyFilter) ([]domain.Property, error) {
					if filter.Search != "cabin" {
						t.Errorf("expected search filter 'cabin', got %q", filter.Search)
					}
					return []domain.Property{
						{ID: 1, Name: "Cabin 1"},
						{ID: 2, Name: "Cabin 2"},
					}, nil
				}
			},
			expectedLen: 2,
		},
		{
			name: "error",
			mockSetup: func(m *propServiceMockPropertyRepo) {
				m.GetAllFunc = func(ctx context.Context, filter domain.PropertyFilter) ([]domain.Property, error) {
					return nil, errors.New("db error")
				}
			},
			expectedErr: "db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mProp := &propServiceMockPropertyRepo{}
			tt.mockSetup(mProp)
			s := NewPropertyService(mProp, &propServiceMockUserRepo{}, &propServiceMockBookingRepo{})
			res, err := s.ListProperties(context.Background(), tt.filter)
			if tt.expectedErr != "" {
				if err == nil || err.Error() != tt.expectedErr {
					t.Fatalf("expected error %q, got %v", tt.expectedErr, err)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if len(res) != tt.expectedLen {
					t.Errorf("expected len %d, got %d", tt.expectedLen, len(res))
				}
			}
		})
	}
}

func TestGetPropertyByID(t *testing.T) {
	tests := []struct {
		name        string
		id          uint
		mockSetup   func(m *propServiceMockPropertyRepo)
		expectedErr string
	}{
		{
			name: "success",
			id:   1,
			mockSetup: func(m *propServiceMockPropertyRepo) {
				m.GetByIDFunc = func(ctx context.Context, id uint) (*domain.Property, error) {
					return &domain.Property{ID: id, Name: "Cabin"}, nil
				}
			},
		},
		{
			name: "error",
			id:   1,
			mockSetup: func(m *propServiceMockPropertyRepo) {
				m.GetByIDFunc = func(ctx context.Context, id uint) (*domain.Property, error) {
					return nil, errors.New("not found")
				}
			},
			expectedErr: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mProp := &propServiceMockPropertyRepo{}
			tt.mockSetup(mProp)
			s := NewPropertyService(mProp, &propServiceMockUserRepo{}, &propServiceMockBookingRepo{})
			res, err := s.GetPropertyByID(context.Background(), tt.id)
			if tt.expectedErr != "" {
				if err == nil || err.Error() != tt.expectedErr {
					t.Fatalf("expected error %q, got %v", tt.expectedErr, err)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if res.ID != tt.id {
					t.Errorf("expected ID %d, got %d", tt.id, res.ID)
				}
			}
		})
	}
}

func TestGetPropertiesByOwner(t *testing.T) {
	tests := []struct {
		name        string
		ownerID     uint
		mockSetup   func(m *propServiceMockPropertyRepo)
		expectedErr string
	}{
		{
			name:    "success",
			ownerID: 5,
			mockSetup: func(m *propServiceMockPropertyRepo) {
				m.GetByOwnerIDFunc = func(ctx context.Context, ownerID uint) ([]domain.Property, error) {
					return []domain.Property{{ID: 1, OwnerID: ownerID}}, nil
				}
			},
		},
		{
			name:    "error",
			ownerID: 5,
			mockSetup: func(m *propServiceMockPropertyRepo) {
				m.GetByOwnerIDFunc = func(ctx context.Context, ownerID uint) ([]domain.Property, error) {
					return nil, errors.New("db error")
				}
			},
			expectedErr: "db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mProp := &propServiceMockPropertyRepo{}
			tt.mockSetup(mProp)
			s := NewPropertyService(mProp, &propServiceMockUserRepo{}, &propServiceMockBookingRepo{})
			res, err := s.GetPropertiesByOwner(context.Background(), tt.ownerID)
			if tt.expectedErr != "" {
				if err == nil || err.Error() != tt.expectedErr {
					t.Fatalf("expected error %q, got %v", tt.expectedErr, err)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if len(res) != 1 || res[0].OwnerID != tt.ownerID {
					t.Errorf("expected owner ID %d, got %+v", tt.ownerID, res)
				}
			}
		})
	}
}

func TestUpdateProperty(t *testing.T) {
	tests := []struct {
		name        string
		property    *domain.Property
		mockSetup   func(m *propServiceMockPropertyRepo)
		expectedErr string
	}{
		{
			name: "error - missing ID",
			property: &domain.Property{
				ID:   0,
				Name: "Cabin",
			},
			expectedErr: "ID de propiedad requerido",
		},
		{
			name: "error - db update fails",
			property: &domain.Property{
				ID:   1,
				Name: "Cabin",
			},
			mockSetup: func(m *propServiceMockPropertyRepo) {
				m.UpdateFunc = func(ctx context.Context, p *domain.Property) error {
					return errors.New("db update error")
				}
			},
			expectedErr: "db update error",
		},
		{
			name: "success",
			property: &domain.Property{
				ID:   1,
				Name: "Cabin",
			},
			mockSetup: func(m *propServiceMockPropertyRepo) {
				m.UpdateFunc = func(ctx context.Context, p *domain.Property) error {
					return nil
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mProp := &propServiceMockPropertyRepo{}
			if tt.mockSetup != nil {
				tt.mockSetup(mProp)
			}
			s := NewPropertyService(mProp, &propServiceMockUserRepo{}, &propServiceMockBookingRepo{})
			err := s.UpdateProperty(context.Background(), tt.property)
			if tt.expectedErr != "" {
				if err == nil || err.Error() != tt.expectedErr {
					t.Fatalf("expected error %q, got %v", tt.expectedErr, err)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
			}
		})
	}
}

func TestDeleteProperty(t *testing.T) {
	tests := []struct {
		name        string
		id          uint
		mockSetup   func(m *propServiceMockPropertyRepo)
		expectedErr string
	}{
		{
			name: "success",
			id:   1,
			mockSetup: func(m *propServiceMockPropertyRepo) {
				m.DeleteFunc = func(ctx context.Context, id uint) error {
					return nil
				}
			},
		},
		{
			name: "error",
			id:   1,
			mockSetup: func(m *propServiceMockPropertyRepo) {
				m.DeleteFunc = func(ctx context.Context, id uint) error {
					return errors.New("db delete error")
				}
			},
			expectedErr: "db delete error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mProp := &propServiceMockPropertyRepo{}
			tt.mockSetup(mProp)
			s := NewPropertyService(mProp, &propServiceMockUserRepo{}, &propServiceMockBookingRepo{})
			err := s.DeleteProperty(context.Background(), tt.id)
			if tt.expectedErr != "" {
				if err == nil || err.Error() != tt.expectedErr {
					t.Fatalf("expected error %q, got %v", tt.expectedErr, err)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
			}
		})
	}
}

func TestAddPropertyImage(t *testing.T) {
	tests := []struct {
		name        string
		image       *domain.PropertyImage
		mockSetup   func(m *propServiceMockPropertyRepo)
		expectedErr string
	}{
		{
			name: "error - missing property ID",
			image: &domain.PropertyImage{
				PropertyID: 0,
				ImageData:  "base64",
			},
			expectedErr: "ID de propiedad y datos de imagen requeridos",
		},
		{
			name: "error - missing image data",
			image: &domain.PropertyImage{
				PropertyID: 1,
				ImageData:  "",
			},
			expectedErr: "ID de propiedad y datos de imagen requeridos",
		},
		{
			name: "error - db fails",
			image: &domain.PropertyImage{
				PropertyID: 1,
				ImageData:  "base64",
			},
			mockSetup: func(m *propServiceMockPropertyRepo) {
				m.AddImageFunc = func(ctx context.Context, image *domain.PropertyImage) error {
					return errors.New("db error")
				}
			},
			expectedErr: "db error",
		},
		{
			name: "success",
			image: &domain.PropertyImage{
				PropertyID: 1,
				ImageData:  "base64",
			},
			mockSetup: func(m *propServiceMockPropertyRepo) {
				m.AddImageFunc = func(ctx context.Context, image *domain.PropertyImage) error {
					return nil
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mProp := &propServiceMockPropertyRepo{}
			if tt.mockSetup != nil {
				tt.mockSetup(mProp)
			}
			s := NewPropertyService(mProp, &propServiceMockUserRepo{}, &propServiceMockBookingRepo{})
			err := s.AddPropertyImage(context.Background(), tt.image)
			if tt.expectedErr != "" {
				if err == nil || err.Error() != tt.expectedErr {
					t.Fatalf("expected error %q, got %v", tt.expectedErr, err)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
			}
		})
	}
}

func TestUpdatePropertyImage(t *testing.T) {
	tests := []struct {
		name        string
		image       *domain.PropertyImage
		mockSetup   func(m *propServiceMockPropertyRepo)
		expectedErr string
	}{
		{
			name: "error - missing image ID",
			image: &domain.PropertyImage{
				ID: 0,
			},
			expectedErr: "ID de imagen requerido",
		},
		{
			name: "error - db fails",
			image: &domain.PropertyImage{
				ID: 1,
			},
			mockSetup: func(m *propServiceMockPropertyRepo) {
				m.UpdateImageFunc = func(ctx context.Context, image *domain.PropertyImage) error {
					return errors.New("db error")
				}
			},
			expectedErr: "db error",
		},
		{
			name: "success",
			image: &domain.PropertyImage{
				ID: 1,
			},
			mockSetup: func(m *propServiceMockPropertyRepo) {
				m.UpdateImageFunc = func(ctx context.Context, image *domain.PropertyImage) error {
					return nil
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mProp := &propServiceMockPropertyRepo{}
			if tt.mockSetup != nil {
				tt.mockSetup(mProp)
			}
			s := NewPropertyService(mProp, &propServiceMockUserRepo{}, &propServiceMockBookingRepo{})
			err := s.UpdatePropertyImage(context.Background(), tt.image)
			if tt.expectedErr != "" {
				if err == nil || err.Error() != tt.expectedErr {
					t.Fatalf("expected error %q, got %v", tt.expectedErr, err)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
			}
		})
	}
}

func TestDeletePropertyImage(t *testing.T) {
	tests := []struct {
		name        string
		imageID     uint
		mockSetup   func(m *propServiceMockPropertyRepo)
		expectedErr string
	}{
		{
			name:    "success",
			imageID: 1,
			mockSetup: func(m *propServiceMockPropertyRepo) {
				m.DeleteImageFunc = func(ctx context.Context, id uint) error {
					return nil
				}
			},
		},
		{
			name:    "error",
			imageID: 1,
			mockSetup: func(m *propServiceMockPropertyRepo) {
				m.DeleteImageFunc = func(ctx context.Context, id uint) error {
					return errors.New("db error")
				}
			},
			expectedErr: "db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mProp := &propServiceMockPropertyRepo{}
			tt.mockSetup(mProp)
			s := NewPropertyService(mProp, &propServiceMockUserRepo{}, &propServiceMockBookingRepo{})
			err := s.DeletePropertyImage(context.Background(), tt.imageID)
			if tt.expectedErr != "" {
				if err == nil || err.Error() != tt.expectedErr {
					t.Fatalf("expected error %q, got %v", tt.expectedErr, err)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
			}
		})
	}
}

func TestCreatePricingRule(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name        string
		rule        *domain.PricingRule
		mockSetup   func(m *propServiceMockBookingRepo)
		expectedErr string
	}{
		{
			name: "error - missing property ID",
			rule: &domain.PricingRule{
				PropertyID: 0,
				StartDate:  now,
				EndDate:    now.Add(time.Hour),
			},
			expectedErr: "propiedad, fecha de inicio y fin son requeridas",
		},
		{
			name: "error - zero StartDate",
			rule: &domain.PricingRule{
				PropertyID: 1,
				StartDate:  time.Time{},
				EndDate:    now.Add(time.Hour),
			},
			expectedErr: "propiedad, fecha de inicio y fin son requeridas",
		},
		{
			name: "error - zero EndDate",
			rule: &domain.PricingRule{
				PropertyID: 1,
				StartDate:  now,
				EndDate:    time.Time{},
			},
			expectedErr: "propiedad, fecha de inicio y fin son requeridas",
		},
		{
			name: "error - price modifier <= 0",
			rule: &domain.PricingRule{
				PropertyID:    1,
				StartDate:     now,
				EndDate:       now.Add(time.Hour),
				PriceModifier: 0,
			},
			expectedErr: "el modificador de precio debe ser mayor a 0",
		},
		{
			name: "error - db error",
			rule: &domain.PricingRule{
				PropertyID:    1,
				StartDate:     now,
				EndDate:       now.Add(time.Hour),
				PriceModifier: 1.2,
			},
			mockSetup: func(m *propServiceMockBookingRepo) {
				m.CreatePricingRuleFunc = func(ctx context.Context, rule *domain.PricingRule) error {
					return errors.New("db error")
				}
			},
			expectedErr: "db error",
		},
		{
			name: "success",
			rule: &domain.PricingRule{
				PropertyID:    1,
				StartDate:     now,
				EndDate:       now.Add(time.Hour),
				PriceModifier: 1.2,
			},
			mockSetup: func(m *propServiceMockBookingRepo) {
				m.CreatePricingRuleFunc = func(ctx context.Context, rule *domain.PricingRule) error {
					return nil
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mBooking := &propServiceMockBookingRepo{}
			if tt.mockSetup != nil {
				tt.mockSetup(mBooking)
			}
			s := NewPropertyService(&propServiceMockPropertyRepo{}, &propServiceMockUserRepo{}, mBooking)
			err := s.CreatePricingRule(context.Background(), tt.rule)
			if tt.expectedErr != "" {
				if err == nil || err.Error() != tt.expectedErr {
					t.Fatalf("expected error %q, got %v", tt.expectedErr, err)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
			}
		})
	}
}

func TestListPricingRulesByProperty(t *testing.T) {
	tests := []struct {
		name        string
		propertyID  uint
		mockSetup   func(m *propServiceMockBookingRepo)
		expectedErr string
	}{
		{
			name:       "success",
			propertyID: 1,
			mockSetup: func(m *propServiceMockBookingRepo) {
				m.GetAllPricingRulesByPropertyIDFunc = func(ctx context.Context, propertyID uint) ([]domain.PricingRule, error) {
					return []domain.PricingRule{{ID: 1, PropertyID: propertyID}}, nil
				}
			},
		},
		{
			name:       "error",
			propertyID: 1,
			mockSetup: func(m *propServiceMockBookingRepo) {
				m.GetAllPricingRulesByPropertyIDFunc = func(ctx context.Context, propertyID uint) ([]domain.PricingRule, error) {
					return nil, errors.New("db error")
				}
			},
			expectedErr: "db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mBooking := &propServiceMockBookingRepo{}
			tt.mockSetup(mBooking)
			s := NewPropertyService(&propServiceMockPropertyRepo{}, &propServiceMockUserRepo{}, mBooking)
			res, err := s.ListPricingRulesByProperty(context.Background(), tt.propertyID)
			if tt.expectedErr != "" {
				if err == nil || err.Error() != tt.expectedErr {
					t.Fatalf("expected error %q, got %v", tt.expectedErr, err)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if len(res) != 1 || res[0].PropertyID != tt.propertyID {
					t.Errorf("expected property ID %d, got %+v", tt.propertyID, res)
				}
			}
		})
	}
}

func TestDeletePricingRule(t *testing.T) {
	tests := []struct {
		name        string
		id          uint
		mockSetup   func(m *propServiceMockBookingRepo)
		expectedErr string
	}{
		{
			name: "success",
			id:   1,
			mockSetup: func(m *propServiceMockBookingRepo) {
				m.DeletePricingRuleFunc = func(ctx context.Context, id uint) error {
					return nil
				}
			},
		},
		{
			name: "error",
			id:   1,
			mockSetup: func(m *propServiceMockBookingRepo) {
				m.DeletePricingRuleFunc = func(ctx context.Context, id uint) error {
					return errors.New("db error")
				}
			},
			expectedErr: "db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mBooking := &propServiceMockBookingRepo{}
			tt.mockSetup(mBooking)
			s := NewPropertyService(&propServiceMockPropertyRepo{}, &propServiceMockUserRepo{}, mBooking)
			err := s.DeletePricingRule(context.Background(), tt.id)
			if tt.expectedErr != "" {
				if err == nil || err.Error() != tt.expectedErr {
					t.Fatalf("expected error %q, got %v", tt.expectedErr, err)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
			}
		})
	}
}

func TestAutoGenerateHighSeasonRules(t *testing.T) {
	tests := []struct {
		name        string
		propertyID  uint
		mockSetup   func(m *propServiceMockBookingRepo)
		expectedErr string
	}{
		{
			name:       "success",
			propertyID: 10,
			mockSetup: func(m *propServiceMockBookingRepo) {
				var count int
				m.CreatePricingRuleFunc = func(ctx context.Context, rule *domain.PricingRule) error {
					count++
					if rule.PropertyID != 10 {
						t.Errorf("expected propertyID 10, got %d", rule.PropertyID)
					}
					if rule.PriceModifier != 1.10 {
						t.Errorf("expected price modifier 1.10, got %f", rule.PriceModifier)
					}
					return nil
				}
			},
		},
		{
			name:       "error - db fails on third call",
			propertyID: 10,
			mockSetup: func(m *propServiceMockBookingRepo) {
				var count int
				m.CreatePricingRuleFunc = func(ctx context.Context, rule *domain.PricingRule) error {
					count++
					if count == 3 {
						return errors.New("db insert limit reached")
					}
					return nil
				}
			},
			expectedErr: "db insert limit reached",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mBooking := &propServiceMockBookingRepo{}
			tt.mockSetup(mBooking)
			s := NewPropertyService(&propServiceMockPropertyRepo{}, &propServiceMockUserRepo{}, mBooking)
			err := s.AutoGenerateHighSeasonRules(context.Background(), tt.propertyID)
			if tt.expectedErr != "" {
				if err == nil || err.Error() != tt.expectedErr {
					t.Fatalf("expected error %q, got %v", tt.expectedErr, err)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
			}
		})
	}
}
