package modelview

type ResetPassword struct {
	Token       string
	NewPassword string
	Email       string
}

type LoginResponse struct {
	MustChangePassword bool `json:"mustChangePassword"`
}
