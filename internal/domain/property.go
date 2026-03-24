package domain

import (
	"context"
	"time"
)

type PropertyStatus string

const (
	PropertyActive      PropertyStatus = "ACTIVE"
	PropertyInactive    PropertyStatus = "INACTIVE"
	PropertyMaintenance PropertyStatus = "MAINTENANCE"
)

type Property struct {
	ID                uint
	OwnerID           uint
	Name              string
	Description       string
	Address           string
	BasePricePerNight float64
	MaxCapacity       int
	Status            PropertyStatus
	CreatedAt         time.Time
	UpdatedAt         time.Time
	Images            []PropertyImage
}

type PropertyImage struct {
	ID         uint
	PropertyID uint
	ImageData  string // Base64
	MimeType   string
	AltText    string
	IsCover    bool
	SortOrder  int
	CreatedAt  time.Time
}

type PropertyRepository interface {
	Create(ctx context.Context, property *Property) error
	GetAll(ctx context.Context) ([]Property, error)
	GetByID(ctx context.Context, id uint) (*Property, error)
	GetByOwnerID(ctx context.Context, ownerID uint) ([]Property, error)
	Update(ctx context.Context, property *Property) error
	Delete(ctx context.Context, id uint) error

	AddImage(ctx context.Context, image *PropertyImage) error
	GetImageByID(ctx context.Context, id uint) (*PropertyImage, error)
	UpdateImage(ctx context.Context, image *PropertyImage) error
	DeleteImage(ctx context.Context, id uint) error
}

type PropertyService interface {
	CreateProperty(ctx context.Context, property *Property) error
	ListProperties(ctx context.Context) ([]Property, error)
	GetPropertyByID(ctx context.Context, id uint) (*Property, error)
	GetPropertiesByOwner(ctx context.Context, ownerID uint) ([]Property, error)
	UpdateProperty(ctx context.Context, property *Property) error
	DeleteProperty(ctx context.Context, id uint) error

	AddPropertyImage(ctx context.Context, image *PropertyImage) error
	UpdatePropertyImage(ctx context.Context, image *PropertyImage) error
	DeletePropertyImage(ctx context.Context, imageID uint) error
}
