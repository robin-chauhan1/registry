package container

import (
	"ondc-registry/internals/api/handlers"
	"ondc-registry/internals/repository"
	"ondc-registry/internals/services"

	"gorm.io/gorm"
)

// Container holds all dependencies
type Container struct {
	DB           *gorm.DB
	RegistryRepo repository.RegistryRepositoryInterface
	RegistrySvc  *services.RegistryService
	RegistryHdl  *handlers.RegistryHandler
	// Add other dependencies here
}

// NewContainer creates a new container with all dependencies
func NewContainer(db *gorm.DB) *Container {
	if db == nil {
		panic("database connection cannot be nil")
	}

	// Initialize repositories
	registryRepo := repository.NewRegistryRepository(db)

	// Initialize services
	registrySvc := services.NewRegistryService(registryRepo)

	// Initialize handlers
	registryHdl := handlers.NewRegistryHandler(registrySvc)

	return &Container{
		DB:           db,
		RegistryRepo: registryRepo,
		RegistrySvc:  registrySvc,
		RegistryHdl:  registryHdl,
	}
}
