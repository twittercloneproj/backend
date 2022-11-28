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
	BusinessUser      = "BusinessUser"
)

var (
	nonregistereduser = UserRole{Id: "", Name: NonregisteredUser}
	user              = UserRole{Id: "", Name: User}
	businessuser      = UserRole{Id: "", Name: BusinessUser}
)
