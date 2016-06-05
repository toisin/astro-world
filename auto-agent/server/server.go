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
	http.Handle(INIT_REQUEST, &ImportDBHandler{})

	workflow.InitWorkflow()
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
	http.ServeFile(w, r, "static"+r.URL.Path)
}

const COV = "/astro-world/"
const COV_STATIC = "/astro-world/js/"
const COV_REACT_STATIC = "/astro-world/react-js/"
const COV_HISTORY = "/astro-world/history"
const COV_NEWUSER = "/astro-world/newuser"
const COV_GETUSER = "/astro-world/getuser"
const COV_SENDRESPONSE = "/astro-world/sendresponse"
const INIT_REQUEST = "/astro-world/importDB"

type GetHandler StaticHandler
type HistoryHandler StaticHandler
type ResponseHandler StaticHandler
type NewUserHandler StaticHandler
type GetUserHandler StaticHandler
type ImportDBHandler StaticHandler

func (covH *ImportDBHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	ImportRecordsDB(c)
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
			fmt.Fprint(os.Stderr, "Why was I here?!\n\n")

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
		// ud, err := MakeUIUserData(c, username)
		// Query to see if user exists
		u, _, err := GetUser(c, username)

		if err != nil {
			fmt.Fprint(os.Stderr, "DB Error Getting User:"+err.Error()+"!\n\n")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		ud := MakeUserData(&u)
		ud.GetUIUserData().History, err = GetHistory(c, username)
		if err != nil {
			fmt.Fprint(os.Stderr, "DB Error Getting list of messages:"+err.Error()+"!\n\n")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		s, err := stringify(*(ud.GetUIUserData()))
		if err != nil {
			fmt.Println("Error converting messages to json", err)
		}
		fmt.Fprint(w, s)

		// fmt.Fprint(w, "{	\"prompt\": {\"type\": \""+workflow.UI_PROMPT_MC+"\", \"text\": \"First Question\", \"workflowStateID\": \"2\", \"options\": [{\"label\": \"health\", \"value\": \"X1\"},{\"label\": \"height\", \"value\": \"X2\"}]}, \"messages\": [{	\"text\": \"" + t + "\",\"type\": \"robot\"},{ \"text\": \"hello22\",\"type\": \"student\"}]}")
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
		// workflowStateID := r.FormValue("workflowStateID")
		promptId := r.FormValue("promptId")
		phaseId := r.FormValue("phaseId")

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
		ud := MakeUserData(&u)
		ud.CurrentPrompt.ProcessResponse(r.FormValue("jsonResponse"))

		responseId := ud.CurrentPrompt.GetResponse().Id
		responseText := ud.CurrentPrompt.GetResponse().Text
		questionText := ud.CurrentPrompt.GetUIPrompt().Display()

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

		// u.CurrentWorkflowStateId = workflowStateID

		//TODO cleanup
		// fmt.Fprint(os.Stderr, "Before UpdateWithNextPrompt, NextPrompt:", ud.CurrentPrompt.GetNextPrompt(), "!\n\n")
		// Move to the next prompt
		ud.UpdateWithNextPrompt()

		err = PutUser(c, u, k)
		if err != nil {
			fmt.Fprint(os.Stderr, "DB Error Put User:"+err.Error()+"!\n\n")
			return
		}

		// Update history
		ud.GetUIUserData().History, err = GetHistory(c, username)
		if err != nil {
			fmt.Fprint(os.Stderr, "DB Error Getting list of messages:"+err.Error()+"!\n\n")
			return
		}

		s, err := stringify(*(ud.GetUIUserData()))
		if err != nil {
			fmt.Println("Error converting messages to json", err)
			return
		}
		fmt.Fprint(w, s)
	}

}

func PutUser(c appengine.Context, u db.User, key *datastore.Key) (err error) {
	_, err = datastore.Put(c, key, &u)

	//    //TODO cleanup
	// fmt.Fprint(os.Stderr, "User", u, "!\n\n")
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
	q := datastore.NewQuery("Message").Ancestor(db.UserHistoryKey(c, username)).Order("MessageNo").Limit(100)
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

func stringify(v interface{}) (s string, err error) {
	b, err := json.Marshal(v)
	if err == nil {
		s = string(b[:])
	}
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
		// TODO some strange error reading from file
		// Temporarily use the const instead
		// f, err := os.Open("cases.csv")
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// r := csv.NewReader(f)

		r := csv.NewReader(strings.NewReader(workflow.CasesStream))
		fmt.Fprintf(os.Stderr, "%s")

		headers, err := r.Read()
		//TODO cleanup
		// fmt.Fprint(os.Stderr, "headers:", headers, "\n\n")

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
			for j := range workflow.GetContentConfig().CausalFactors {
				if headers[i] == workflow.GetContentConfig().CausalFactors[j].Id {
					factorColIndex[headers[i]] = i
					break
				}
			}
			for j := range workflow.GetContentConfig().NonCausalFactors {
				if headers[i] == workflow.GetContentConfig().NonCausalFactors[j].Id {
					factorColIndex[headers[i]] = i
					break
				}
			}
			if headers[i] == workflow.GetContentConfig().OutcomeVariable.Id {
				outcomeColIndex = i
				break
			}
		}

		arecord, err := r.Read()
		if err == io.EOF {
			fmt.Fprint(os.Stderr, "Record file is empty!\n\n")
			log.Fatal(err)
			return
		} else if err != nil {
			fmt.Fprint(os.Stderr, "Error reading record file!\n\n")
			log.Fatal(err)
			return
		}

		factorIds := make([]string, len(factorColIndex))
		factorLevels := make([]string, len(factorColIndex))
		i := 0
		for k, v := range factorColIndex {
			factorIds[i] = k
			factorLevels[i] = arecord[v]
			i++
		}
		outcomeLevel := arecord[outcomeColIndex]
		//TODO cleanup
		fmt.Fprint(os.Stderr, "factorIds:", factorLevels, "\n\n")
		fmt.Fprint(os.Stderr, "outcomeLevel:", outcomeLevel, "\n\n")
		// recordNo := headers[0] // First column is the record number

		//   record := []db.Record{
		// db.Record{
		// 	RecordNo: arecord[0],
		// },
		// RecordNo        int
		// ID              string
		// Name            string
		// NumberOfFactors int
		// FactorIds       []string
		// FactorLevels    []string

		//TODO cleanup
		// fmt.Fprint(os.Stderr, "headers:", factorColIndex, "\n\n")
	}
}

//     m := []db.Message{
// 		db.Message{
// 			Text: questionText,
// 			Mtype: db.ROBOT,
//     Date: time.Now(),
//     MessageNo: rc1,
// 		},
// 		db.Message{
// 			Id: responseId,
// 			Text: responseText,
// 			Mtype: db.HUMAN,
//     Date: time.Now(),
//     MessageNo: rc2,
// 		}}
// 	fmt.Fprintf(os.Stderr, "%s", record)
// }

// //check if db exists
// Record.find({}, function (err, count){
//   if (count<1) {
//     var stream = fs.readFileSync(filename, 'utf-8');
//     //console.log(stream);
//     var lines = stream.split(/\n|\r/);

//     for (var i = 1; i < lines.length; i++) {
//       var t = lines[i].split(',');
//       var c = new Cart({
//         trips: [t[4]],
//         handleLength: t[0],
//         wheelSize: t[3],
//         bucketSize: t[2],
//         bucketPlacement: t[1],
//       });

//       c.save(function(err, result) {
//         if(err)
//          console.error(err);
//         else
//          console.log(result)
//       });
//     }
//   }
// });
// };
