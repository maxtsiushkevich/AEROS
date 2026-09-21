package postgres

import (
	"users/internal/domain/models"
)

func ToDomain(u *User) *models.User {
	if u == nil {
		return nil
	}

	return &models.User{
		Id:          u.ID,
		Name:        u.Name,
		Birthday:    u.Birthday,
		Email:       u.Email,
		IsActivated: u.IsActivated,
		ActivatedAt: u.ActivatedAt,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}

func ToPersistence(u *models.User) *User {
	if u == nil {
		return nil
	}

	return &User{
		ID:          u.Id,
		Name:        u.Name,
		Birthday:    u.Birthday,
		Email:       u.Email,
		IsActivated: u.IsActivated,
		ActivatedAt: u.ActivatedAt,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}
