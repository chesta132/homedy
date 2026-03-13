package samba

type Bool string

const (
	YES Bool = "yes"
	NO  Bool = "no"
)

type Share struct {
	Path        string   `ini:"path,omitempty" json:"path" example:"/var/www/html"`
	ReadOnly    Bool     `ini:"read only,omitempty" json:"read_only" example:"no"`
	Browsable   Bool     `ini:"browsable,omitempty" json:"browsable" example:"yes"`
	GuestUsers  []string `ini:"guest users,omitempty" json:"guest_users" example:"['guest']"`
	AdminUsers  []string `ini:"admin users,omitempty" json:"admin_users" example:"['root']"`
	Permissions []int    `ini:"-" json:"permissions" example:"[7, 7, 7]"`
}

type Shares map[string]Share
