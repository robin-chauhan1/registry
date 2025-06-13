package handlers

import (
	"fmt"
	"net/http"
	"ondc-registry/internals/repository"
	"ondc-registry/internals/services"
	"time"

	"github.com/gin-gonic/gin"
)

type Registry struct {
	ID               uint      `json:"id" gorm:"primaryKey"`
	UserID           uint      `json:"user_id" gorm:"not null"` // Foreign key field
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

type RegistryHandler struct {
	registryService *services.RegistryService
}

func NewRegistryHandler(registryService *services.RegistryService) *RegistryHandler {
	if registryService == nil {
		panic("registry service cannot be nil")
	}
	return &RegistryHandler{
		registryService: registryService,
	}
}

func (h *RegistryHandler) GetAllRegistries(c *gin.Context) {
	userID := c.Query("user_id")
	registries, err := h.registryService.GetAllRegistries(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch registries"})
		return
	}
	c.JSON(http.StatusOK, registries)
}

func (h *RegistryHandler) CreateRegistry(c *gin.Context) {
	var registry Registry
	if err := c.ShouldBindJSON(&registry); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	created, err := h.registryService.CreateRegistry(c.Request.Context(), &repository.Registry{
		UserID:           5,
		SubscriberID:     registry.SubscriberID,
		Status:           registry.Status,
		UKID:             registry.UKID,
		SubscriberURL:    registry.SubscriberURL,
		Country:          registry.Country,
		Domain:           registry.Domain,
		ValidFrom:        registry.ValidFrom,
		ValidUntil:       registry.ValidUntil,
		Type:             registry.Type,
		SigningPublicKey: registry.SigningPublicKey,
		EncrPublicKey:    registry.EncrPublicKey,
		BRID:             registry.BRID,
		City:             registry.City,
	})
	if err != nil {
		switch err {
		case repository.ErrInvalidInput:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			fmt.Println("Error creating registry:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create registry"})
		}
		return
	}
	c.JSON(http.StatusCreated, created)
}

func (h *RegistryHandler) UpdateRegistry(c *gin.Context) {
	id := c.Param("id")
	var registry Registry
	if err := c.ShouldBindJSON(&registry); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	fmt.Println("Updating registry with ID:", registry)
	updated, err := h.registryService.UpdateRegistry(c.Request.Context(), id, &repository.Registry{
		ID:               registry.ID,
		UserID:           registry.UserID,
		SubscriberID:     registry.SubscriberID,
		Status:           registry.Status,
		UKID:             registry.UKID,
		SubscriberURL:    registry.SubscriberURL,
		Country:          registry.Country,
		Domain:           registry.Domain,
		ValidFrom:        registry.ValidFrom,
		ValidUntil:       registry.ValidUntil,
		Type:             registry.Type,
		SigningPublicKey: registry.SigningPublicKey,
		EncrPublicKey:    registry.EncrPublicKey,
		BRID:             registry.BRID,
		City:             registry.City,
	})

	fmt.Println("Registry update result:", updated, "Error:", err)
	if err != nil {
		switch err {
		case repository.ErrRegistryNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Registry not found"})
		case repository.ErrInvalidInput:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			fmt.Println("Error updating registry:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update registry"})
		}
		return
	}
	c.JSON(http.StatusOK, updated)
}

func (h *RegistryHandler) DeleteRegistry(c *gin.Context) {
	id := c.Param("id")
	err := h.registryService.DeleteRegistry(c.Request.Context(), id)
	if err != nil {
		switch err {
		case repository.ErrRegistryNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Registry not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete registry"})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Registry deleted successfully"})
}
