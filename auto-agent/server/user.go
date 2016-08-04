package server

import (
	"db"
	"fmt"
	"os"
	"workflow"
)

type UserData struct {
	uiUserData    *workflow.UIUserData
	CurrentPrompt workflow.Prompt
	user          db.User
}

func MakeLoginUserData(u db.User) *UserData {
	return makeAllUserData(u, true)
}

func MakeUserData(u db.User) *UserData {
	return makeAllUserData(u, false)
}

func makeAllUserData(u db.User, isNewLogin bool) *UserData {
	// Process submitted answer
	ud := &UserData{}
	ud.user = u
	ud.uiUserData = workflow.MakeUIUserData(u)
	ud.uiUserData.ArchiveHistoryLength = ud.user.ArchiveHistoryLength
	// Construct Prompt appropriately
	if (u.CurrentPromptId == "") || (u.CurrentPhaseId == "") {
		// No existing prompt, make the first one
		ud.CurrentPrompt = workflow.MakeFirstPrompt(ud.uiUserData)
		ud.user.CurrentPhaseId = ud.uiUserData.CurrentPhaseId
		ud.user.CurrentPromptId = ud.CurrentPrompt.GetPromptId()
		ud.user.CurrentSequenceOrder = ud.CurrentPrompt.GetSequenceOrder()
	} else {
		// Returning user with existing prompt, reconstruct it
		phaseId := u.CurrentPhaseId
		var promptId string
		if isNewLogin {
			// Instead of using the stored currentPromptId, use the first prompt of the sequence
			ud.user.CurrentPromptId = workflow.GetFirstPromptConfigInSequence(u.CurrentSequenceOrder, phaseId).Id
			promptId = ud.user.CurrentPromptId
		} else {
			promptId = u.CurrentPromptId
		}
		ud.CurrentPrompt = workflow.MakePrompt(promptId, phaseId, ud.uiUserData)
	}

	// update UserData with latest prompt & Ui related members
	ud.uiUserData.CurrentUIAction = ud.CurrentPrompt.GetUIAction()
	ud.uiUserData.CurrentPhaseId = ud.CurrentPrompt.GetPhaseId()
	ud.uiUserData.CurrentUIPrompt = ud.CurrentPrompt.GetUIPrompt()

	return ud
}

func (ud *UserData) UpdateWithNextPrompt() {
	ud.CurrentPrompt = ud.CurrentPrompt.GetNextPrompt()

	if ud.CurrentPrompt != nil {
		ud.uiUserData.CurrentUIPrompt = ud.CurrentPrompt.GetUIPrompt()
		ud.uiUserData.CurrentUIAction = ud.CurrentPrompt.GetUIAction()
		ud.uiUserData.CurrentPhaseId = ud.CurrentPrompt.GetPhaseId()
		ud.user.CurrentPhaseId = ud.CurrentPrompt.GetPhaseId()
		ud.user.CurrentPromptId = ud.CurrentPrompt.GetPromptId()
		ud.user.CurrentSequenceOrder = ud.CurrentPrompt.GetSequenceOrder()
	}

	if ud.uiUserData.State != nil {
		s, err := stringify(ud.uiUserData.State)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error converting StateEntities to json: %s\n\n", err)
			return
		}
		ud.user.UIState = s
	}
}
