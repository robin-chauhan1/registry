package models

import (
	"time"
)

// Registry represents the schema for the registry model
type Registry struct {
	ID               uint      `json:"id" gorm:"primaryKey"`
	UserID           uint      `json:"user_id" gorm:"not null"` // Foreign key field
	User             User      `json:"user" gorm:"foreignKey:UserID;references:ID" validate:"-"`
	SubscriberID     string    `json:"subscriber_id" gorm:"not null" validate:"required"`
	Status           string    `json:"status" gorm:"not null" validate:"required"`
	UKID             string    `json:"ukId" gorm:"not null" validate:"required"`
	SubscriberURL    string    `json:"subscriber_url" gorm:"not null" validate:"required"`
	Country          string    `json:"country" gorm:"not null" validate:"required"`
	Domain           string    `json:"domain" gorm:"not null" validate:"required"`
	ValidFrom        time.Time `json:"valid_from" gorm:"not null" validate:"required"`
	ValidUntil       time.Time `json:"valid_until" gorm:"not null" validate:"required"`
	Type             string    `json:"type" gorm:"not null" validate:"required"`
	SigningPublicKey string    `json:"signing_public_key" gorm:"not null" validate:"required"`
	EncrPublicKey    string    `json:"encr_public_key" gorm:"not null" validate:"required"`
	Created          time.Time `json:"created" gorm:"autoCreateTime"`
	Updated          time.Time `json:"updated" gorm:"autoUpdateTime"`
	BRID             string    `json:"br_id" gorm:"not null" validate:"required"`
	City             string    `json:"city" gorm:"not null" validate:"required"`
}
