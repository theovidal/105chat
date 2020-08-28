package http

import (
	"errors"
	"github.com/asaskevich/govalidator"
	"net/http"

	"github.com/theovidal/105chat/db"
	"github.com/theovidal/105chat/http/controllers"
	"github.com/theovidal/105chat/ws"
)

// GetUser returns information about a specific user thanks to their ID
func GetUser(w http.ResponseWriter, r *http.Request) {
	user, err := ParseUserFromURL(&w, r)
	if err != nil {
		return
	}

	Response(w, http.StatusOK, user)
}

// UpdateUserProfile is used by a user to edit their profile
func UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	userToUpdate, err := ParseUserFromURL(&w, r)
	if err != nil {
		return
	}
	authenticatedUser := r.Context().Value("user").(*db.User)
	if userToUpdate.ID != authenticatedUser.ID && !authenticatedUser.HasGlobalPermission(db.MANAGE_USERS) {
		Response(w, http.StatusForbidden, nil)
		return
	}

	var payload UserProfileUpdatePayload
	if err = ParseBody(r, &payload); err != nil {
		Response(w, http.StatusBadRequest, nil)
		return
	}

	payload.Name = govalidator.Trim(payload.Name, "")
	payload.AvatarURL = govalidator.Trim(payload.AvatarURL, "")
	payload.Description = govalidator.Trim(payload.Description, "")
	if err = db.Database.Model(userToUpdate).Updates(payload).Error; err != nil {
		Response(w, http.StatusBadRequest, nil)
		return
	}

	ws.Pipeline <- ws.Event{
		Event: ws.USER_UPDATE,
		Data:  &userToUpdate,
	}
	ws.Station.EditUser(userToUpdate)
	Response(w, http.StatusNoContent, nil)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	userToUpdate, err := ParseUserFromURL(&w, r)
	if err != nil {
		return
	}
	authenticatedUser := r.Context().Value("user").(*db.User)
	if !authenticatedUser.HasGlobalPermission(db.MANAGE_USERS) {
		Response(w, http.StatusForbidden, nil)
		return
	}

	var payload UserUpdatePayload
	if err = ParseBody(r, &payload); err != nil {
		Response(w, http.StatusBadRequest, nil)
		return
	}

	userToUpdate.Muted = payload.Muted
	userToUpdate.Disabled = payload.Disabled
	db.Database.Model(&userToUpdate).Updates(payload)

	ws.Pipeline <- ws.Event{
		Event: ws.USER_UPDATE,
		Data:  &userToUpdate,
	}
	ws.Station.EditUser(userToUpdate)
	Response(w, http.StatusNoContent, nil)
}

func GetUserGroup(w http.ResponseWriter, r *http.Request) {
	userToFetch, err := ParseUserFromURL(&w, r)
	if err != nil {
		return
	}
	authenticatedUser := r.Context().Value("user").(*db.User)

	if userToFetch.ID != authenticatedUser.ID && !authenticatedUser.HasGlobalPermission(db.MANAGE_USERS) {
		Response(w, http.StatusForbidden, nil)
		return
	}

	var group db.Group
	db.Database.Find(&group, userToFetch.GroupID)
	db.FetchPermissions(&group, userToFetch.GroupID)

	Response(w, http.StatusOK, &group)
}

func UpdateUserGroup(w http.ResponseWriter, r *http.Request) {
	userToUpdate, err := ParseUserFromURL(&w, r)
	if err != nil {
		return
	}
	authenticatedUser := r.Context().Value("user").(*db.User)

	if !authenticatedUser.HasGlobalPermission(db.MANAGE_USERS) {
		Response(w, http.StatusForbidden, nil)
		return
	}

	var payload UserGroupUpdatePayload
	if err = ParseBody(r, &payload); err != nil {
		Response(w, http.StatusBadRequest, nil)
		return
	}

	var group db.Group
	if err = db.Database.Find(&group, payload.GroupID).Error; err != nil {
		Response(w, http.StatusNotFound, nil)
		return
	}

	userToUpdate.GroupID = group.ID
	db.Database.Save(&userToUpdate)

	ws.Pipeline <- ws.Event{
		Event: ws.USER_UPDATE,
		Data:  &userToUpdate,
	}
	ws.Station.EditUser(userToUpdate)
	Response(w, http.StatusNoContent, nil)
}

// ParseUserFromURL checks for errors in the passed user ID inside request's URL
func ParseUserFromURL(w *http.ResponseWriter, r *http.Request) (user *db.User, err error) {
	user, err = controllers.FindUserFromURL(r)
	if errors.Is(err, controllers.InvalidType) {
		Response(*w, http.StatusBadRequest, nil)
	} else if errors.Is(err, controllers.UnknownUser) {
		Response(*w, http.StatusNotFound, nil)
	}
	return
}
