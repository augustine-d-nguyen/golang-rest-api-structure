package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//OFFLINE : player status
const (
	OFFLINE   = "Offline"
	ONLINE    = "Online"
	WAITMATCH = "WaitMatch"
	INMATCH   = "InMatch"
)

//Status contain status
type Status struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	DeviceID     string             `bson:"device_id,omitempty" json:"device_id,omitempty"`
	PlayerName   string             `bson:"player_name,omitempty" json:"player_name,omitempty"`
	PlayerStatus string             `bson:"player_status,omitempty" json:"player_status,omitempty"`
	PlayerNation string             `bson:"player_nation,omitempty" json:"player_nation,omitempty"`
	PlayerMMR    int64              `bson:"player_mmr,omitempty" json:"player_mmr,omitempty"`
	UpdatedTime  time.Time          `bson:"updated_time,omitempty" json:"updated_time,omitempty"`
	CreatedTime  time.Time          `bson:"created_time,omitempty" json:"created_time,omitempty"`
}
