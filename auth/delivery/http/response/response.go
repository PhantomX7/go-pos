package response

type AuthResponse struct {
	Token    string `json:"token"`
	Username string `json:"username"`
}

type GetMeResponse struct {
	Username string `json:"username"`
	Role     string `json:"role"`
}
