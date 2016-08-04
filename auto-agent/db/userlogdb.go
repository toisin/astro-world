package db

import (
	// "appengine"
	// "appengine/datastore"
	"time"
)

// Generic message with no additional phase specific details
type UserLog struct {
	Username     string // Storing this in case if Users are deleted for some reasons
	Id           string
	PromptId     string
	PhaseId      string
	QuestionText string
	JsonResponse string
	ResponseId   string
	ResponseText string
	Mtype        string // ROBOT | HUMAN
	Date         time.Time
	URL          string
}
