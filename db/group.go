package db

import "github.com/theovidal/105chat/cache"

// Constants corresponding to bitwise permissions
const (
	READ_MESSAGES      = 1
	WRITE_MESSAGES     = 2
	SEND_ANNOUNCEMENTS = 4

	MANAGE_MESSAGES = 8
	MANAGE_ROOM     = 16
	MANAGE_USERS    = 32
	MANAGE_GROUPS   = 64
)

// Group defines a group of users with certain permissions
type Group struct {
	// Identifier of the group, Twitter snowflake
	ID uint `json:"id" gorm:"primary_key"`
	// Name of the group
	Name string `json:"name" gorm:"size:32" valid:"length(2|32)"`
	// Color of the group (hex string)
	Color string `json:"color" valid:"hexcolor"`
	// Global permissions the group has
	Permissions uint `json:"permissions"`
	// Permissions specific to rooms (used internally)
	RoomPermissions map[uint]uint `json:"room_permissions" gorm:"-"`
	// All the groups this group inherit from (used internally)
	Inheritances []uint `json:"inheritances,omitempty" gorm:"-"`
}

// FetchPermissions fetches all the permissions of the group by descending inheritance
func FetchPermissions(group *Group, id uint) {
	var groupToMerge Group
	Client.First(&groupToMerge, id)
	group.Permissions |= groupToMerge.Permissions

	AppendRoomPermissions(group, groupToMerge.ID)
	for _, inheritance := range FindGroupInheritances(id) {
		FetchPermissions(group, inheritance.ChildGroupID)
	}

	group.ID = id
}

func SetAllGroupsCache() {
	var groups []Group
	Client.Find(&groups)

	for _, group := range groups {
		SetGroupCache(&group)
	}
}

func SetGroupCache(group *Group) {
	FetchPermissions(group, group.ID)
	cache.SetGroupPermissions(group.ID, group.Permissions)
	for room, permissions := range group.RoomPermissions {
		cache.SetGroupRoomPermissions(group.ID, room, permissions)
	}
}
