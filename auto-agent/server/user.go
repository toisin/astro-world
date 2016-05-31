package server

import (
	"db"
	"workflow"
)

type UserData struct {
	User db.User
	History []db.Message
	CurrentPrompt workflow.OldState
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

