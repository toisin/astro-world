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
	Username        string
	Screenname      string
	Date            time.Time // Time when user is created
	CurrentPhaseId  string
	CurrentPromptId string
	CurrentFactorId string
}

// Generic message with no additional phase specific details
type Message struct {
	// Username string // Does not really need to store this everytime
	Id        string
	Text      string
	Mtype     string // ROBOT | HUMAN
	Date      time.Time
	MessageNo int // in the order of the message
}

type CovMessage struct {
	MessageId   string // The message this cov message is linked to
	FactorId    string // not empty if message is related to an investigating factor
	RecordNoOne string
	RecordNoTwo string
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
