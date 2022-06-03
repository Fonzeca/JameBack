package modelview

type ResetPassword struct {
	Token       string
	NewPassword string
	Email       string
}

type Token struct {
	Token              string `json:"token"`
	MustChangePassword bool   `json:"mustChangePassword"`
}
