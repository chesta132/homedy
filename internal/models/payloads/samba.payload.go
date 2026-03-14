package payloads

import "homedy/internal/models"

type RequestCreateShare struct {
	TemplateShareName
	models.Share
}

type RequestUpdateShare struct {
	TemplateShareName
	models.Share
}

type RequestDeleteShare struct {
	TemplateShareName
}

type TemplateShareName struct {
	Name string `json:"name" uri:"name" example:"apache_source" validate:"required,share_name"`
}
