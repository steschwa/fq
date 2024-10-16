package firestore

import "cloud.google.com/go/firestore"

func GetFirestoreDirection(desc bool) firestore.Direction {
	if desc {
		return firestore.Desc
	}
	return firestore.Asc
}
