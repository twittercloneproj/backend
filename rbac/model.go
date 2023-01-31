package rbac

type UserRole struct {
	Id   string `gorm:"primaryKey"`
	Name string `gorm:"unique"`
}

type Permission struct {
	Id   string `gorm:"primaryKey"`
	Name string `gorm:"unique"`
}

type RolePermission struct {
	RoleId       string `gorm:"primaryKey"`
	PermissionId string `gorm:"primaryKey"`
}

const (
	NonregisteredUser = "NonregisteredUser"
	User              = "User"
)

var (
	nonregistereduser = UserRole{Id: "1", Name: NonregisteredUser}
	user              = UserRole{Id: "2", Name: User}
)
