package cache

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/fatih/color"
	"github.com/go-redis/redis/v8"
)

func SetGroupPermissions(group uint, permissions uint) {
	err := Client.Set(context.Background(), GroupPermissionsKey(group), permissions, 0).Err()
	if err != nil {
		log.Println(color.RedString("‼ Error while setting cache for group permissions:", err.Error()))
	}
}

func SetGroupRoomPermissions(group uint, room uint, permissions uint) {
	err := Client.HSet(context.Background(), GroupRoomPermissionsKey(group), room, permissions).Err()
	if err != nil {
		log.Println(color.RedString("‼ Error while setting cache for group room permissions:", err.Error()))
	}
}

func GetGroupPermissions(group uint) uint {
	result, err := Client.Get(context.Background(), GroupPermissionsKey(group)).Result()
	if err == redis.Nil {
		result = "0"
	}
	permissions, _ := strconv.Atoi(result)
	fmt.Println(permissions)
	return uint(permissions)
}

func GetAllGroupRoomPermissions(group uint) map[uint]uint {
	raw, err := Client.HGetAll(context.Background(), GroupRoomPermissionsKey(group)).Result()
	if err == redis.Nil {
		return map[uint]uint{}
	}

	data := make(map[uint]uint)
	for room, permission := range raw {
		roomID, _ := strconv.Atoi(room)
		permissionID, _ := strconv.Atoi(permission)
		data[uint(roomID)] = uint(permissionID)
	}
	return data
}

func GetGroupRoomPermissions(group uint, room uint) uint {
	result, _ := Client.HGet(context.Background(), GroupRoomPermissionsKey(group), strconv.Itoa(int(room))).Result()
	permissions, _ := strconv.Atoi(result)
	return uint(permissions)
}

func RemoveGroupCache(group uint) {
	_ = Client.Del(context.Background(), GroupPermissionsKey(group))
	_ = Client.Del(context.Background(), GroupRoomPermissionsKey(group))
}

func GroupPermissionsKey(group uint) string {
	return fmt.Sprintf("group-permissions-%d", group)
}

func GroupRoomPermissionsKey(group uint) string {
	return fmt.Sprintf("group-room-permissions-%d", group)
}
