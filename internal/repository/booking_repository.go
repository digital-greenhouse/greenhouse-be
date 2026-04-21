package repository

import (
	"context"
	"time"

	"digital-greenhouse/greenhouse-be/internal/domain"

	"gorm.io/gorm"
)

// DB Models
type PricingRuleDBModel struct {
	ID            uint      `gorm:"primaryKey;autoIncrement"`
	PropertyID    uint      `gorm:"not null"`
	Name          string    `gorm:"type:varchar(100);not null"`
	StartDate     time.Time `gorm:"type:date;not null"`
	EndDate       time.Time `gorm:"type:date;not null"`
	PriceModifier float64   `gorm:"type:decimal(5,2);not null"`
	Description   string    `gorm:"type:varchar(255)"`
	IsActive      bool      `gorm:"default:true"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
}

func (PricingRuleDBModel) TableName() string {
	return "pricing_rules"
}

type QuoteDBModel struct {
	ID                uint      `gorm:"primaryKey;autoIncrement"`
	PropertyID        uint      `gorm:"not null"`
	ClientID          *uint     `gorm:"default:null"`
	CheckInDate       time.Time `gorm:"type:date;not null"`
	CheckOutDate      time.Time `gorm:"type:date;not null"`
	GuestCount        int       `gorm:"not null"`
	CalculatedTotal   float64   `gorm:"type:decimal(10,2);not null"`
	NightsCount       int       `gorm:"not null"`
	AppliedModifier   float64   `gorm:"type:decimal(5,2);default:1.00"`
	Status            string    `gorm:"type:enum('ACTIVE','CONVERTED','EXPIRED','ABANDONED');default:'ACTIVE'"`
	AbandonmentReason string    `gorm:"type:varchar(255)"`
	ExpiresAt         *time.Time
	CreatedAt         time.Time `gorm:"autoCreateTime"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime"`
}

func (QuoteDBModel) TableName() string {
	return "quotes"
}

type BookingDBModel struct {
	ID                 uint      `gorm:"primaryKey;autoIncrement"`
	PropertyID         uint      `gorm:"not null"`
	ClientID           uint      `gorm:"not null"`
	QuoteID            *uint     `gorm:"default:null"`
	CheckInDate        time.Time `gorm:"type:date;not null"`
	CheckOutDate       time.Time `gorm:"type:date;not null"`
	GuestCount         int       `gorm:"not null"`
	NightsCount        int       `gorm:"not null"`
	TotalPrice         float64   `gorm:"type:decimal(10,2);not null"`
	Status             string    `gorm:"type:enum('PENDING_PAYMENT','CONFIRMED','CANCELLED','COMPLETED');default:'PENDING_PAYMENT'"`
	CancellationReason string    `gorm:"type:varchar(255)"`
	SpecialRequests    string    `gorm:"type:text"`
	CreatedAt          time.Time `gorm:"autoCreateTime"`
	UpdatedAt          time.Time `gorm:"autoUpdateTime"`
}

func (BookingDBModel) TableName() string {
	return "bookings"
}

// Mappers
func toDomainPricingRule(m PricingRuleDBModel) domain.PricingRule {
	return domain.PricingRule{
		ID:            m.ID,
		PropertyID:    m.PropertyID,
		Name:          m.Name,
		StartDate:     m.StartDate,
		EndDate:       m.EndDate,
		PriceModifier: m.PriceModifier,
		Description:   m.Description,
		IsActive:      m.IsActive,
		CreatedAt:     m.CreatedAt,
	}
}

func fromDomainPricingRule(d *domain.PricingRule) PricingRuleDBModel {
	return PricingRuleDBModel{
		ID:            d.ID,
		PropertyID:    d.PropertyID,
		Name:          d.Name,
		StartDate:     d.StartDate,
		EndDate:       d.EndDate,
		PriceModifier: d.PriceModifier,
		Description:   d.Description,
		IsActive:      d.IsActive,
	}
}

func toDomainQuote(m QuoteDBModel) *domain.Quote {
	return &domain.Quote{
		ID:                m.ID,
		PropertyID:        m.PropertyID,
		ClientID:          m.ClientID,
		CheckInDate:       m.CheckInDate,
		CheckOutDate:      m.CheckOutDate,
		GuestCount:        m.GuestCount,
		CalculatedTotal:   m.CalculatedTotal,
		NightsCount:       m.NightsCount,
		AppliedModifier:   m.AppliedModifier,
		Status:            domain.QuoteStatus(m.Status),
		AbandonmentReason: m.AbandonmentReason,
		ExpiresAt:         m.ExpiresAt,
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
	}
}

func toDomainBooking(m BookingDBModel) *domain.Booking {
	return &domain.Booking{
		ID:                 m.ID,
		PropertyID:         m.PropertyID,
		ClientID:           m.ClientID,
		QuoteID:            m.QuoteID,
		CheckInDate:        m.CheckInDate,
		CheckOutDate:       m.CheckOutDate,
		GuestCount:         m.GuestCount,
		NightsCount:        m.NightsCount,
		TotalPrice:         m.TotalPrice,
		Status:             domain.BookingStatus(m.Status),
		CancellationReason: m.CancellationReason,
		SpecialRequests:    m.SpecialRequests,
		CreatedAt:          m.CreatedAt,
		UpdatedAt:          m.UpdatedAt,
	}
}

// Implementation
type bookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) domain.BookingRepository {
	return &bookingRepository{db: db}
}

func (r *bookingRepository) CreateQuote(ctx context.Context, q *domain.Quote) error {
	m := QuoteDBModel{
		PropertyID:      q.PropertyID,
		ClientID:        q.ClientID,
		CheckInDate:     q.CheckInDate,
		CheckOutDate:    q.CheckOutDate,
		GuestCount:      q.GuestCount,
		CalculatedTotal: q.CalculatedTotal,
		NightsCount:     q.NightsCount,
		AppliedModifier: q.AppliedModifier,
		Status:          string(q.Status),
		ExpiresAt:       q.ExpiresAt,
	}
	err := r.db.WithContext(ctx).Create(&m).Error
	if err == nil {
		q.ID = m.ID
		q.CreatedAt = m.CreatedAt
		q.UpdatedAt = m.UpdatedAt
	}
	return err
}

func (r *bookingRepository) GetQuoteByID(ctx context.Context, id uint) (*domain.Quote, error) {
	var m QuoteDBModel
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}
	return toDomainQuote(m), nil
}

func (r *bookingRepository) GetQuotesByClientID(ctx context.Context, clientID uint) ([]domain.Quote, error) {
	var models []QuoteDBModel
	if err := r.db.WithContext(ctx).Where("client_id = ?", clientID).Find(&models).Error; err != nil {
		return nil, err
	}
	quotes := make([]domain.Quote, len(models))
	for i, m := range models {
		quotes[i] = *toDomainQuote(m)
	}
	return quotes, nil
}

func (r *bookingRepository) UpdateQuoteStatus(ctx context.Context, id uint, status domain.QuoteStatus) error {
	return r.db.WithContext(ctx).Model(&QuoteDBModel{}).Where("id = ?", id).Update("status", string(status)).Error
}

func (r *bookingRepository) CreateBooking(ctx context.Context, b *domain.Booking) error {
	m := BookingDBModel{
		PropertyID:      b.PropertyID,
		ClientID:        b.ClientID,
		QuoteID:         b.QuoteID,
		CheckInDate:     b.CheckInDate,
		CheckOutDate:    b.CheckOutDate,
		GuestCount:      b.GuestCount,
		NightsCount:     b.NightsCount,
		TotalPrice:      b.TotalPrice,
		Status:          string(b.Status),
		SpecialRequests: b.SpecialRequests,
	}
	err := r.db.WithContext(ctx).Create(&m).Error
	if err == nil {
		b.ID = m.ID
		b.CreatedAt = m.CreatedAt
		b.UpdatedAt = m.UpdatedAt
	}
	return err
}

func (r *bookingRepository) GetBookingByID(ctx context.Context, id uint) (*domain.Booking, error) {
	var m BookingDBModel
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}
	return toDomainBooking(m), nil
}

func (r *bookingRepository) GetBookingsByClientID(ctx context.Context, clientID uint) ([]domain.Booking, error) {
	var models []BookingDBModel
	if err := r.db.WithContext(ctx).Where("client_id = ?", clientID).Find(&models).Error; err != nil {
		return nil, err
	}
	bookings := make([]domain.Booking, len(models))
	for i, m := range models {
		bookings[i] = *toDomainBooking(m)
	}
	return bookings, nil
}

func (r *bookingRepository) GetBookingsByPropertyID(ctx context.Context, propertyID uint) ([]domain.Booking, error) {
	var models []BookingDBModel
	if err := r.db.WithContext(ctx).Where("property_id = ?", propertyID).Find(&models).Error; err != nil {
		return nil, err
	}
	bookings := make([]domain.Booking, len(models))
	for i, m := range models {
		bookings[i] = *toDomainBooking(m)
	}
	return bookings, nil
}

func (r *bookingRepository) UpdateBookingStatus(ctx context.Context, id uint, status domain.BookingStatus, reason string) error {
	updates := map[string]interface{}{
		"status": string(status),
	}
	if reason != "" {
		updates["cancellation_reason"] = reason
	}
	return r.db.WithContext(ctx).Model(&BookingDBModel{}).Where("id = ?", id).Updates(updates).Error
}

func (r *bookingRepository) CheckAvailability(ctx context.Context, propertyID uint, checkIn, checkOut time.Time) (bool, error) {
	var count int64
	// Una reserva se solapa si: (nueva_entrada < reserva_salida) Y (nueva_salida > reserva_entrada)
	// Solo contamos reservas CONFIRMED o PENDING_PAYMENT
	err := r.db.WithContext(ctx).Model(&BookingDBModel{}).
		Where("property_id = ? AND status IN (?, ?) AND check_in_date < ? AND check_out_date > ?",
			propertyID, string(domain.BookingConfirmed), string(domain.BookingPending), checkOut, checkIn).
		Count(&count).Error

	return count == 0, err
}

func (r *bookingRepository) CreatePricingRule(ctx context.Context, rule *domain.PricingRule) error {
	m := fromDomainPricingRule(rule)
	err := r.db.WithContext(ctx).Create(&m).Error
	if err == nil {
		rule.ID = m.ID
		rule.CreatedAt = m.CreatedAt
	}
	return err
}

func (r *bookingRepository) GetPricingRulesByPropertyID(ctx context.Context, propertyID uint, start, end time.Time) ([]domain.PricingRule, error) {
	var models []PricingRuleDBModel
	// Buscamos reglas que se solapen con el rango de fechas pedido y estén activas
	err := r.db.WithContext(ctx).
		Where("property_id = ? AND is_active = TRUE AND ((start_date <= ? AND end_date >= ?) OR (start_date <= ? AND end_date >= ?))",
			propertyID, end, start, end, start).
		Find(&models).Error
	if err != nil {
		return nil, err
	}
	rules := make([]domain.PricingRule, len(models))
	for i, m := range models {
		rules[i] = toDomainPricingRule(m)
	}
	return rules, nil
}

func (r *bookingRepository) GetAllPricingRulesByPropertyID(ctx context.Context, propertyID uint) ([]domain.PricingRule, error) {
	var models []PricingRuleDBModel
	err := r.db.WithContext(ctx).Where("property_id = ?", propertyID).Order("start_date ASC").Find(&models).Error
	if err != nil {
		return nil, err
	}
	rules := make([]domain.PricingRule, len(models))
	for i, m := range models {
		rules[i] = toDomainPricingRule(m)
	}
	return rules, nil
}

func (r *bookingRepository) DeletePricingRule(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&PricingRuleDBModel{}, id).Error
}
