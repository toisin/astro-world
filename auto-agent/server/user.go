package server

import (
	"db"
	"workflow"
	// "fmt"
	// "os"
)

type UserData struct {
	uiUserData UIUserData
	CurrentPrompt workflow.Prompt
}

func MakeUserData(u *db.User) *UserData {
  // Process submitted answer
	ud := UserData {}
	ud.uiUserData.User = u
	if ((u.CurrentPromptId == "") || (u.CurrentPhaseId == "")) {
		ud.CurrentPrompt = workflow.MakeFirstPrompt()
		ud.uiUserData.User.CurrentPhaseId = ud.CurrentPrompt.GetPhaseId()
		ud.uiUserData.User.CurrentPromptId = ud.CurrentPrompt.GetPromptId()
		ud.uiUserData.CurrentUIPrompt = ud.CurrentPrompt.GetUIPrompt()
	} else {
		phaseId := u.CurrentPhaseId
		promptId := u.CurrentPromptId
		ud.CurrentPrompt = workflow.MakePrompt(promptId, phaseId)
		ud.uiUserData.CurrentUIPrompt = ud.CurrentPrompt.GetUIPrompt()
	}

	return &ud
}

func (ud *UserData)GetUIUserData() *UIUserData {
	return &ud.uiUserData
}

func (ud *UserData)UpdateWithNextPrompt() {
	ud.CurrentPrompt = ud.CurrentPrompt.GetNextPrompt()
	
  if (ud.CurrentPrompt != nil) {
		ud.uiUserData.CurrentUIPrompt = ud.CurrentPrompt.GetUIPrompt()
		ud.uiUserData.User.CurrentPhaseId = ud.CurrentPrompt.GetPhaseId()
		ud.uiUserData.User.CurrentPromptId = ud.CurrentPrompt.GetPromptId()
	}
}

// Includes only the variables that are needed on the client side
type UIUserData struct {
	User *db.User
	History []db.Message
	CurrentUIPrompt workflow.UIPrompt
}



