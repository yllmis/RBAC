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
