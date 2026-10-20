package otp

type OTPPurpose string

const (
	ActivateProfile OTPPurpose = "ActivateProfile"
	ChangeEmail     OTPPurpose = "ChangeEmail"
)
