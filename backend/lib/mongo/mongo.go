package mongo

import (
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OID = primitive.ObjectID

// OIDFromHex wraps converting string to OID
func OIDFromHex(id string) (OID, error) {
	return primitive.ObjectIDFromHex(id)
}

// for use when we pre-validate the hex string
func OIDFromValidHex(id string) OID {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("OIDFromValidHex: called with invalid hex: " + id)
	}

	return oid
}
