package samba

type Bool string

const (
	YES Bool = "yes"
	NO  Bool = "no"
)

type Share struct {
	Path       string   `ini:"path,omitempty" json:"path"`
	ReadOnly   Bool     `ini:"read only,omitempty" json:"read_only"`
	Browsable  Bool     `ini:"browsable,omitempty" json:"browsable"`
	GuestUsers []string `ini:"guest users,omitempty" json:"guest_users"`
	AdminUsers []string `ini:"admin users,omitempty" json:"admin_users"`
}

type Shares map[string]Share
