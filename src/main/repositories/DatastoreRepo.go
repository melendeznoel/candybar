package repositories

import (
	"cloud.google.com/go/datastore"
)

func KeyExists(transaction *datastore.Transaction, key *datastore.Key) bool {
	var entityType interface{}

	if err := transaction.Get(key, &entityType); err != datastore.ErrNoSuchEntity {
		return false
	} else {
		return true
	}
}
