package dto

type AuthRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type TokensResponse struct {
	Access_token string `json:"access_token" `
}

type PasswordUpdateRequest struct {
	OldPassword *string `json:"old_password" validate:"required,min=8"`
	NewPassword *string `json:"new_password" validate:"required,min=8"`
}
