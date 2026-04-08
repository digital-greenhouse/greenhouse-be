package repository

import (
	"context"
	"time"

	"digital-greenhouse/greenhouse-be/internal/domain"

	"gorm.io/gorm"
)

type PaymentDBModel struct {
	ID              uint      `gorm:"primaryKey;autoIncrement"`
	BookingID       uint      `gorm:"not null"`
	Amount          float64   `gorm:"type:decimal(10,2);not null"`
	PaymentMethod   string    `gorm:"type:enum('TRANSFERENCIA','EFECTIVO');not null"`
	ProofData       string    `gorm:"type:longtext"`
	ProofMimeType   string    `gorm:"type:varchar(50)"`
	Status          string    `gorm:"type:enum('PENDING_VERIFICATION','VERIFIED','REJECTED');default:'PENDING_VERIFICATION'"`
	RejectionReason string    `gorm:"type:varchar(255)"`
	VerifiedBy      *uint     `gorm:"default:null"`
	PaymentDate     time.Time `gorm:"autoCreateTime"`
	VerifiedAt      *time.Time
}

func (PaymentDBModel) TableName() string {
	return "payments"
}

func toDomainPayment(m PaymentDBModel) *domain.Payment {
	return &domain.Payment{
		ID:              m.ID,
		BookingID:       m.BookingID,
		Amount:          m.Amount,
		PaymentMethod:   domain.PaymentMethod(m.PaymentMethod),
		ProofData:       m.ProofData,
		ProofMimeType:   m.ProofMimeType,
		Status:          domain.PaymentStatus(m.Status),
		RejectionReason: m.RejectionReason,
		VerifiedBy:      m.VerifiedBy,
		PaymentDate:     m.PaymentDate,
		VerifiedAt:      m.VerifiedAt,
	}
}

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) domain.PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) Create(ctx context.Context, p *domain.Payment) error {
	m := PaymentDBModel{
		BookingID:     p.BookingID,
		Amount:        p.Amount,
		PaymentMethod: string(p.PaymentMethod),
		ProofData:     p.ProofData,
		ProofMimeType: p.ProofMimeType,
		Status:        string(p.Status),
	}
	err := r.db.WithContext(ctx).Create(&m).Error
	if err == nil {
		p.ID = m.ID
		p.PaymentDate = m.PaymentDate
	}
	return err
}

func (r *paymentRepository) GetByID(ctx context.Context, id uint) (*domain.Payment, error) {
	var m PaymentDBModel
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}
	return toDomainPayment(m), nil
}

func (r *paymentRepository) GetByBookingID(ctx context.Context, bookingID uint) ([]domain.Payment, error) {
	var models []PaymentDBModel
	if err := r.db.WithContext(ctx).Where("booking_id = ?", bookingID).Find(&models).Error; err != nil {
		return nil, err
	}
	payments := make([]domain.Payment, len(models))
	for i, m := range models {
		payments[i] = *toDomainPayment(m)
	}
	return payments, nil
}

func (r *paymentRepository) UpdateStatus(ctx context.Context, id uint, status domain.PaymentStatus, verifierID *uint, reason string) error {
	updates := map[string]interface{}{
		"status": string(status),
	}
	if verifierID != nil {
		updates["verified_by"] = verifierID
		now := time.Now()
		updates["verified_at"] = &now
	}
	if reason != "" {
		updates["rejection_reason"] = reason
	}
	return r.db.WithContext(ctx).Model(&PaymentDBModel{}).Where("id = ?", id).Updates(updates).Error
}
