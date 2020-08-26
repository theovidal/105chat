package http

// AuthenticatePayload is sent to authenticate to the API and get back a token
type AuthenticatePayload struct {
	Email    string
	Password string
}

// UserProfileUpdatePayload is sent to update user's profile
type UserProfileUpdatePayload struct {
	Name        string
	AvatarURL   string `json:"avatar_url"`
	Description string
	Color       uint
}

// UserUpdatePayload is sent by a user with the MANAGE_USERS permission to update a user
type UserUpdatePayload struct {
	Muted    bool
	Disabled bool
}

// UserGroupUpdatePayload is sent to update the group in which a user is
type UserGroupUpdatePayload struct {
	GroupID uint `json:"group_id" valid:"required"`
}

// GroupPayload is sent to create or update a group and its permissions, room permissions and inheritances
type GroupPayload struct {
	Name            string
	Color           uint
	Permissions     uint
	RoomPermissions map[uint]uint `json:"room_permissions" gorm:"-"`
	Inheritances    []uint        `json:"inheritances" gorm:"-"`
}

// GroupDeletePayload is sent to delete a group a assign another group to users
type GroupDeletePayload struct {
	FallbackGroupID uint `json:"fallback_group_id"`
}

// MessageCreatePayload is sent to create a message in a room
type MessageCreatePayload struct {
	Content      string
	Announcement bool
}

// MessageUpdatePayload is sent to update a message in a room
type MessageUpdatePayload struct {
	Content string
}
