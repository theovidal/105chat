package http

import (
	"errors"
	"net/http"

	"github.com/asaskevich/govalidator"

	"github.com/theovidal/105chat/db"
	"github.com/theovidal/105chat/http/controllers"
	"github.com/theovidal/105chat/utils"
	"github.com/theovidal/105chat/ws"
)

// CreateGroup creates a new group with its own permissions
func CreateGroup(w http.ResponseWriter, r *http.Request) {
	if user := r.Context().Value("user").(*db.User); !user.HasGlobalPermission(db.MANAGE_GROUPS) {
		Response(w, http.StatusForbidden, nil)
		return
	}

	var payload GroupPayload
	if err := ParseBody(r, &payload); err != nil {
		Response(w, http.StatusBadRequest, nil)
		return
	}

	group := db.Group{
		ID:    utils.GenerateSnowflake(),
		Name:  payload.Name,
		Color: payload.Color,
	}
	err := db.Database.Create(&group).Error
	if err != nil {
		Response(w, http.StatusBadRequest, nil)
	}

	var errorsList []utils.Error
	errorsList = append(errorsList, controllers.UpdateGroupRoomPermissions(&group, payload.RoomPermissions)...)
	errorsList = append(errorsList, controllers.UpdateGroupInheritances(&group, payload.Inheritances)...)

	ws.Pipeline <- ws.Event{
		Event: ws.GROUP_CREATE,
		Data:  &group,
	}

	var status int
	if len(errorsList) > 0 {
		status = http.StatusPartialContent
	} else {
		status = http.StatusCreated
	}

	Response(w, status, utils.H{
		"id":     group.ID,
		"errors": errorsList,
	})
}

// GetGroups returns all the available groups
func GetGroups(w http.ResponseWriter, _ *http.Request) {
	var groups []db.Group
	db.Database.Find(&groups)

	for index, group := range groups {
		db.AppendRoomPermissions(&group, group.ID)
		db.AppendGroupInheritances(&group)
		groups[index] = group
	}

	Response(w, http.StatusOK, groups)
}

// GetGroup returns information about a specific group
func GetGroup(w http.ResponseWriter, r *http.Request) {
	group, err := ParseGroupFromURL(&w, r)
	if err != nil {
		return
	}

	db.AppendRoomPermissions(group, group.ID)
	db.AppendGroupInheritances(group)
	Response(w, http.StatusOK, group)
}

// UpdateGroup updates a group with its permissions, room permissions and inheritances
func UpdateGroup(w http.ResponseWriter, r *http.Request) {
	if user := r.Context().Value("user").(*db.User); !user.HasGlobalPermission(db.MANAGE_GROUPS) {
		Response(w, http.StatusForbidden, nil)
		return
	}

	group, err := ParseGroupFromURL(&w, r)
	if err != nil {
		return
	}
	db.AppendRoomPermissions(group, group.ID)
	db.AppendGroupInheritances(group)

	var payload GroupPayload
	if err = ParseBody(r, &payload); err != nil {
		Response(w, http.StatusBadRequest, nil)
		return
	}

	var errorsList []utils.Error
	errorsList = append(errorsList, controllers.UpdateGroupRoomPermissions(group, payload.RoomPermissions)...)
	errorsList = append(errorsList, controllers.UpdateGroupInheritances(group, payload.Inheritances)...)

	payload.Name = govalidator.Trim(payload.Name, "")
	if err = db.Database.Model(&group).Updates(payload).Error; err != nil {
		errorsList = append(errorsList, utils.Error{
			Key:     "dataError",
			Message: err.Error(),
		})
		return
	}

	ws.Pipeline <- ws.Event{
		Event: ws.GROUP_UPDATE,
		Data:  &group,
	}

	var status int
	if len(errorsList) > 0 {
		status = http.StatusPartialContent
	} else {
		status = http.StatusCreated
	}
	Response(w, status, utils.H{
		"id":     group.ID,
		"errors": errorsList,
	})
}

// DeleteGroup deletes a group and assigns its users to a fallback group
func DeleteGroup(w http.ResponseWriter, r *http.Request) {
	if user := r.Context().Value("user").(*db.User); !user.HasGlobalPermission(db.MANAGE_GROUPS) {
		Response(w, http.StatusForbidden, nil)
		return
	}

	group, err := ParseGroupFromURL(&w, r)
	if err != nil {
		return
	}

	var payload GroupDeletePayload
	if err = ParseBody(r, &payload); err != nil {
		Response(w, http.StatusBadRequest, nil)
		return
	}

	var count int
	db.Database.First(&db.Group{}, payload.FallbackGroupID).Count(&count)
	if count != 1 {
		Response(w, http.StatusBadRequest, nil)
		return
	}

	remainingData := utils.H{
		"group_id":          group.ID,
		"fallback_group_id": payload.FallbackGroupID,
	}

	ws.Pipeline <- ws.Event{
		Event: ws.GROUP_DELETE,
		Data:  &remainingData,
	}
	Response(w, http.StatusAccepted, &remainingData)

	var usersToUpdate []db.User
	db.Database.Where("group_id = ?", group.ID).Find(&usersToUpdate)
	for _, user := range usersToUpdate {
		user.GroupID = payload.FallbackGroupID
		db.Database.Save(&user)
	}
}

// ParseGroupFromURL checks for errors in the passed group ID inside request's URL
func ParseGroupFromURL(w *http.ResponseWriter, r *http.Request) (group *db.Group, err error) {
	group, err = controllers.FindGroupFromURL(r)
	if errors.Is(err, controllers.InvalidType) {
		Response(*w, http.StatusBadRequest, nil)
	} else if errors.Is(err, controllers.UnknownRoom) {
		Response(*w, http.StatusNotFound, nil)
	}
	return
}
