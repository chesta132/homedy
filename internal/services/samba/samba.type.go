package samba

type Bool string

const (
	YES Bool = "yes"
	NO  Bool = "no"
)

var Bools = []Bool{YES, NO}

type Share struct {
	Path        string   `ini:"path,omitempty" json:"path" example:"/var/www/html" validate:"required,abs_path"`
	ReadOnly    Bool     `ini:"read only,omitempty" json:"read_only" example:"no" validate:"required,samba_bool"`
	Browsable   Bool     `ini:"browsable,omitempty" json:"browsable" example:"yes" validate:"required,samba_bool"`
	ValidUsers  []string `ini:"valid users,omitempty" json:"valid_users" example:"['guest']" validate:"required,min=1,dive,required"`
	AdminUsers  []string `ini:"admin users,omitempty" json:"admin_users" example:"['root']" validate:"required,min=1,dive,required"`
	Permissions []int    `ini:"homedy permission,omitempty" json:"permissions" example:"[7, 7, 7]" validate:"required,file_permission"`
}

type Shares map[string]Share

type ShareMap map[string]string
type ShareMaps map[string]ShareMap
