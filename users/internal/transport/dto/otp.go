package dto

import "users/internal/domain/otp"

type OTPRequest struct {
	Purpose otp.OTPPurpose `json:"purpose" binding:"required,oneof=ActivateProfile ChangeEmail"`
}
