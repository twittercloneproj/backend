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
		userRoles := []UserRole{nonregistereduser, user}
		result := db.Create(&userRoles)
		if result.Error != nil {
			return result.Error
		}

		permissions := []Permission{
			register, registerBusiness, login, getAllTweets, createTweet, getUserInfo, changePrivacy,
			getMyTweets, createRetweet, likeTweet, unlikeTweet, getTweetLikes, getMyHomeFeed, followUser, checkFollow,
			canIAccessTweet, acceptRejectRequest, unfollowUser, getMyFollowRequests, getMyFollowers, getUsersIFollow,
			getMySuggestions,
		}
		result = db.Create(&permissions)
		if result.Error != nil {
			return result.Error
		}

		rolePermissions := []RolePermission{
			nonregistereduserregister, nonregistereduserregisterBusiness, userlogin, usergetalltweets, usercreatetweet,
			usergetuserinfo, userchangeprivacy, usergetmytweets, usercreateretweet, userliketweet, userunliketweet,
			usergettweetlikes, usergetmyhomefeed, userfollowuser, usercheckfollow, usercaniaccesstweet, useracceptrejectrequest,
			userunfollowuser, usergetmyfollowrequests, usergetmyfollowers, usergetusersifollow, usergetmysuggestions,
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

var (
	register            = Permission{Id: "1", Name: "Register"}
	registerBusiness    = Permission{Id: "2", Name: "RegisterBusiness"}
	login               = Permission{Id: "3", Name: "Login"}
	getAllTweets        = Permission{Id: "4", Name: "GetAllTweets"}
	createTweet         = Permission{Id: "5", Name: "CreateTweet"}
	getUserInfo         = Permission{Id: "6", Name: "GetUserInfo"}
	changePrivacy       = Permission{Id: "7", Name: "ChangePrivacy"}
	getMyTweets         = Permission{Id: "8", Name: "GetMyTweets"}
	createRetweet       = Permission{Id: "9", Name: "CreateRetweet"}
	likeTweet           = Permission{Id: "10", Name: "LikeTweet"}
	unlikeTweet         = Permission{Id: "11", Name: "UnlikeTweet"}
	getTweetLikes       = Permission{Id: "12", Name: "GetTweetLikes"}
	getMyHomeFeed       = Permission{Id: "13", Name: "GetMyHomeFeed"}
	followUser          = Permission{Id: "14", Name: "FollowUser"}
	checkFollow         = Permission{Id: "15", Name: "CheckFollow"}
	canIAccessTweet     = Permission{Id: "16", Name: "CanIAccessTweet"}
	acceptRejectRequest = Permission{Id: "17", Name: "AcceptRejectRequest"}
	unfollowUser        = Permission{Id: "18", Name: "UnfollowUser"}
	getMyFollowRequests = Permission{Id: "19", Name: "GetMyFollowRequests"}
	getMyFollowers      = Permission{Id: "20", Name: "GetMyFollowers"}
	getUsersIFollow     = Permission{Id: "21", Name: "GetUsersIFollow"}
	getMySuggestions    = Permission{Id: "22", Name: "GetMySuggestions"}
)

var (
	nonregistereduserregister         = RolePermission{RoleId: nonregistereduser.Id, PermissionId: register.Id}
	nonregistereduserregisterBusiness = RolePermission{RoleId: nonregistereduser.Id, PermissionId: registerBusiness.Id}
	userlogin                         = RolePermission{RoleId: user.Id, PermissionId: login.Id}
	usergetalltweets                  = RolePermission{RoleId: user.Id, PermissionId: getAllTweets.Id}
	usercreatetweet                   = RolePermission{RoleId: user.Id, PermissionId: createTweet.Id}
	usergetuserinfo                   = RolePermission{RoleId: user.Id, PermissionId: getUserInfo.Id}
	userchangeprivacy                 = RolePermission{RoleId: user.Id, PermissionId: changePrivacy.Id}
	usergetmytweets                   = RolePermission{RoleId: user.Id, PermissionId: getMyTweets.Id}
	usercreateretweet                 = RolePermission{RoleId: user.Id, PermissionId: createRetweet.Id}
	userliketweet                     = RolePermission{RoleId: user.Id, PermissionId: likeTweet.Id}
	userunliketweet                   = RolePermission{RoleId: user.Id, PermissionId: unlikeTweet.Id}
	usergettweetlikes                 = RolePermission{RoleId: user.Id, PermissionId: getTweetLikes.Id}
	usergetmyhomefeed                 = RolePermission{RoleId: user.Id, PermissionId: getMyHomeFeed.Id}
	userfollowuser                    = RolePermission{RoleId: user.Id, PermissionId: followUser.Id}
	usercheckfollow                   = RolePermission{RoleId: user.Id, PermissionId: checkFollow.Id}
	usercaniaccesstweet               = RolePermission{RoleId: user.Id, PermissionId: canIAccessTweet.Id}
	useracceptrejectrequest           = RolePermission{RoleId: user.Id, PermissionId: acceptRejectRequest.Id}
	userunfollowuser                  = RolePermission{RoleId: user.Id, PermissionId: unfollowUser.Id}
	usergetmyfollowrequests           = RolePermission{RoleId: user.Id, PermissionId: getMyFollowRequests.Id}
	usergetmyfollowers                = RolePermission{RoleId: user.Id, PermissionId: getMyFollowers.Id}
	usergetusersifollow               = RolePermission{RoleId: user.Id, PermissionId: getUsersIFollow.Id}
	usergetmysuggestions              = RolePermission{RoleId: user.Id, PermissionId: getMySuggestions.Id}
)
