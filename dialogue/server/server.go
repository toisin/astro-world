package server

import (
    "fmt"
    "net/http"
    "os"
    "time"
    "encoding/json"
    "errors"
    "strings"
 
    "workflow"
    "db"

    "appengine"
    "appengine/datastore"
)

func init() {
    http.Handle("/", &StaticHandler{})

    http.Handle(COV, &GetHandler{})
    http.Handle(COV_STATIC, &StaticHandler{})
    http.Handle(COV_HISTORY, &HistoryHandler{})
    http.Handle(COV_NEWUSER, &NewUserHandler{})
    http.Handle(COV_GETUSER, &GetUserHandler{})
    http.Handle(COV_SENDRESPONSE, &ResponseHandler{})

    workflow.InitWorkflowMaps()
}

type TextHandler string

func (t TextHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, t)
}

type StaticHandler struct {

}

func (staticH *StaticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Note that the path must not start with / for some reasons
	// i.e. "/static..." does not work. Has to be "static..."
    http.ServeFile(w, r, "static" + r.URL.Path)
}

const COV = "/cov/"
const COV_STATIC = "/cov/js/"
const COV_HISTORY = "/cov/history"
const COV_NEWUSER = "/cov/newuser"
const COV_GETUSER = "/cov/getuser"
const COV_SENDRESPONSE = "/cov/sendresponse"

type GetHandler StaticHandler
type HistoryHandler StaticHandler
type ResponseHandler StaticHandler
type NewUserHandler StaticHandler
type GetUserHandler StaticHandler

func (covH *GetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	// Check if user has logged in (by simply providing a username in the query parameter for now
	// because we are not checking password at the moment)
	// If logged in, serve the requested static file (assuming all URL request
	// without its own path handler is a request to a file in the static folder
	// If not logged in, redirect to the parent index page for login
	if r.URL.Query()["user"] != nil {

		// Check to make sure that the provided user actually exists
		username := strings.ToLower(r.URL.Query()["user"][0])
		u, _, err := GetUser(c, username)

	    if err != nil {
			fmt.Fprint(os.Stderr, "DB Error Getting User:" + err.Error() + "!\n\n")
	        http.ServeFile(w, r, "static/index.html")
	        return
	    }

	    if u.Username == "" {
			fmt.Fprint(os.Stderr, "Why was I hear?!\n\n")

	    	http.ServeFile(w, r, "static/index.html")
	    	return
		}

		if len(r.URL.Path[len(COV):]) != 0 {
			http.ServeFile(w, r, "static/cov" + r.URL.Path)
			return
		} else {
			http.ServeFile(w, r, "static/cov/index.html")
			return
		}
	} else {
		http.ServeFile(w, r, "static/index.html")
		return
	}
}

func (covH *HistoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if r.URL.Query()["user"] != nil {
		// Always handle username in lowercase
		username := strings.ToLower(r.URL.Query()["user"][0])
		ud, err := UpdateUserData(c, username, "")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		s, err := stringify(ud)
		if err != nil {
			fmt.Println("Error converting messages to json", err)
		}
		fmt.Fprint(w, s)

		// fmt.Fprint(w, "{	\"prompt\": {\"type\": \""+workflow.PROMPT_MC+"\", \"text\": \"First Question\", \"workflowStateID\": \"2\", \"options\": [{\"label\": \"health\", \"value\": \"X1\"},{\"label\": \"height\", \"value\": \"X2\"}]}, \"messages\": [{	\"text\": \"" + t + "\",\"type\": \"robot\"},{ \"text\": \"hello22\",\"type\": \"student\"}]}")
	} else {
		fmt.Fprint(os.Stderr, "Error: username not provided for getting history!\n\n")
	}
}

func (newuserH *NewUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	if r.FormValue("user") != "" {
		// Always handle username in lowercase
		username := strings.ToLower(r.FormValue("user"))

	    // Query to see if user exists
 		u, _, err := GetUser(c, username)

	    if err != nil {
			fmt.Fprint(os.Stderr, "DB Error Getting User:" + err.Error() + "!\n\n")
	        http.Error(w, err.Error(), http.StatusInternalServerError)
	        return
	    }

	    if u.Username != "" {
	    	http.Error(w, "Cannot create user. Username already exists!", 500)
	    	return
		}

        u = db.User{
				Username: username,
				Screenname: r.FormValue("screenname"),
			    Date: time.Now(),
        }

	    key := db.UserKey(c)
	    _, err = datastore.Put(c, key, &u)
	    if err != nil {
			fmt.Fprint(os.Stderr, "DB Error Creating User:" + err.Error() + "!\n\n")
	        http.Error(w, err.Error(), http.StatusInternalServerError)
	        return
	    }

		s, err := stringify(u)
		if err != nil {
			fmt.Println("Error converting user object to json", err)
		}
		fmt.Fprint(w, s)

	}

}

func (newuserH *GetUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	if r.FormValue("user") != "" {
		// Always handle username in lowercase
		username := strings.ToLower(r.FormValue("user"))

	    q := datastore.NewQuery("User").
	            Filter("Username=", username)

	    // To retrieve the results,
	    // you must execute the Query using its GetAll or Run methods.
	    rc, err := q.Count(c)

	    if err != nil {
			fmt.Fprint(os.Stderr, "DB Error Getting User:" + err.Error() + "!\n\n")
	        http.Error(w, err.Error(), http.StatusInternalServerError)
	        return
	    }

	    if rc != 1 {
	    	http.Error(w, "There is a problem with the username!", 500)
	    	return
		}

		fmt.Fprint(w, "")
	}

}

func (covH *ResponseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	if r.FormValue("user") != "" {
		// Always handle username in lowercase
		username := strings.ToLower(r.FormValue("user"))
		workflowStateID := r.FormValue("workflowStateID")
        
        // Get the count of existing messages
		rc, err := GetHistoryCount(c, username)
	    if err != nil {
	    	fmt.Fprint(os.Stderr, "DB Error Getting count of messages:" + err.Error() + "!\n\n")
	        return
	    }

	    rc++
	    rc1:= rc
	    rc++
	    rc2:= rc

        m := []db.Message{
        		db.Message{
					Text: r.FormValue("questionText"),
					Mtype: db.ROBOT,
					WorflowStateID: workflowStateID,
				    Date: time.Now(),
				    RecordNo: rc1,
        		},
        		db.Message{
					Value: r.FormValue("responseValue"),
					Text: r.FormValue("responseText"),
					Mtype: db.HUMAN,
					WorflowStateID: workflowStateID,
				    Date: time.Now(),
				    RecordNo: rc2,
        		}}

	    // TODO what does this comment mean?
	    // We set the same parent key on every Message entity to ensure each Message
	    // is in the same entity group. Queries across the single entity group
	    // will be consistent. However, the write rate to a single entity group
	    // should be limited to ~1/second.
        var keys = []*datastore.Key{
        			datastore.NewIncompleteKey(c, "Message", db.UserHistoryKey(c, username)),
	    			datastore.NewIncompleteKey(c, "Message", db.UserHistoryKey(c, username))}
	    			
	    _, err = datastore.PutMulti(c, keys, m)
	    if err != nil {
			fmt.Fprint(os.Stderr, "DB Error Adding Messages:" + err.Error() + "!\n\n")
	        //http.Error(w, err.Error(), http.StatusInternalServerError)
	        //return
	    }

	    // Query to see if user exists
 		u, k, err := GetUser(c, username)

	    if err != nil {
			fmt.Fprint(os.Stderr, "DB Error Getting User:" + err.Error() + "!\n\n")
	        http.Error(w, err.Error(), http.StatusInternalServerError)
	        return
	    }

        u.CurrentWorkflowStateId = workflowStateID

	    err = PutUser(c, u, k)
	    if err != nil {
			fmt.Fprint(os.Stderr, "DB Error Put User:" + err.Error() + "!\n\n")
	        //http.Error(w, err.Error(), http.StatusInternalServerError)
	        //return
	    }

		// var history = make([]db.Message, 10)
		//TODO Probably should not be creating new records of history everytime?
		// var u = User{username:r.FormValue("user"), history:history}
		// history := GetHistory(w, c, username)

		// fmt.Fprint(w, "{	\"prompt\": {\"type\": \"TEXT\", \"text\": \"First Question\", \"workflowStateID\": \"2\"}, \"messages\": [{	\"text\": \"" + history[0].Text + "\",\"type\": \"robot\"},{ \"text\": \"hello22\",\"type\": \"student\"}]}")

		ud, err := UpdateUserData(c, username, workflowStateID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		s, err := stringify(ud)
		if err != nil {
			fmt.Println("Error converting messages to json", err)
		}
		fmt.Fprint(w, s)
	}

}


func PutUser(c appengine.Context, u db.User, key *datastore.Key) (err error){
    _, err = datastore.Put(c, key, &u)
    return
}


func GetUser(c appengine.Context, username string) (u db.User, k *datastore.Key, err error) {
    q := datastore.NewQuery("User").Ancestor(db.UserListKey(c)).
            Filter("Username=", username)

    var users []db.User
    // To retrieve the results,
    // you must execute the Query using its GetAll or Run methods.
    ks, err := q.GetAll(c, &users)

    if len(users) > 1 {
    	err = errors.New("Error getting history: More than one user found!")
    	return
	} else if len(users) != 0 {
		u = users[0]
		k = ks[0]
	}
    return
}

func GetHistory(c appengine.Context, username string) (messages []db.Message, err error) {
    q := datastore.NewQuery("Message").Ancestor(db.UserHistoryKey(c, username)).Order("RecordNo").Limit(100)
    // [END query]
    // [START getall]
    messages = make([]db.Message, 0, 100)
    _, err = q.GetAll(c, &messages)
    return
}

func GetHistoryCount(c appengine.Context, username string) (rc int, err error) {
    q := datastore.NewQuery("Message").Ancestor(db.UserHistoryKey(c, username)).Limit(100)
    rc, err = q.Count(c)
    return
}

func UpdateUserData(c appengine.Context, username string, workflowStateID string) (ud UserData, err error) {

	ud = UserData {}

	u, _, err := GetUser(c, username)

    if err != nil {
		fmt.Fprint(os.Stderr, "DB Error Getting User to update:" + err.Error() + "!\n\n")
        return
    }

    if u.Username == "" {
    	err = errors.New("Error updating userdata: Cannot find user!")
    	return
	}

	ud.User = u
    // TODO What does this comment mean...
    // Ancestor queries, as shown here, are strongly consistent with the High
    // Replication Datastore. Queries that span entity groups are eventually
    // consistent. If we omitted the .Ancestor from this query there would be
    // a slight chance that Message that had just been written would not
    // show up in a query.
    // [START query]
 
	ud.History, err = GetHistory(c, username)
    if err != nil {
    	fmt.Fprint(os.Stderr, "DB Error Getting list of messages:" + err.Error() + "!\n\n")
        return
    }

    var currentWorkflowStateId string
    if workflowStateID == "" {
	    currentWorkflowStateId = ud.User.CurrentWorkflowStateId
	} else {
		currentWorkflowStateId = workflowStateID
	}

	if currentWorkflowStateId == "" {
		// Start from the beginning
		ud.CurrentPrompt = workflow.GetFirstState()
		ud.User.CurrentWorkflowStateId = ud.CurrentPrompt.GetId()
    	//fmt.Fprint(os.Stderr, ud.CurrentPrompt.Ptype())
	} else if (currentWorkflowStateId) != "" && (currentWorkflowStateId != workflow.PROMPT_END) {
		nid := workflow.GetStateMap()[currentWorkflowStateId].GetNextStateId()
		ud.CurrentPrompt = workflow.GetStateMap()[nid]
		ud.User.CurrentWorkflowStateId = nid
	}

    return
}

func stringify(v interface{}) (s string, err error) {
	b, err := json.Marshal(v)
	if err == nil {
		s = string(b[:])
	}
	return
}
