package requests

type UserCreateRequest struct {
	Email    string `json:"email" validate:"required,email,unique"`
	Password string `json:"password" validate:"required,min=6,max=32"`
}
