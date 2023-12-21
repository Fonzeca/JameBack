package modelview

type ResetPassword struct {
	Token       string
	NewPassword string
	Email       string
}

type LoginResponse struct {
	MustChangePassword bool     `json:"mustChangePassword"`
	Username           string   `json:"username"`
	Admin              bool     `json:"admin"`
	Roles              []string `json:"roles"`
	FullName           string   `json:"fullName"`
	Expiration         int64    `json:"expiration"`
}
