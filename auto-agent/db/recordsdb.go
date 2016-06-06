package db

import (
	"appengine"
	"appengine/datastore"
)

type Record struct {
	RecordNo     string
	ID           string
	Name         string
	FactorIds    []string
	FactorLevels []string
	OutcomeLevel string
}

// RecordKey returns the key used for all records.
func RecordKey(c appengine.Context, appname string) *datastore.Key {
	return datastore.NewKey(c, "Records", appname, 0, nil)
}
