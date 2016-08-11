package db

import (
	"appengine"
	"appengine/datastore"
	"errors"
	// "fmt"
	// "os"
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

type Memo struct {
	Id       string
	FactorId string
	Memo     string
	Evidence string
	PhaseId  string
	Date     time.Time
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

// userlistKey returns the key used for all user entries.
func UserMemoKey(c appengine.Context, username string) *datastore.Key {
	return datastore.NewKey(c, "Memo", username, 0, nil)
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

func PutUser(c appengine.Context, u User, key *datastore.Key) (err error) {
	_, err = datastore.Put(c, key, &u)
	return
}

func GetUser(c appengine.Context, username string) (u User, k *datastore.Key, err error) {
	q := datastore.NewQuery("User").Ancestor(UserListKey(c)).
		Filter("Username=", username)

	var users []User
	// To retrieve the results,
	// you must execute the Query using its GetAll or Run methods.
	ks, err := q.GetAll(c, &users)

	if len(users) > 1 {
		err = errors.New("Error getting history: More than one User found!")
		return
	} else if len(users) != 0 {
		u = users[0]
		k = ks[0]
	}
	return
}

func GetHistory(c appengine.Context, username string) (messages []*Message, err error) {
	q := datastore.NewQuery("Message").Ancestor(UserHistoryKey(c, username)).Order("MessageNo").Limit(100)
	// [END query]
	// [START getall]
	messages = make([]*Message, 0, 100)
	_, err = q.GetAll(c, &messages)
	return
}

func GetHistoryCount(c appengine.Context, username string) (rc int, err error) {
	q := datastore.NewQuery("Message").Ancestor(UserHistoryKey(c, username)).Limit(100)
	rc, err = q.Count(c)
	return
}

func PutMemo(c appengine.Context, username string, m Memo) (err error) {
	memos := []Memo{m}
	var keys = []*datastore.Key{
		datastore.NewIncompleteKey(c, "Memo", UserMemoKey(c, username))}

	_, err = datastore.PutMulti(c, keys, memos)
	return
}
