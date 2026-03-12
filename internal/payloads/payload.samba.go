package payloads

import "homedy/internal/services/samba"

type RequestAddShare struct {
	Name       string     `json:"name" example:"apache_source"`
	Path       string     `json:"path" example:"/var/www/html"`
	ReadOnly   samba.Bool `json:"read_only" example:"no"`
	Browsable  samba.Bool `json:"browsable" example:"yes"`
	GuestUsers []string   `json:"guest_users" example:"['guest']"`
	AdminUsers []string   `json:"admin_users" example:"['root']"`
	Permission []int      `json:"permission" example:"[7, 7, 7]"`
}
