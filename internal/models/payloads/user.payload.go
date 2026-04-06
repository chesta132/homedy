package payloads

import "homedy/internal/models"

type RequestGetUser struct {
	ID string `uri:"id" validate:"required,uuid4"`
}

type ResponseGetUser struct {
	models.Base
	Username string `json:"username" example:"chesta_ardiona"`
	Email    string `json:"email" example:"chestaardi4@gmail.com"`
}
