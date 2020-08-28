package db

import "github.com/theovidal/105chat/cache"

// User represents a user of 105chat, who communicates to others via messages in rooms.
type User struct {
	// Identifier of the user, Twitter snowflake
	ID uint `json:"id" gorm:"primary_key"`
	// Name of the user, 2~32 characters
	Name string `json:"name" gorm:"size:32" valid:"required,length(2|32)"`
	// URL pointing to user's avatar
	AvatarURL string `json:"avatar_url" valid:"url"`
	// Description of the user, 0~512 characters
	Description string `json:"description" gorm:"size:512" valid:"length(0|512)"`
	// Color of the user (stored as an integer for less complexity)
	Color uint `json:"color" valid:"range(0|16777215)"`
	// When the user was created (via registration or administrator action)
	Timestamp int64 `json:"timestamp"`
	// Identifier of user's group, Twitter snowflake
	GroupID uint `json:"group_id"`
	// Entity of user's group, used in the code for a better simplicity
	Group Group `json:"-" gorm:"-"`
	// Whether the user is muted by a moderator
	Muted bool `json:"muted"`
	// Whether the user is disabled by a moderator
	Disabled bool `json:"disabled"`

	// Email of the user, used to communicate and authenticate
	Email string `json:"-" gorm:"unique" valid:"required,email"`
	// Password of the user, used to authenticate
	Password string `json:"-" valid:"required"`
	// Token of the user, used to interact with the API (WS and HTTP)
	// Is obtained after user's login via email and password
	Token string `json:"-"`
}

// FindUserByToken searches for a user with a specific token
func FindUserByToken(token string) (user User, err error) {
	err = Client.Where("token = ?", token).First(&user).Error
	return
}

func (u *User) HasGlobalPermission(permission uint) bool {
	return cache.GetGroupPermissions(u.GroupID)&permission != 0
}

func (u *User) HasRoomPermission(room uint, permission uint) bool {
	return cache.GetGroupRoomPermissions(u.GroupID, room)&permission != 0
}

func (u *User) HasAnyPermission(room uint, permission uint) bool {
	return u.HasGlobalPermission(permission) || u.HasRoomPermission(room, permission)
}
