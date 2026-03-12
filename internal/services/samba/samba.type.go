package samba

type Bool string

const (
	YES Bool = "yes"
	NO  Bool = "no"
)

type Share struct {
	Path       string   `ini:"path"`
	ReadOnly   Bool     `ini:"read only"`
	Browsable  Bool     `ini:"browsable"`
	GuestUsers []string `ini:"guest users"`
	AdminUsers []string `ini:"admin users"`
}

type Shares map[string]Share
