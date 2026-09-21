package dto

import "users/internal/domain/models"

func ToUserResponse(u *models.User) *UserResponse {
	r := &UserResponse{
		Id:          u.Id,
		Name:        u.Name,
		Birthday:    u.Birthday,
		Email:       u.Email,
		IsActivated: u.IsActivated,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}

	if u.ActivatedAt != nil {
		r.ActivatedAt = u.ActivatedAt
	}

	return r
}
