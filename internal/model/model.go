package model

// User
type User struct {
	Id       int64  `gorm:"column:id;primary_key;AUTO_INCREMENT;NOT NULL"`
	Name     string `gorm:"column:name;default:NULL"`
	Accout   string `gorm:"column:accout;default:NULL"`
	Password string `gorm:"column:password;default:NULL"`
}

// TableName 表名
func (u *User) TableName() string {
	return "user"
}

// Role
type Role struct {
	Id       int64  `gorm:"column:id;primary_key;AUTO_INCREMENT;NOT NULL"`
	RoleName string `gorm:"column:role_name;default:NULL;comment:'角色名称'"`
}

// TableName 表名
func (r *Role) TableName() string {
	return "role"
}

// UserRole
type UserRole struct {
	Id     int64 `gorm:"column:id;primary_key;AUTO_INCREMENT;NOT NULL"`
	UserId int64 `gorm:"column:user_id;default:NULL;comment:'用户id'"`
	RoleId int64 `gorm:"column:role_id;default:NULL;comment:'角色id'"`
}

// TableName 表名
func (u *UserRole) TableName() string {
	return "user_role"
}

// Permission
type Permission struct {
	Id       int64  `gorm:"column:id;primary_key;AUTO_INCREMENT;NOT NULL"`
	PermName string `gorm:"column:perm_name;default:NULL"`
	PermCode string `gorm:"column:perm_code;default:NULL"`

	ApiPath string `gorm:"column:api_path;default:NULL"` // 例如：/api/users
	Method  string `gorm:"column:method;default:NULL"`
}

// TableName 表名
func (p *Permission) TableName() string {
	return "permission"
}

// RolePerm
type RolePerm struct {
	Id     int64 `gorm:"column:id;primary_key;AUTO_INCREMENT;NOT NULL"`
	RoleId int64 `gorm:"column:role_id;default:NULL"`
	PermId int64 `gorm:"column:perm_id;default:NULL"`
}

// TableName 表名
func (r *RolePerm) TableName() string {
	return "role_perm"
}
