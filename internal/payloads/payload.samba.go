package payloads

import "homedy/internal/services/samba"

type RequestAddShare struct {
	Name string `json:"name" example:"apache_source"`
	samba.Share
}

func (r *RequestAddShare) ToShare() samba.Share {
	return r.Share
}
