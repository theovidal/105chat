package controllers

import (
	"fmt"
	"net/http"

	"github.com/theovidal/105chat/db"
	"github.com/theovidal/105chat/utils"
)

// FindGroupFromURL parses request's URL and find the corresponding group thanks to the ID
func FindGroupFromURL(r *http.Request) (*db.Group, error) {
	groupID, err := FindIDFromURL(r, "group")
	if err != nil {
		return &db.Group{}, err
	}

	var group db.Group
	if err = db.Client.First(&group, groupID).Error; err != nil {
		return &db.Group{}, UnknownRoom
	}

	return &group, nil
}

// UpdateGroupRoomPermissions updates room permissions of a group inside the database
func UpdateGroupRoomPermissions(group *db.Group, roomPermissions map[uint]uint) (errorsList []utils.Error) {
	for roomID, newPermissions := range roomPermissions {
		var count int
		db.Client.First(&db.Room{}, roomID).Count(&count)
		if count != 1 {
			errorsList = append(errorsList, utils.Error{
				Key:     fmt.Sprint("roomPermissions.unknownRoom.", roomID),
				Message: fmt.Sprint("Unknown room with ID ", roomID),
			})
			continue
		}

		actualPermissions, exists := group.RoomPermissions[roomID]
		if exists {
			if actualPermissions != newPermissions {
				var roomPermission db.RoomPermission
				db.Client.Where("group_id = ? AND room_id = ?", group.ID, roomID).Find(&roomPermission)
				roomPermission.Permissions = newPermissions
				db.Client.Save(&roomPermission)
			}
		} else {
			roomPermission := db.RoomPermission{
				ID:          utils.GenerateSnowflake(),
				RoomID:      roomID,
				GroupID:     group.ID,
				Permissions: newPermissions,
			}
			db.Client.Create(&roomPermission)
		}
		delete(group.RoomPermissions, roomID)
	}
	// Reminding permissions are to delete
	for roomID := range group.RoomPermissions {
		db.Client.Where(
			"group_id = ? AND room_id = ?",
			group.ID,
			roomID,
		).Delete(&db.RoomPermission{})
	}

	return
}

// UpdateGroupInheritances updates inheritances of a group inside the database
func UpdateGroupInheritances(group *db.Group, inheritances []uint) (errorsList []utils.Error) {
	for _, childGroupID := range inheritances {
		if childGroupID == group.ID {
			errorsList = append(errorsList, utils.Error{
				Key:     "inheritance.selfInherit",
				Message: "You can't make a group inherit from itself",
			})
			continue
		}

		var count int
		if db.Client.First(&db.Group{}, childGroupID).Count(&count); count != 1 {
			errorsList = append(errorsList, utils.Error{
				Key:     fmt.Sprint("inheritance.unknownGroup.", childGroupID),
				Message: fmt.Sprint("Unknown group with ID ", childGroupID),
			})
			continue
		}

		if !utils.Contains(group.Inheritances, childGroupID) {
			inheritance := db.GroupInheritance{
				ParentGroupID: group.ID,
				ChildGroupID:  childGroupID,
			}
			db.Client.Create(&inheritance)
		}
	}
	for _, childGroupID := range group.Inheritances {
		if !utils.Contains(inheritances, childGroupID) {
			db.Client.Where(
				"parent_group_id = ? AND child_group_id = ?",
				group.ID,
				childGroupID,
			).Delete(&db.GroupInheritance{})
		}
	}

	return
}
