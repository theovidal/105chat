package db

const (
	READ_MESSAGES      = 1
	WRITE_MESSAGES     = 2
	SEND_ANNOUNCEMENTS = 4

	MANAGE_MESSAGES = 8
	MANAGE_ROOM     = 16
	MANAGE_USERS    = 32
	MANAGE_ROLES    = 64
)

type Group struct {
	ID              uint          `json:"id" gorm:"primary_key"`
	Name            string        `json:"name"`
	Color           uint          `json:"color" valid:"range(0|16777215)"`
	Permissions     uint          `json:"permissions"`
	RoomPermissions map[uint]uint `json:"room_permissions" gorm:"-"`
	Inheritances    []uint        `json:"inheritances,omitempty" gorm:"-"`
}

func FetchPermissions(group *Group, id uint) {
	var groupToMerge Group
	Database.First(&groupToMerge, id)
	group.Permissions |= groupToMerge.Permissions

	AppendRoomPermissions(group, groupToMerge.ID)
	for _, inheritance := range FindGroupInheritances(id) {
		FetchPermissions(group, inheritance.ChildGroupID)
	}

	group.ID = id
}
