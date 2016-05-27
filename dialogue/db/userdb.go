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
	Username string
	Screenname string
    Date time.Time
    CurrentWorkflowStateId string
	CurrentX string
	Case1 string
	Case2 string
	// CompletedX []string // What is that?
}

type Message struct {
	// Username string // Does not really need to store this everytime
	Value string
	Text string
	Mtype string // ROBOT | HUMAN
	WorflowStateID string
    Date time.Time
    RecordNo int
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
