package services

import (
	"context"
	"errors"
	"ondc-registry/internals/repository"
)

type RegistryService struct {
	repo repository.RegistryRepositoryInterface
}

// NewRegistryService creates a new registry service with the given repository
func NewRegistryService(repo repository.RegistryRepositoryInterface) *RegistryService {
	if repo == nil {
		panic("repository cannot be nil")
	}
	return &RegistryService{
		repo: repo,
	}
}

func (s *RegistryService) GetAllRegistries(ctx context.Context, userID string) ([]repository.Registry, error) {
	return s.repo.GetAllRegistries(ctx, userID)
}

// CreateRegistry creates a new registry.
func (s *RegistryService) CreateRegistry(ctx context.Context, r *repository.Registry) (*repository.Registry, error) {
	if r == nil {
		return nil, errors.New("registry cannot be nil")
	}
	return s.repo.CreateRegistry(ctx, r)
}

// UpdateRegistry updates an existing registry by ID.
func (s *RegistryService) UpdateRegistry(ctx context.Context, id string, r *repository.Registry) (*repository.Registry, error) {
	return s.repo.UpdateRegistry(ctx, id, r)
}

// DeleteRegistry deletes a registry by ID.
func (s *RegistryService) DeleteRegistry(ctx context.Context, id string) error {
	return s.repo.DeleteRegistry(ctx, id)
}
