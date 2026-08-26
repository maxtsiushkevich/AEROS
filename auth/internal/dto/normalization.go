package dto

import "strings"

func (r *AuthRequest) Normalize() {
	r.Email = strings.TrimSpace(r.Email)
}
