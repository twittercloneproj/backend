package rbac

import (
	"gorm.io/gorm"
)

func SetupContentRBAC(db *gorm.DB) error {
	dropContentTables(db)
	err := db.AutoMigrate(&UserRole{}, Permission{}, RolePermission{})
	if err != nil {
		return err
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		userRoles := []UserRole{nonregistereduser, user, businessuser}
		result := db.Create(&userRoles)
		if result.Error != nil {
			return result.Error
		}

		permissions := []Permission{
			register, registerBusiness, login, getAllTweets, createTweet,
		}
		result = db.Create(&permissions)
		if result.Error != nil {
			return result.Error
		}

		rolePermissions := []RolePermission{
			nonregistereduserregister, nonregistereduserregisterBusiness, userlogin, businessuserlogin,
			usergetAllTweets, userCreateTweet, businessusergetAllTweets, businessuserCreateTweet,
		}

		result = db.Create(&rolePermissions)
		if result.Error != nil {
			return result.Error
		}

		return nil
	})

	return err
}

func dropContentTables(db *gorm.DB) {
	if db.Migrator().HasTable(&UserRole{}) {
		db.Migrator().DropTable(&UserRole{})
	}
	if db.Migrator().HasTable(&Permission{}) {
		db.Migrator().DropTable(&Permission{})
	}
	if db.Migrator().HasTable(&RolePermission{}) {
		db.Migrator().DropTable(&RolePermission{})
	}
}

// Content RBAC
var (
	register         = Permission{Id: "", Name: "Register"}
	registerBusiness = Permission{Id: "", Name: "RegisterBusiness"}
	login            = Permission{Id: "", Name: "Login"}
	getAllTweets     = Permission{Id: "", Name: "GetAllTweets"}
	createTweet      = Permission{Id: "", Name: "CreateTweet"}
)

var (
	// Posts
	nonregistereduserregister         = RolePermission{RoleId: nonregistereduser.Id, PermissionId: register.Id}
	nonregistereduserregisterBusiness = RolePermission{RoleId: nonregistereduser.Id, PermissionId: registerBusiness.Id}
	userlogin                         = RolePermission{RoleId: user.Id, PermissionId: login.Id}
	businessuserlogin                 = RolePermission{RoleId: businessuser.Id, PermissionId: login.Id}
	userCreateTweet                   = RolePermission{RoleId: user.Id, PermissionId: createTweet.Id}
	businessuserCreateTweet           = RolePermission{RoleId: businessuser.Id, PermissionId: createTweet.Id}
	usergetAllTweets                  = RolePermission{RoleId: user.Id, PermissionId: getAllTweets.Id}
	businessusergetAllTweets          = RolePermission{RoleId: businessuser.Id, PermissionId: getAllTweets.Id}
)
