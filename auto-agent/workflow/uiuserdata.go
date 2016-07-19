package workflow

import (
	"db"
)

// Includes only the variables that are needed on the client side
type UIUserData struct {
	User            *db.User
	History         []db.Message
	CurrentUIPrompt UIPrompt
	CurrentUIAction UIAction
	State           StateEntities
}

type StateEntities interface {
	GetPhaseId() string
}
