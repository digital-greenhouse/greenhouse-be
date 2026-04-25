package dto

type PropertyImageDTO struct {
	ID         uint   `json:"id,omitempty"`
	PropertyID uint   `json:"property_id,omitempty"`
	ImageData  string `json:"image_data"`
	MimeType   string `json:"mime_type"`
	AltText    string `json:"alt_text,omitempty"`
	IsCover    bool   `json:"is_cover"`
	SortOrder  int    `json:"sort_order"`
}

type CreatePropertyRequest struct {
	OwnerID           uint               `json:"owner_id"`
	Name              string             `json:"name"`
	Description       string             `json:"description"`
	Address           string             `json:"address"`
	BasePricePerNight float64            `json:"base_price_per_night"`
	MaxCapacity       int                `json:"max_capacity"`
	Images            []PropertyImageDTO `json:"images,omitempty"`
}

type PropertyResponse struct {
	ID                uint               `json:"id"`
	OwnerID           uint               `json:"owner_id"`
	Name              string             `json:"name"`
	Description       string             `json:"description"`
	Address           string             `json:"address"`
	BasePricePerNight float64            `json:"base_price_per_night"`
	MaxCapacity       int                `json:"max_capacity"`
	Status            string             `json:"status"`
	CreatedAt         string             `json:"created_at"`
	UpdatedAt         string             `json:"updated_at"`
	Images            []PropertyImageDTO `json:"images,omitempty"`
}

type PricingRuleDTO struct {
	ID            uint    `json:"id"`
	PropertyID    uint    `json:"property_id"`
	Name          string  `json:"name"`
	StartDate     string  `json:"start_date"`
	EndDate       string  `json:"end_date"`
	PriceModifier float64 `json:"price_modifier"`
	Description   string  `json:"description"`
	IsActive      bool    `json:"is_active"`
}

type CreatePricingRuleRequest struct {
	Name          string  `json:"name"`
	StartDate     string  `json:"start_date"` // format: 2006-01-02
	EndDate       string  `json:"end_date"`   // format: 2006-01-02
	PriceModifier float64 `json:"price_modifier"`
	Description   string  `json:"description"`
}
