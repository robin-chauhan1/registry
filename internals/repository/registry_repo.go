package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

var (
	ErrRegistryNotFound = errors.New("registry not found")
	ErrInvalidInput     = errors.New("invalid input")
)

type RegistryRepositoryInterface interface {
	GetAllRegistries(ctx context.Context, userID string) ([]Registry, error)
	CreateRegistry(ctx context.Context, r *Registry) (*Registry, error)
	UpdateRegistry(ctx context.Context, id string, r *Registry) (*Registry, error)
	DeleteRegistry(ctx context.Context, id string) error
}

type RegistryRepository struct {
	db *gorm.DB
}

type Registry struct {
	ID               uint      `json:"id" gorm:"primaryKey"`
	UserID           uint      `json:"user_id" gorm:"not null"` // Foreign key field
	SubscriberID     string    `json:"subscriber_id" gorm:"not null"`
	Status           string    `json:"status" gorm:"not null"`
	UKID             string    `json:"ukId" gorm:"not null"`
	SubscriberURL    string    `json:"subscriber_url" gorm:"not null"`
	Country          string    `json:"country" gorm:"not null"`
	Domain           string    `json:"domain" gorm:"not null"`
	ValidFrom        time.Time `json:"valid_from" gorm:"not null"`
	ValidUntil       time.Time `json:"valid_until" gorm:"not null"`
	Type             string    `json:"type" gorm:"not null"`
	SigningPublicKey string    `json:"signing_public_key" gorm:"not null"`
	EncrPublicKey    string    `json:"encr_public_key" gorm:"not null"`
	Created          time.Time `json:"created" gorm:"autoCreateTime"`
	Updated          time.Time `json:"updated" gorm:"autoUpdateTime"`
	BRID             string    `json:"br_id" gorm:"not null"`
	City             string    `json:"city" gorm:"not null"`
}

// NewRegistryRepository creates a new registry repository with the given database connection
func NewRegistryRepository(db *gorm.DB) RegistryRepositoryInterface {
	if db == nil {
		panic("database connection cannot be nil")
	}
	return &RegistryRepository{db: db}
}

// validateRegistry checks if the registry has all required fields
func validateRegistry(r *Registry) error {
	if r == nil {
		return ErrInvalidInput
	}
	if r.SubscriberID == "" || r.UKID == "" || r.BRID == "" {
		return fmt.Errorf("%w: missing required fields", ErrInvalidInput)
	}
	if r.ValidFrom.IsZero() || r.ValidUntil.IsZero() {
		return fmt.Errorf("%w: invalid validity period", ErrInvalidInput)
	}
	if r.ValidUntil.Before(r.ValidFrom) {
		return fmt.Errorf("%w: valid_until must be after valid_from", ErrInvalidInput)
	}
	return nil
}

func (rr *RegistryRepository) GetAllRegistries(ctx context.Context, userID string) ([]Registry, error) {
	var registries []Registry
	if err := rr.db.WithContext(ctx).Where("user_id = ?", userID).Find(&registries).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch registries: %w", err)
	}
	return registries, nil
}

func (rr *RegistryRepository) CreateRegistry(ctx context.Context, r *Registry) (*Registry, error) {
	if err := validateRegistry(r); err != nil {
		return nil, err
	}

	if err := rr.db.WithContext(ctx).Create(r).Error; err != nil {
		return nil, fmt.Errorf("failed to create registry: %w", err)
	}
	return r, nil
}

func (rr *RegistryRepository) UpdateRegistry(ctx context.Context, id string, r *Registry) (*Registry, error) {
	// Load the existing registry
	var existing Registry
	if err := rr.db.WithContext(ctx).First(&existing, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRegistryNotFound
		}
		return nil, fmt.Errorf("failed to find registry: %w", err)
	}

	// Update only non-zero fields from the input registry
	if err := rr.db.WithContext(ctx).Model(&existing).Updates(r).Error; err != nil {
		return nil, fmt.Errorf("failed to update registry: %w", err)
	}

	return &existing, nil
}

func (rr *RegistryRepository) DeleteRegistry(ctx context.Context, id string) error {
	result := rr.db.WithContext(ctx).Delete(&Registry{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete registry: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrRegistryNotFound
	}
	return nil
}
