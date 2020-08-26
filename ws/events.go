package ws

// Event is the contained for all events that'll be sent inside the pipeline
type Event struct {
	// The type of the event, as defined in the constants below
	Event string `json:"event"`
	// The data related to the event (message, rooom, user...)
	Data interface{} `json:"data,omitempty"`
}

// API to Client events
const (
	AUTHENTICATION_NEEDED  = "AUTHENTICATION_NEEDED"
	AUTHENTICATION_SUCCESS = "AUTHENTICATION_SUCCESS"
	AUTHENTICATION_FAIL    = "AUTHENTICATION_FAIL"
	PONG                   = "PONG"
	CLOSE                  = "CLOSE"
	ERROR                  = "ERROR"

	ROOM_CREATE = "ROOM_CREATE"
	ROOM_UPDATE = "ROOM_UPDATE"
	ROOM_DELETE = "ROOM_DELETE"

	MESSAGE_CREATE = "MESSAGE_CREATE"
	MESSAGE_UPDATE = "MESSAGE_UPDATE"
	MESSAGE_DELETE = "MESSAGE_DELETE"

	USER_UPDATE         = "USER_UPDATE"
	USER_PROFILE_UPDATE = "USER_PROFILE_UPDATE"
	USER_CONNECT        = "USER_CONNECT"
	USER_DISCONNECT     = "USER_DISCONNECT"
)

// Error represents an error in event processing
type Error struct {
	// The code of the error (4xx for client errors, 5xx for server errors)
	Code int `json:"code"`
	// The message that explains the error
	Message string `json:"message"`
}

// Error codes sent to the client
var (
	ERROR400 = Error{400, "Unknown client error"}
	ERROR401 = Error{401, "You must authenticate in order to send and receive events"}
	ERROR402 = Error{402, "Invalid event data"}
	ERROR403 = Error{403, "You don't have the permission to send this event"}
	ERROR404 = Error{404, "Unknown event"}
	ERROR405 = Error{405, "A moderator disabled your account. Contact the moderation team to know further details"}

	ERROR500 = Error{400, "Unknown server error"}
)
