package repository

import (
	"context"
	"time"

	"digital-greenhouse/greenhouse-be/internal/domain"

	"gorm.io/gorm"
)

type PropertyDBModel struct {
	ID                uint               `gorm:"primaryKey;autoIncrement"`
	OwnerID           uint               `gorm:"not null"`
	Name              string             `gorm:"type:varchar(150);not null"`
	Description       string             `gorm:"type:text"`
	Address           string             `gorm:"type:varchar(255)"`
	BasePricePerNight float64            `gorm:"type:decimal(10,2);not null"`
	MaxCapacity       int                `gorm:"not null"`
	Status            string             `gorm:"type:enum('ACTIVE', 'INACTIVE', 'MAINTENANCE');default:'ACTIVE'"`
	CreatedAt         time.Time          `gorm:"autoCreateTime"`
	UpdatedAt         time.Time          `gorm:"autoUpdateTime"`
	Images            []PropertyImageDBModel `gorm:"foreignKey:PropertyID"`
}

func (PropertyDBModel) TableName() string {
	return "properties"
}

type PropertyImageDBModel struct {
	ID         uint      `gorm:"primaryKey;autoIncrement"`
	PropertyID uint      `gorm:"not null"`
	ImageData  string    `gorm:"type:longtext;not null"`
	MimeType   string    `gorm:"type:varchar(50);not null"`
	AltText    string    `gorm:"type:varchar(150)"`
	IsCover    bool      `gorm:"default:false"`
	SortOrder  int       `gorm:"default:0"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
}

func (PropertyImageDBModel) TableName() string {
	return "property_images"
}

func toDomainProperty(m PropertyDBModel) *domain.Property {
	p := &domain.Property{
		ID:                m.ID,
		OwnerID:           m.OwnerID,
		Name:              m.Name,
		Description:       m.Description,
		Address:           m.Address,
		BasePricePerNight: m.BasePricePerNight,
		MaxCapacity:       m.MaxCapacity,
		Status:            domain.PropertyStatus(m.Status),
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
	}

	for _, img := range m.Images {
		p.Images = append(p.Images, *toDomainPropertyImage(img))
	}

	return p
}

func toDomainPropertyImage(m PropertyImageDBModel) *domain.PropertyImage {
	return &domain.PropertyImage{
		ID:         m.ID,
		PropertyID: m.PropertyID,
		ImageData:  m.ImageData,
		MimeType:   m.MimeType,
		AltText:    m.AltText,
		IsCover:    m.IsCover,
		SortOrder:  m.SortOrder,
		CreatedAt:  m.CreatedAt,
	}
}

type propertyRepository struct {
	db *gorm.DB
}

func NewPropertyRepository(db *gorm.DB) domain.PropertyRepository {
	return &propertyRepository{db: db}
}

func (r *propertyRepository) Create(ctx context.Context, p *domain.Property) error {
	m := PropertyDBModel{
		OwnerID:           p.OwnerID,
		Name:              p.Name,
		Description:       p.Description,
		Address:           p.Address,
		BasePricePerNight: p.BasePricePerNight,
		MaxCapacity:       p.MaxCapacity,
		Status:            string(p.Status),
	}

	err := r.db.WithContext(ctx).Create(&m).Error
	if err == nil {
		p.ID = m.ID
		p.CreatedAt = m.CreatedAt
		p.UpdatedAt = m.UpdatedAt
	}
	return err
}

func (r *propertyRepository) GetAll(ctx context.Context) ([]domain.Property, error) {
	var models []PropertyDBModel
	err := r.db.WithContext(ctx).Preload("Images").Find(&models).Error
	if err != nil {
		return nil, err
	}

	var props []domain.Property
	for _, m := range models {
		props = append(props, *toDomainProperty(m))
	}
	return props, nil
}

func (r *propertyRepository) GetByID(ctx context.Context, id uint) (*domain.Property, error) {
	var m PropertyDBModel
	err := r.db.WithContext(ctx).Preload("Images").First(&m, id).Error
	if err != nil {
		return nil, err
	}
	return toDomainProperty(m), nil
}

func (r *propertyRepository) GetByOwnerID(ctx context.Context, ownerID uint) ([]domain.Property, error) {
	var models []PropertyDBModel
	err := r.db.WithContext(ctx).Preload("Images").Where("owner_id = ?", ownerID).Find(&models).Error
	if err != nil {
		return nil, err
	}

	var props []domain.Property
	for _, m := range models {
		props = append(props, *toDomainProperty(m))
	}
	return props, nil
}

func (r *propertyRepository) Update(ctx context.Context, p *domain.Property) error {
	m := PropertyDBModel{
		ID:                p.ID,
		OwnerID:           p.OwnerID,
		Name:              p.Name,
		Description:       p.Description,
		Address:           p.Address,
		BasePricePerNight: p.BasePricePerNight,
		MaxCapacity:       p.MaxCapacity,
		Status:            string(p.Status),
	}
	return r.db.WithContext(ctx).Save(&m).Error
}

func (r *propertyRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&PropertyDBModel{}, id).Error
}

func (r *propertyRepository) AddImage(ctx context.Context, img *domain.PropertyImage) error {
	m := PropertyImageDBModel{
		PropertyID: img.PropertyID,
		ImageData:  img.ImageData,
		MimeType:   img.MimeType,
		AltText:    img.AltText,
		IsCover:    img.IsCover,
		SortOrder:  img.SortOrder,
	}
	err := r.db.WithContext(ctx).Create(&m).Error
	if err == nil {
		img.ID = m.ID
		img.CreatedAt = m.CreatedAt
	}
	return err
}

func (r *propertyRepository) GetImageByID(ctx context.Context, id uint) (*domain.PropertyImage, error) {
	var m PropertyImageDBModel
	err := r.db.WithContext(ctx).First(&m, id).Error
	if err != nil {
		return nil, err
	}
	return toDomainPropertyImage(m), nil
}

func (r *propertyRepository) UpdateImage(ctx context.Context, img *domain.PropertyImage) error {
	m := PropertyImageDBModel{
		ID:         img.ID,
		PropertyID: img.PropertyID,
		ImageData:  img.ImageData,
		MimeType:   img.MimeType,
		AltText:    img.AltText,
		IsCover:    img.IsCover,
		SortOrder:  img.SortOrder,
	}
	return r.db.WithContext(ctx).Save(&m).Error
}

func (r *propertyRepository) DeleteImage(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&PropertyImageDBModel{}, id).Error
}
