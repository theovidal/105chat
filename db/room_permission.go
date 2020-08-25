package db

type RoomPermission struct {
	RoomID      uint `json:"room_id"`
	GroupID     uint `json:"group_id"`
	Permissions uint `json:"permissions"`
}

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
