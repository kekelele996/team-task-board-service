package dto

// RegisterRequest is the payload used to create a new account.
type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=64"`
	Email    string `json:"email" validate:"required,email,max=128"`
	Password string `json:"password" validate:"required,min=6,max=128"`
}

// LoginRequest is the payload used to obtain a JWT.
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=128"`
}

// TokenResponse is returned after login/register.
type TokenResponse struct {
	Token string  `json:"token"`
	User  UserDTO `json:"user"`
}

// UserDTO is the public user representation.
type UserDTO struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
