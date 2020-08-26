package db

// RoomPermission defines permissions of a group for a specific room
type RoomPermission struct {
	// Identifier of the room permission, Twitter snowflake
	ID uint `json:"-" gorm:"primary_key,AUTO_INCREMENT"`
	// Identifier of the concerned room, Twitter snowflake
	RoomID uint `json:"room_id"`
	// Identifier of the concerned group, Twitter snowflake
	GroupID uint `json:"group_id"`
	// Permissions the group has in this room
	Permissions uint `json:"permissions"`
}

// AppendRoomPermissions appends all the room permissions of a group to its structure
func AppendRoomPermissions(group *Group, id uint) {
	if group.RoomPermissions == nil {
		group.RoomPermissions = make(map[uint]uint)
	}

	var roomPermissions []RoomPermission
	Database.Where("group_id = ?", id).Find(&roomPermissions)
	for _, room := range roomPermissions {
		group.RoomPermissions[room.RoomID] |= room.Permissions
	}
}
