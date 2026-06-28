package user

type LoginRequest struct {
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,password"`
}

type CreateUserRequest struct {
	Username     string `json:"username" binding:"required,max=16"`
	Email    string `json:"email" binding:"required,min=8,email"`
	Password string `json:"password" binding:"required,password"`
}

type DeleteUserRequest struct {
	Id string `json:"id" binding:"required"`
}