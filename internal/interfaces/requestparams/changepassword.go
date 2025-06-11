package requestparams

// ChangePasswordParams defines the parameters for password change
type ChangePasswordParams struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required"`
	ConfirmPassword string `json:"confirm_password" validate:"required"`
}
