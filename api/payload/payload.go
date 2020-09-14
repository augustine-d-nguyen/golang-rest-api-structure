package payload

//Coment :
const (
	WebRTCOfferType      = "off"
	WebRTCCandidatesType = "can"
	WebRTCAnswerType     = "ans"
)

//ResResult :
type ResResult struct {
	Result interface{} `json:"result,omitempty"`
}

//ReqUpsertStatus :
type ReqUpsertStatus struct {
	DeviceID     string `json:"device_id"`
	PlayerName   string `json:"player_name"`
	PlayerNation string `json:"player_nation"`
	PlayerStatus string `json:"player_status"`
}

//ReqFindMatch :
type ReqFindMatch struct {
	DeviceID string `json:"device_id"`
}

//ResFindMatch :
type ResFindMatch struct {
	MatchID string `json:"match_id,omitempty"`
}

//ReqReadyMatch :
type ReqReadyMatch struct {
	DeviceID string `json:"device_id"`
	MatchID  string `json:"match_id"`
}

//ResReadyMatch :
type ResReadyMatch struct {
	MatchID     string `json:"match_id"`
	FirstTurn   bool   `json:"first_turn"`
	EnemyID     string `json:"enemy_id,omitempty"`
	EnemyName   string `json:"enemy_name,omitempty"`
	EnemyNation string `json:"enemy_nation,omitempty"`
}

//ReqSendRTCMessage :
type ReqSendRTCMessage struct {
	DeviceID      string `json:"device_id"`
	MatchID       string `json:"match_id"`
	WebRTCType    string `json:"webrtc_type"`
	WebRTCMessage string `json:"webrtc_message"`
}

//ReqReceiveRTCMessage :
type ReqReceiveRTCMessage struct {
	DeviceID   string `json:"device_id"`
	MatchID    string `json:"match_id"`
	WebRTCType string `json:"webrtc_type"`
}

//ResReceiveRTCMessage :
type ResReceiveRTCMessage struct {
	MatchID       string `json:"match_id,omitempty"`
	WebRTCType    string `json:"webrtc_type,omitempty"`
	WebRTCMessage string `json:"webrtc_message,omitempty"`
}

//ReqMatchResult :
type ReqMatchResult struct {
	DeviceID string `json:"device_id"`
	MatchID  string `json:"match_id"`
	Winner   bool   `json:"winner"`
}

//ReqGetRank :
type ReqGetRank struct {
	DeviceID   string `json:"device_id"`
	MatchLimit int64  `json:"match_limit"`
}

//ResGetRank :
type ResGetRank struct {
	TopRankPlayerID     string            `json:"top_rank_id,omitempty"`
	TopRankPlayerName   string            `json:"top_rank_name,omitempty"`
	TopRankPlayerNation string            `json:"top_rank_nation,omitempty"`
	YourCurrentRank     int64             `json:"your_rank,omitempty"`
	LastestMatches      []ResGetRankMatch `json:"lastest_matches,omitempty"`
}

//ResGetRankMatch :
type ResGetRankMatch struct {
	MatchID     string `json:"match_id,omitempty"`
	MatchDate   string `json:"match_date,omitempty"`
	EnemyID     string `json:"enemy_id,omitempty"`
	EnemyName   string `json:"enemy_name,omitempty"`
	EnemyNation string `json:"enemy_nation,omitempty"`
	Win         bool   `json:"win"`
}

//ReqSendMove :
type ReqSendMove struct {
	MatchID  string `json:"match_id"`
	DeviceID string `json:"device_id"`
	Sequence int    `json:"sequence"`
	Step     string `json:"step"`
}

//ReqReceiveMove :
type ReqReceiveMove struct {
	MatchID  string `json:"match_id"`
	DeviceID string `json:"device_id"`
	Sequence int    `json:"sequence"`
}

//ResReceiveMove :
type ResReceiveMove struct {
	MatchID  string `json:"match_id,omitempty"`
	Sequence int    `json:"sequence,omitempty"`
	Step     string `json:"step,omitempty"`
}
