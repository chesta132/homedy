package samba

type Bool string

const (
	YES Bool = "yes"
	NO  Bool = "no"
)

type Share struct {
	Path       string   `ini:"path" json:"path"`
	ReadOnly   Bool     `ini:"read only" json:"read_only"`
	Browsable  Bool     `ini:"browsable" json:"browsable"`
	GuestUsers []string `ini:"guest users" json:"guest_users"`
	AdminUsers []string `ini:"admin users" json:"admin_users"`
}

type Shares map[string]Share
