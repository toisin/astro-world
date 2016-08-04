package db

import (
	"appengine"
	"appengine/datastore"
	"time"
)

const (
	// Message Type: Constants for Message.Mtype
	ROBOT = "robot"
	HUMAN = "student"
)

type User struct {
	Username             string
	Screenname           string
	Date                 time.Time // Time when user is created
	CurrentPhaseId       string
	CurrentSequenceOrder int
	CurrentPromptId      string
	CurrentFactorId      string
	UIState              []byte // Do not store as string because string type has a limit of 500 characters
	ArchiveHistoryLength int
}

// Generic message with no additional phase specific details
type Message struct {
	// Username string // Does not really need to store this everytime
	Id        string
	Texts     []string
	Mtype     string // ROBOT | HUMAN
	Date      time.Time
	MessageNo int // in the order of the message
}

// userlistKey returns the key used for all user entries.
func UserLogsKey(c appengine.Context, username string) *datastore.Key {
	return datastore.NewKey(c, "Logs", username, 0, nil)
}

// userlistKey returns the key used for all user entries.
func UserHistoryKey(c appengine.Context, username string) *datastore.Key {
	return datastore.NewKey(c, "History", username, 0, nil)
}

// userKey returns the key used for all user entries.
func UserKey(c appengine.Context) *datastore.Key {
	return datastore.NewIncompleteKey(c, "User", UserListKey(c))
}

// userListKey returns the key used as the ancestor for all user entries.
func UserListKey(c appengine.Context) *datastore.Key {
	// The string "default_guestbook" here could be varied to have multiple guestbooks.
	return datastore.NewKey(c, "UserList", "default_userlist", 0, nil)
}
