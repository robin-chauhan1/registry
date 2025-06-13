package models

import (
	"time"
)

// Registry represents the schema for the registry model
type Registry struct {
	ID               uint      `json:"id" gorm:"primaryKey"`
	UserID           uint      `json:"user_id" gorm:"not null"` // Foreign key field
	User             User      `json:"user" gorm:"foreignKey:UserID;references:ID" validate:"-"`
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
