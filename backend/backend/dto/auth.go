package dto

// LoginRequest represents a login request
type LoginRequest struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RememberMe bool   `json:"rememberMe"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}
