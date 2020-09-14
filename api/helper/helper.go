package helper

import (
	"encoding/json"
	"io"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//GenerateObjID :
func GenerateObjID() *primitive.ObjectID {
	id := primitive.NewObjectID()
	return &id
}

//HexToObjID :
func HexToObjID(hex string) (*primitive.ObjectID, error) {
	id, err := primitive.ObjectIDFromHex(hex)
	return &id, err
}

//IsValidHexID :
func IsValidHexID(id string) bool {
	_, err := primitive.ObjectIDFromHex(id)
	return err == nil
}

//IsValidObjID :
func IsValidObjID(id *primitive.ObjectID) bool {
	return id.Hex() != "" && IsValidHexID(id.Hex())
}

//DecodeReqBody :
func DecodeReqBody(reqBody io.ReadCloser, payload interface{}) error {
	dec := json.NewDecoder(reqBody)
	dec.DisallowUnknownFields()
	return dec.Decode(payload)
}
