package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//OFFLINE : player status
const (
	INIT  = "Init"
	WAIT  = "Wait"
	START = "Start"
	END   = "End"
	ERR   = "Error"
	INV   = "Invalid"
)

//Match contains match info.
type Match struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Device1ID        string             `bson:"device1_id,omitempty" json:"device1_id,omitempty"`
	Device2ID        string             `bson:"device2_id,omitempty" json:"device2_id,omitempty"`
	FirstConnectID   string             `bson:"first_connect_id",omitempty" json:"first_connect_id,omitempty"`
	MatchStatus      string             `bson:"match_status,omitempty" json:"match_status,omitempty"`
	WinnerID         string             `bson:"winner_id,omitempty" json:"winner_id,omitempty"`
	LoserID          string             `bson:"loser_id,omitempty" json:"loser_id,omitempty"`
	FirstTurnID      string             `bson:"first_turn_id,omitempty" json:"first_turn_id,omitempty"`
	WebRTCOffer      string             `bson:"webrtc_offer,omitempty" json:"webrtc_offer,omitempty"`
	WebRTCCandidates string             `bson:"webrtc_candidates,omitempty" json:"webrtc_candidates,omitempty"`
	WebRTCAnswer     string             `bson:"webrtc_answer,omitempty" json:"webrtc_answer,omitempty"`
	Moves            []Move             `bson:"moves,omitempty" json:"moves,omitempty"`
	CreatedTime      time.Time          `bson:"created_time,omitempty" json:"created_time,omitempty"`
	UpdatedTime      time.Time          `bson:"updated_time,omitempty" json:"updated_time,omitempty"`
}

//Move :
type Move struct {
	DeviceID string `bson:"device_id,omitempty" json:"device_id,omitempty"`
	Sequence int    `bson:"sequence,omitempty" json:"sequence,omitempty"`
	Step     string `bson:"step,omitempty" json:"step,omitempty"`
}
