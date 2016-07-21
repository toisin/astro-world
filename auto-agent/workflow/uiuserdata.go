package workflow

import (
	"db"
)

// Includes only the variables that are needed on the client side
type UIUserData struct {
	Username        string
	Screenname      string
	CurrentPhaseId  string
	CurrentFactorId string
	History         []db.Message
	CurrentUIPrompt UIPrompt
	CurrentUIAction UIAction
	State           StateEntities
}

type StateEntities interface {
	setUsername(string)
	setScreenname(string)
	setTargetFactor(*FactorState)
	GetPhaseId() string
}
