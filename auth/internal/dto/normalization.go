package dto

import "strings"

func (r *AuthRequest) Normalize() {
	r.Email = strings.TrimSpace(r.Email)
}

func (r *PasswordUpdateRequest) Normalize() {
	r.OldPassword = trim(r.OldPassword)
	r.NewPassword = trim(r.NewPassword)
}

func trim(s *string) *string {
	if s == nil {
		return nil
	}
	res := strings.TrimSpace(*s)
	return &res
}
