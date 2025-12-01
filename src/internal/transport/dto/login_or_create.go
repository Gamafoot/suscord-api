package dto

type LoginOrCreateRequest struct {
	Username string `json:"username" validate:"required,gte=1,lte=15"`
	Password string `json:"password" validate:"required"`
}
