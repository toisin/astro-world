package server

import (
	"db"
	// "fmt"
	// "os"
	"workflow"
)

type UserData struct {
	uiUserData    workflow.UIUserData
	CurrentPrompt workflow.Prompt
}

func MakeUserData(u *db.User) *UserData {
	// Process submitted answer
	ud := UserData{}
	ud.uiUserData.User = u
	if (u.CurrentPromptId == "") || (u.CurrentPhaseId == "") {
		ud.CurrentPrompt = workflow.MakeFirstPrompt()
		ud.uiUserData.User.CurrentPhaseId = ud.CurrentPrompt.GetPhaseId()
		ud.uiUserData.User.CurrentPromptId = ud.CurrentPrompt.GetPromptId()
		ud.uiUserData.CurrentUIPrompt = ud.CurrentPrompt.GetUIPrompt(&ud.uiUserData)
		ud.uiUserData.CurrentUIAction = ud.CurrentPrompt.GetUIAction()
	} else {
		phaseId := u.CurrentPhaseId
		promptId := u.CurrentPromptId
		ud.CurrentPrompt = workflow.MakePrompt(promptId, phaseId)
		ud.uiUserData.CurrentUIPrompt = ud.CurrentPrompt.GetUIPrompt(&ud.uiUserData)
		ud.uiUserData.CurrentUIAction = ud.CurrentPrompt.GetUIAction()
	}

	return &ud
}

func (ud *UserData) GetUIUserData() *workflow.UIUserData {
	return &ud.uiUserData
}

func (ud *UserData) UpdateWithNextPrompt() {
	ud.CurrentPrompt = ud.CurrentPrompt.GetNextPrompt()

	if ud.CurrentPrompt != nil {
		// TODO cleanup -- Order might matter now but probably shouldn't
		ud.uiUserData.CurrentUIPrompt = ud.CurrentPrompt.GetUIPrompt(&ud.uiUserData)
		ud.uiUserData.CurrentUIAction = ud.CurrentPrompt.GetUIAction()
		ud.uiUserData.User.CurrentPhaseId = ud.CurrentPrompt.GetPhaseId()
		ud.uiUserData.User.CurrentPromptId = ud.CurrentPrompt.GetPromptId()
	}
}
