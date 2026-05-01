package dto

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"` //邮箱
}

type ForgotPasswordResponse struct {
	Message string `json:"message"`
}
type ResetPasswordRequest struct {
	Token                string `json:"token" binding:"required"`
	Password             string `json:"password" binding:"required,min=6"`
	PasswordConfirmation string `json:"password_confirmation" binding:"required,eqfield=Password"`
}

type ResetPasswordResponse struct {
	Message string `json:"message"`
}

type ValidateTokenResponse struct {
	Valid bool `json:"valid"`
}

type ResetPasswordWithTokenRequest struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

type ResetPasswordWithTokenResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
