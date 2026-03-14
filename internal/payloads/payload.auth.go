package payloads

import "homedy/internal/models"

type RequestSignUp struct {
	Username   string `json:"username" validate:"required,username" example:"chesta_ardiona"`
	Email      string `json:"email" validate:"required,email" example:"chestaardi4@gmail.com"`
	Password   string `json:"password" validate:"required,password" example:"YourPassword123"`
	RememberMe bool   `json:"remember_me"`
}

func (p *RequestSignUp) ToUser() models.User {
	return models.User{
		Username: p.Username,
		Email:    p.Email,
		Password: p.Password,
	}
}
