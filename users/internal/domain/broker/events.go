package broker

import "users/internal/domain/otp"

type SendOTPEvent struct {
	Email   string         `json:"email"`
	OTP     string         `json:"otp"`
	Purpose otp.OTPPurpose `json:"purpose"`
}

func GetOTPEventMessage(email string, otp string, purpose otp.OTPPurpose) SendOTPEvent {
	return SendOTPEvent{
		Email:   email,
		OTP:     otp,
		Purpose: purpose,
	}
}
