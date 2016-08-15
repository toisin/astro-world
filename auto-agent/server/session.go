package server

// import (
// 	"github.com/gorilla/sessions"
// 	"net/http"
// )

// var store = sessions.NewCookieStore([]byte("something-very-secret"))

// func NewSession(uiUserData UiUserData, username string) {

// }

// func MyHandler(w http.ResponseWriter, r *http.Request) {
// 	// Get a session. We're ignoring the error resulted from decoding an
// 	// existing session: Get() always returns a session, even if empty.
// 	session, err := store.Get(r, "session-name")
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// Session expires in 1 day
// 	session.Options = &sessions.Options{
// 		Path:     "/",
// 		MaxAge:   86400,
// 		HttpOnly: true,
// 	}
// 	// Set some session values.
// 	session.Values["foo"] = "bar"
// 	session.Values[42] = 43
// 	// Save it before we write to the response/return from the handler.
// 	session.Save(r, w)
// }
