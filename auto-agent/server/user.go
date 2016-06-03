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
	// fmt.Fprint(os.Stderr, "Inside MakeUserData, u:", u, "!\n\n")
	// fmt.Fprint(os.Stderr, "Inside MakeUserData, ud.uiUserData.User", u, "!\n\n")

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

// func (u *User)GetHistory() []db.Message {


//     c := appengine.NewContext(r)
//     // Ancestor queries, as shown here, are strongly consistent with the High
//     // Replication Datastore. Queries that span entity groups are eventually
//     // consistent. If we omitted the .Ancestor from this query there would be
//     // a slight chance that Greeting that had just been written would not
//     // show up in a query.
//     // [START query]
//     q := datastore.NewQuery("Greeting").Ancestor(guestbookKey(c)).Order("-Date").Limit(10)
//     // [END query]
//     // [START getall]
//     greetings := make([]Greeting, 0, 10)
//     if _, err := q.GetAll(c, &greetings); err != nil {
//             http.Error(w, err.Error(), http.StatusInternalServerError)
//             return
//     }
//     // [END getall]
//     if err := guestbookTemplate.Execute(w, greetings); err != nil {
//             http.Error(w, err.Error(), http.StatusInternalServerError)
//     }


// 	u.history[0] = db.Message{
// 						Username: u.username,
// 						Text: workflow.StateMap["1"].Text(),
// 						Mtype:db.ROBOT,
// 						WorkflowStateID: workflow.StateMap["1"].Id()}
// 	u.history[1] = db.Message{
// 						Username: u.username,
// 						Text: "",
// 						Mtype:db.HUMAN,
// 						WorkflowStateID: workflow.StateMap["2"].Id()}
// 	return u.history
// }

