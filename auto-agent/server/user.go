package server

import (
	"db"
	"fmt"
	"os"
	"workflow"
)

type UserData struct {
	uiUserData    workflow.UIUserData
	CurrentPrompt workflow.Prompt
	user          db.User
}

func MakeUserData(u db.User) *UserData {
	// Process submitted answer
	ud := &UserData{}
	ud.user = u
	ud.uiUserData = *workflow.MakeUIUserData(u)
	// ud.uiUserData = workflow.UIUserData{}

	// // update new UserData with everything that is available on db.User
	// ud.uiUserData.Username = u.Username
	// ud.uiUserData.Screenname = u.Screenname
	// ud.uiUserData.CurrentFactorId = u.CurrentFactorId
	// ud.uiUserData.CurrentPhaseId = u.CurrentPhaseId

	// if u.UIState != nil {
	// 	s, err := workflow.UnstringifyState(u.UIState, u.CurrentPhaseId)

	// 	if err != nil {
	// 		fmt.Fprintf(os.Stderr, "Error converting json to StateEntities: %s\n\n", err)
	// 	}
	// 	ud.uiUserData.State = s
	// }

	// Construct Prompt appropriately
	if (u.CurrentPromptId == "") || (u.CurrentPhaseId == "") {
		// No existing prompt, make the first one
		ud.CurrentPrompt = workflow.MakeFirstPrompt(&ud.uiUserData)
		ud.user.CurrentPromptId = ud.CurrentPrompt.GetPromptId()
	} else {
		// Returning user with existing prompt, reconstruct it
		phaseId := u.CurrentPhaseId
		promptId := u.CurrentPromptId
		ud.CurrentPrompt = workflow.MakePrompt(promptId, phaseId, &ud.uiUserData)
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
