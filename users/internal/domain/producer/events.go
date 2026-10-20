package producer

type SendOTPEvent struct {
	Email   string     `json:"email"`
	OTP     string     `json:"otp"`
	Purpose OTPPurpose `json:"purpose"`
}

type OTPPurpose string

const (
	ActivateProfile OTPPurpose = "ActivateProfile"
	ChangeEmail     OTPPurpose = "ChangeEmail"
)

func GetOTPEventMessage(email string, otp string, purpose OTPPurpose) SendOTPEvent {
	return SendOTPEvent{
		Email:   email,
		OTP:     otp,
		Purpose: purpose,
	}
}
