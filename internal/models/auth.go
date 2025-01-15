package models

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=30,alphanum"`
	Password string `json:"password" binding:"required,min=6,max=50"`
	Email    string `json:"email" binding:"required,email"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required,min=3,max=30"`
	Password string `json:"password" binding:"required,min=6"`
}
