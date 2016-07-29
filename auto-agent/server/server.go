package server

import (
	"encoding/csv"
	"encoding/json"

	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"db"
	"log"
	"workflow"

	"appengine"
	"appengine/datastore"
)

const (
	APP_NAME = "auto-agent"
)

func init() {
	http.Handle("/", &StaticHandler{})

	http.Handle(COV, &GetHandler{})
	http.Handle(COV_STATIC, &StaticHandler{})
	http.Handle(COV_REACT_STATIC, &StaticHandler{})
	http.Handle(COV_HISTORY, &HistoryHandler{})
	http.Handle(COV_NEWUSER, &NewUserHandler{})
	http.Handle(COV_GETUSER, &GetUserHandler{})
	http.Handle(COV_SENDRESPONSE, &ResponseHandler{})

	//TODO should not rely on a separate http request but it only needs to happen once
	// needs to find a better place
	http.Handle(IMPORTDB_REQUEST, &ImportRecordDBHandler{})
	http.Handle(CLEARDB_REQUEST, &ClearRecordDBHandler{})
	http.Handle(CLEARALLUSERS_REQUEST, &ClearAllUsersDBHandler{})

	workflow.InitWorkflow()
}

type TextHandler string

func (t *TextHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, t)
}

type StaticHandler struct {
}

func (staticH *StaticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Note that the path must not start with / for some reasons
	// i.e. "/static..." does not work. Has to be "static..."
	http.ServeFile(w, r, "static"+r.URL.Path)
}

const COV = "/astro-world/"
const COV_STATIC = "/astro-world/js/"
const COV_REACT_STATIC = "/astro-world/react-js/"
const COV_HISTORY = "/astro-world/history"
const COV_NEWUSER = "/astro-world/newuser"
const COV_GETUSER = "/astro-world/getuser"
const COV_SENDRESPONSE = "/astro-world/sendresponse"
const IMPORTDB_REQUEST = "/astro-world/importDB"
const CLEARDB_REQUEST = "/astro-world/clearDB"
const CLEARALLUSERS_REQUEST = "/astro-world/clearAllUsersDB"

type GetHandler StaticHandler
type HistoryHandler StaticHandler
type ResponseHandler StaticHandler
type NewUserHandler StaticHandler
type GetUserHandler StaticHandler
type ImportRecordDBHandler StaticHandler
type ClearRecordDBHandler StaticHandler
type ClearAllUsersDBHandler StaticHandler

func (covH *ImportRecordDBHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	ImportRecordsDB(c)
	http.ServeFile(w, r, "static/index.html")
}

func (covH *ClearRecordDBHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	ClearAllRecordsDB(c)
	http.ServeFile(w, r, "static/index.html")
}

func (covH *ClearAllUsersDBHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	ClearAllUsersDB(c)
	http.ServeFile(w, r, "static/index.html")
}

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
			fmt.Fprint(os.Stderr, "DB Error Getting User:"+err.Error()+"!\n\n")
			http.ServeFile(w, r, "static/index.html")
			return
		}

		if u.Username == "" {
			fmt.Fprint(os.Stderr, "User does not exist.\n\n")

			http.ServeFile(w, r, "static/index.html")
			return
		}

		if len(r.URL.Path[len(COV):]) != 0 {
			http.ServeFile(w, r, "static/astro-world"+r.URL.Path)
			return
		} else {
			http.ServeFile(w, r, "static/astro-world/index.html")
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
		// Query to see if user exists
		u, _, err := GetUser(c, username)

		if err != nil {
			fmt.Fprint(os.Stderr, "DB Error Getting User:"+err.Error()+"!\n\n")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		ud := MakeUserData(u)
		ud.uiUserData.History, err = GetHistory(c, username)
		if err != nil {
			fmt.Fprint(os.Stderr, "DB Error Getting list of messages:"+err.Error()+"!\n\n")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		s, err := stringify(ud.uiUserData)
		if err != nil {
			fmt.Println("Error converting messages to json", err)
		}
		fmt.Fprint(w, string(s[:]))

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
			fmt.Fprint(os.Stderr, "DB Error Getting User:"+err.Error()+"!\n\n")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if u.Username != "" {
			http.Error(w, "Cannot create user. Username already exists!", 500)
			return
		}

		u = db.User{
			Username:   username,
			Screenname: r.FormValue("screenname"),
			Date:       time.Now(),
		}

		key := db.UserKey(c)
		_, err = datastore.Put(c, key, &u)
		if err != nil {
			fmt.Fprint(os.Stderr, "DB Error Creating User:"+err.Error()+"!\n\n")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		s, err := stringify(u)
		if err != nil {
			fmt.Println("Error converting user object to json", err)
		}
		fmt.Fprint(w, string(s[:]))

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
			fmt.Fprint(os.Stderr, "DB Error Getting User:"+err.Error()+"!\n\n")
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
		promptId := r.FormValue("promptId")
		phaseId := r.FormValue("phaseId")
		questionText := r.FormValue("questionText")

		// Query to see if user exists
		u, k, err := GetUser(c, username)

		if err != nil {
			fmt.Fprint(os.Stderr, "DB Error Getting User:"+err.Error()+"!\n\n")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// First time user will have empty currentPhaseId & currentPromptId
		if u.CurrentPhaseId != "" {
			if u.CurrentPhaseId != phaseId {
				fmt.Fprint(os.Stderr, "Out of sync error! User info and DB are out of sync.\n\n. Revert to what's in the DB")
			}
		}
		if u.CurrentPromptId != "" {
			if u.CurrentPromptId != promptId {
				fmt.Fprint(os.Stderr, "Out of sync error! User info and DB are out of sync.\n\n. Revert to what's in the DB")
			}
		}

		// Process submitted answers
		ud := MakeUserData(u)
		ud.CurrentPrompt.ProcessResponse(r.FormValue("jsonResponse"), &ud.user, &ud.uiUserData, c)

		responseId := ud.CurrentPrompt.GetResponseId()
		responseText := ud.CurrentPrompt.GetResponseText()

		// Get the count of existing messages
		rc, err := GetHistoryCount(c, username)
		if err != nil {
			fmt.Fprint(os.Stderr, "DB Error Getting count of messages:"+err.Error()+"!\n\n")
			return
		}

		//TODO need to find a way to save the responses that are not text
		//Process submitted answers
		rc++
		rc1 := rc
		rc++
		rc2 := rc

		m := []db.Message{
			db.Message{
				Text:      questionText,
				Mtype:     db.ROBOT,
				Date:      time.Now(),
				MessageNo: rc1,
			},
			db.Message{
				Id:        responseId,
				Text:      responseText,
				Mtype:     db.HUMAN,
				Date:      time.Now(),
				MessageNo: rc2,
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
			fmt.Fprint(os.Stderr, "DB Error Adding Messages:"+err.Error()+"!\n\n")
			return
			//http.Error(w, err.Error(), http.StatusInternalServerError)
			//return
		}

		// Move to the next prompt
		ud.UpdateWithNextPrompt()

		// Store updated user in DB
		err = PutUser(c, ud.user, k)
		if err != nil {
			fmt.Fprint(os.Stderr, "DB Error Put User:"+err.Error()+"!\n\n")
			return
		}

		// Update history
		ud.uiUserData.History, err = GetHistory(c, username)
		if err != nil {
			fmt.Fprint(os.Stderr, "DB Error Getting list of messages:"+err.Error()+"!\n\n")
			return
		}

		s, err := stringify(ud.uiUserData)
		if err != nil {
			fmt.Println("Error converting messages to json", err)
			return
		}
		fmt.Fprint(w, string(s[:]))
	}

}

func PutUser(c appengine.Context, u db.User, key *datastore.Key) (err error) {
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

func GetHistory(c appengine.Context, username string) (messages []*db.Message, err error) {
	q := datastore.NewQuery("Message").Ancestor(db.UserHistoryKey(c, username)).Order("MessageNo").Limit(100)
	// [END query]
	// [START getall]
	messages = make([]*db.Message, 0, 100)
	_, err = q.GetAll(c, &messages)
	return
}

func GetHistoryCount(c appengine.Context, username string) (rc int, err error) {
	q := datastore.NewQuery("Message").Ancestor(db.UserHistoryKey(c, username)).Limit(100)
	rc, err = q.Count(c)
	return
}

func stringify(v interface{}) (b []byte, err error) {
	b, err = json.Marshal(v)
	return
}

// DOC
// Expect the first column to be record number
// The rest of column headers should match the factor ids in configuration
func ImportRecordsDB(c appengine.Context) {

	q := datastore.NewQuery("Record")
	rc, err := q.Count(c)
	if err != nil {
		fmt.Fprint(os.Stderr, "DB Error Getting count of records:"+err.Error()+"!\n\n")
		return
	}

	if rc < 1 {
		f, err := os.Open(workflow.GetContentConfig().RecordFileName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s", err)
			log.Fatal(err)
		}
		// Note: Do not use mac version of csv has problems with the end of line symbols
		r := csv.NewReader(f)

		headers, err := r.Read()

		if err == io.EOF {
			fmt.Fprint(os.Stderr, "Record file is empty!\n\n")
			log.Fatal(err)
			return
		} else if err != nil {
			fmt.Fprint(os.Stderr, "Error reading record file!\n\n")
			log.Fatal(err)
			return
		}

		factorColIndex := make(map[string]int)
		var outcomeColIndex int
		for i := range headers {
			if workflow.GetFactorConfig(headers[i]).Id != "" {
				factorColIndex[headers[i]] = i
			} else if headers[i] == workflow.GetContentConfig().OutcomeVariable.Id {
				outcomeColIndex = i
			}
		}

		records := make([]db.Record, workflow.GetContentConfig().RecordSize)
		var keys = make([]*datastore.Key, workflow.GetContentConfig().RecordSize)

		ri := 0
		for {
			arecord, err := r.Read()
			if err == io.EOF {
				break
			} else if err != nil {
				fmt.Fprint(os.Stderr, "Error reading record file!\n\n", err)
				log.Fatal(err)
				return
			}
			recordNo := arecord[0]  // First column is the record number
			firstname := arecord[1] // Second column is the first name
			lastname := arecord[2]  // Third column is the last name
			factorIds := make([]string, len(factorColIndex))
			factorLevels := make([]string, len(factorColIndex))
			for k, v := range factorColIndex {
				f := workflow.GetFactorConfig(k)
				i := f.DBIndex
				factorIds[i] = k
				factorLevels[i] = arecord[v]
			}
			outcomeLevel := arecord[outcomeColIndex]
			records[ri] = db.Record{
				RecordNo:     recordNo,
				Firstname:    firstname,
				Lastname:     lastname,
				FactorId0:    factorIds[0],
				FactorId1:    factorIds[1],
				FactorId2:    factorIds[2],
				FactorId3:    factorIds[3],
				FactorId4:    factorIds[4],
				FactorLevel0: factorLevels[0],
				FactorLevel1: factorLevels[1],
				FactorLevel2: factorLevels[2],
				FactorLevel3: factorLevels[3],
				FactorLevel4: factorLevels[4],
				OutcomeLevel: outcomeLevel,
			}
			keys[ri] = datastore.NewIncompleteKey(c, "Record", db.RecordKey(c, APP_NAME))
			ri++
		}
		_, err = datastore.PutMulti(c, keys, records)
		if err != nil {
			fmt.Fprint(os.Stderr, "DB Error Adding Records:"+err.Error()+"!\n\n")
			log.Fatal(err)
			return
		}
	}
}

func ClearAllRecordsDB(c appengine.Context) {
	q := datastore.NewQuery("Record")

	var records []db.Record
	// To retrieve the results,
	// you must execute the Query using its GetAll or Run methods.
	ks, err := q.GetAll(c, &records)
	if err != nil {
		fmt.Fprint(os.Stderr, "DB Error Getting Record:"+err.Error()+"!\n\n")
		log.Fatal(err)
		return
	}

	err = datastore.DeleteMulti(c, ks)
	if err != nil {
		fmt.Fprint(os.Stderr, "DB Error Deleting Records:"+err.Error()+"!\n\n")
		log.Fatal(err)
		return
	}
}

func ClearAllUsersDB(c appengine.Context) {
	q := datastore.NewQuery("User")

	var us []db.User
	// To retrieve the results,
	// you must execute the Query using its GetAll or Run methods.
	ks, err := q.GetAll(c, &us)
	if err != nil {
		fmt.Fprint(os.Stderr, "DB Error Getting All Users:"+err.Error()+"!\n\n")
		log.Fatal(err)
		return
	}

	err = datastore.DeleteMulti(c, ks)
	if err != nil {
		fmt.Fprint(os.Stderr, "DB Error Deleting All Users:"+err.Error()+"!\n\n")
		log.Fatal(err)
		return
	}

	q = datastore.NewQuery("Message")

	var ms []db.Message
	// To retrieve the results,
	// you must execute the Query using its GetAll or Run methods.
	ks, err = q.GetAll(c, &ms)
	if err != nil {
		fmt.Fprint(os.Stderr, "DB Error Getting All Messages:"+err.Error()+"!\n\n")
		log.Fatal(err)
		return
	}

	err = datastore.DeleteMulti(c, ks)
	if err != nil {
		fmt.Fprint(os.Stderr, "DB Error Deleting All Messages:"+err.Error()+"!\n\n")
		log.Fatal(err)
		return
	}
}
