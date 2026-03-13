package payloads

import "homedy/internal/services/samba"

type RequestAddShare struct {
	TemplateShareName
	samba.Share
}

func (r *RequestAddShare) ToShare() samba.Share {
	return r.Share
}

type RequestUpdateShare struct {
	TemplateShareName
	samba.Share
}

type TemplateShareName struct {
	Name string `json:"name" uri:"name" example:"apache_source" validate:"required,share_name"`
}
