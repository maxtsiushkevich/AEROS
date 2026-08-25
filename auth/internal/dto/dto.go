package dto

type AuthRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type TokensResponse struct {
	Access_token string `json:"access_token" `
}
