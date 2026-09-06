package dto

import "strings"

func (r *AuthRequest) Normalize() {
	r.Email = strings.TrimSpace(r.Email)
}

func (r *PasswordUpdateRequest) Normalize() {
	r.OldPassword = strings.TrimSpace(r.OldPassword)
	r.NewPassword = strings.TrimSpace(r.NewPassword)
}
