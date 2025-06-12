package services

import (
	"ondc-registry/internals/models"
	"ondc-registry/internals/repository"
)

func GetAllUsers() ([]models.User, error) {
	return repository.FindAllUsers()
}
