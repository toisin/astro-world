package db

// import (
//     "appengine"
// 	"appengine/datastore"
//     "time"
// )

type Record struct {
	RecordNo     int
	ID           string
	Name         string
	FactorIds    []string
	FactorLevels []string
	OutcomeLevel string
}
