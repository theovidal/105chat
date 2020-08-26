package http

type AuthenticatePayload struct {
	Email    string
	Password string
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

// UserProfileUpdatePayload is sent to update user's profile
type UserProfileUpdatePayload struct {
	Name        string
	AvatarURL   string `json:"avatar_url"`
	Description string
	Color       uint
}
