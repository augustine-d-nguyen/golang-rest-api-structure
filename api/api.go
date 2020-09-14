package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"earthshaker/api/config"
	"earthshaker/api/dao"
	"earthshaker/api/helper"
	"earthshaker/api/models"
	"earthshaker/api/payload"

	"github.com/gorilla/mux"
)

var cfg = config.Config{}
var statusDAO = dao.StatusDAO{}
var matchDAO = dao.MatchDAO{}

//UpsertStatusEndPoint : If new device id => insert, otherwise update.
func UpsertStatusEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var reqPayload payload.ReqUpsertStatus
	if err := helper.DecodeReqBody(r.Body, &reqPayload); err != nil {
		log.Fatal(err)
		RespondWithError(w, http.StatusBadRequest, payload.ResResult{Result: "Invalid request payload"})
		return
	}
	if reqPayload.PlayerStatus == models.WAITMATCH {
		err := matchDAO.CleanMatchOf(reqPayload.DeviceID)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, payload.ResResult{Result: err.Error()})
			return
		}
	}
	var statusModel = models.Status{
		DeviceID:     reqPayload.DeviceID,
		PlayerName:   reqPayload.PlayerName,
		PlayerStatus: reqPayload.PlayerStatus,
		PlayerNation: reqPayload.PlayerNation,
	}
	if err := statusDAO.Upsert(statusModel); err != nil {
		RespondWithError(w, http.StatusInternalServerError, payload.ResResult{Result: err.Error()})
		return
	}

	RespondWithJSON(w, http.StatusOK, payload.ResResult{Result: "Success"})
}

//GetOnlinePlayersEndpoint : Get the number of online players.
func GetOnlinePlayersEndpoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if num, err := statusDAO.CountOnlinePlayers(); err != nil {
		RespondWithError(w, http.StatusInternalServerError, payload.ResResult{Result: err.Error()})
	} else {
		RespondWithJSON(w, http.StatusOK, payload.ResResult{Result: num})
	}
}

//GetMatchInfoEndPoint : Get the math info
func GetMatchInfoEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var player payload.ReqFindMatch
	if err := helper.DecodeReqBody(r.Body, &player); err != nil {
		RespondWithError(w, http.StatusBadRequest, payload.ResResult{Result: "Invalid request payload"})
		return
	}
	mch, err := matchDAO.FindMatchOf(player.DeviceID)
	if err != nil {
		if err.Error() == "NotFound" {
			RespondWithJSON(w, http.StatusOK, payload.ResFindMatch{})
			return
		}
		RespondWithError(w, http.StatusInternalServerError, payload.ResResult{Result: err.Error()})
	}
	// enemyID := mch.Device1ID
	// if player.DeviceID == mch.Device1ID {
	// 	enemyID = mch.Device2ID
	// }
	// enemy, err := statusDAO.FindByID(enemyID)
	// if err != nil {
	// 	RespondWithError(w, http.StatusInternalServerError, payload.ResResult{Result: err.Error()})
	// 	return
	// }
	var response = payload.ResFindMatch{
		MatchID: mch.ID.Hex(),
		// FirstTurn:   player.DeviceID == mch.FirstTurnID,
		// EnemyID:     enemy.DeviceID,
		// EnemyName:   enemy.PlayerName,
		// EnemyNation: enemy.PlayerNation,
	}
	RespondWithJSON(w, http.StatusOK, response)

}

//GetMatchReadyEndPoint : Get the math if it is ready
func GetMatchReadyEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var player payload.ReqReadyMatch
	if err := helper.DecodeReqBody(r.Body, &player); err != nil {
		log.Println("Invalid request payload")
		RespondWithError(w, http.StatusBadRequest, payload.ResResult{Result: "Invalid request payload"})
		return
	}
	mch, err := matchDAO.IsReadyMatch(player.DeviceID, player.MatchID)
	if err != nil {
		if err.Error() == "NotReady" {
			log.Println("Match is not ready")
			RespondWithJSON(w, http.StatusOK, payload.ResFindMatch{})
			return
		}
		RespondWithError(w, http.StatusInternalServerError, payload.ResResult{Result: err.Error()})
	}
	enemyID := mch.Device1ID
	if player.DeviceID == mch.Device1ID {
		enemyID = mch.Device2ID
	}
	enemy, err := statusDAO.FindByID(enemyID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, payload.ResResult{Result: err.Error()})
		return
	}
	var response = payload.ResReadyMatch{
		MatchID:     mch.ID.Hex(),
		FirstTurn:   player.DeviceID == mch.FirstTurnID,
		EnemyID:     enemy.DeviceID,
		EnemyName:   enemy.PlayerName,
		EnemyNation: enemy.PlayerNation,
	}
	RespondWithJSON(w, http.StatusOK, response)

}

//SendWebRTCMsgEndPoint : Send a WebRTC offer/candidate/answer
func SendWebRTCMsgEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var reqPayload payload.ReqSendRTCMessage
	if err := helper.DecodeReqBody(r.Body, &reqPayload); err != nil {
		RespondWithError(w, http.StatusBadRequest, payload.ResResult{Result: "Invalid request payload"})
		return
	}

	matchObjID, err := helper.HexToObjID(reqPayload.MatchID)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, payload.ResResult{Result: "Invalid match id"})
		return
	}
	var matchModel = models.Match{
		ID: *matchObjID,
	}
	if reqPayload.WebRTCType == payload.WebRTCOfferType {
		matchModel.WebRTCOffer = reqPayload.WebRTCMessage
		matchModel.MatchStatus = models.WAIT
	} else if reqPayload.WebRTCType == payload.WebRTCCandidatesType {
		matchModel.WebRTCCandidates = reqPayload.WebRTCMessage
	} else if reqPayload.WebRTCType == payload.WebRTCAnswerType {
		matchModel.WebRTCAnswer = reqPayload.WebRTCMessage
	} else {
		RespondWithError(w, http.StatusBadRequest, payload.ResResult{Result: "Invalid message type"})
		return
	}

	if err := matchDAO.Upsert(matchModel); err != nil {
		RespondWithError(w, http.StatusInternalServerError, payload.ResResult{Result: err.Error()})
		return
	}

	RespondWithJSON(w, http.StatusOK, payload.ResResult{Result: "Success"})
}

//ReceiveWebRTCMsgEndPoint : Receive a WebRTC offer/candidate/answer
func ReceiveWebRTCMsgEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var reqPayload payload.ReqReceiveRTCMessage
	if err := helper.DecodeReqBody(r.Body, &reqPayload); err != nil {
		RespondWithError(w, http.StatusBadRequest, payload.ResResult{Result: "Invalid request payload"})
		return
	}

	match, err := matchDAO.FindByID(reqPayload.MatchID)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, payload.ResResult{Result: "Invalid match id"})
		return
	}

	var resPayload = payload.ResReceiveRTCMessage{
		MatchID:    reqPayload.MatchID,
		WebRTCType: reqPayload.WebRTCType,
	}
	if reqPayload.WebRTCType == payload.WebRTCOfferType {
		resPayload.WebRTCMessage = match.WebRTCOffer
	} else if reqPayload.WebRTCType == payload.WebRTCCandidatesType {
		resPayload.WebRTCMessage = match.WebRTCCandidates
	} else if reqPayload.WebRTCType == payload.WebRTCAnswerType {
		resPayload.WebRTCMessage = match.WebRTCAnswer

		match.MatchStatus = models.START
		if err := matchDAO.Upsert(match); err != nil {
			RespondWithError(w, http.StatusInternalServerError, payload.ResResult{Result: err.Error()})
			return
		}
	}

	RespondWithJSON(w, http.StatusOK, resPayload)
}

//SendMoveEndPoint :
func SendMoveEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var reqPayload payload.ReqSendMove
	if err := helper.DecodeReqBody(r.Body, &reqPayload); err != nil {
		RespondWithError(w, http.StatusBadRequest, payload.ResResult{Result: "Invalid request payload"})
		return
	}

	mv := models.Move{}
	mv.DeviceID = reqPayload.DeviceID
	mv.Sequence = reqPayload.Sequence
	mv.Step = reqPayload.Step

	err := matchDAO.AppendMove(reqPayload.MatchID, mv)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, payload.ResResult{Result: err.Error()})
		return
	}
	RespondWithJSON(w, http.StatusOK, payload.ResResult{Result: "Success"})
}

//ReceiveMoveEndPoint :
func ReceiveMoveEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var reqPayload payload.ReqReceiveMove
	if err := helper.DecodeReqBody(r.Body, &reqPayload); err != nil {
		RespondWithError(w, http.StatusBadRequest, payload.ResResult{Result: "Invalid request payload"})
		return
	}

	mch, err := matchDAO.FindMove(reqPayload.MatchID, reqPayload.Sequence)
	if err != nil {
		if err.Error() == "NotFound" {
			RespondWithJSON(w, http.StatusOK, payload.ResReceiveMove{})
			return
		}
		RespondWithError(w, http.StatusInternalServerError, payload.ResResult{Result: err.Error()})
		return
	}
	resPayload := payload.ResReceiveMove{}
	resPayload.MatchID = reqPayload.MatchID
	resPayload.Sequence = reqPayload.Sequence
	resPayload.Step = mch.Moves[0].Step
	RespondWithJSON(w, http.StatusOK, resPayload)
}

//UpdateMatchResultEndPoint : Update match result.
func UpdateMatchResultEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var reqPayload payload.ReqMatchResult
	if err := helper.DecodeReqBody(r.Body, &reqPayload); err != nil {
		RespondWithError(w, http.StatusBadRequest, payload.ResResult{Result: "Invalid request payload"})
		return
	}

	matchObjID, err := helper.HexToObjID(reqPayload.MatchID)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, payload.ResResult{Result: "Invalid match id"})
		return
	}

	match, err := matchDAO.FindByID(reqPayload.MatchID)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, payload.ResResult{Result: "Invalid match id"})
		return
	}

	if match.MatchStatus == models.END || match.MatchStatus == models.ERR || match.MatchStatus == models.INV {
		RespondWithError(w, http.StatusBadRequest, payload.ResResult{Result: "Invalid request"})
		return
	}

	var matchModel = models.Match{
		ID: *matchObjID,
	}

	if reqPayload.Winner {
		if len(match.WinnerID) > 0 {
			RespondWithError(w, http.StatusBadRequest, payload.ResResult{Result: "Invalid request"})
			return
		}
		matchModel.WinnerID = reqPayload.DeviceID
	} else {
		if len(match.LoserID) > 0 {
			RespondWithError(w, http.StatusBadRequest, payload.ResResult{Result: "Invalid request"})
			return
		}
		matchModel.LoserID = reqPayload.DeviceID
	}

	if err := matchDAO.Upsert(matchModel); err != nil {
		RespondWithError(w, http.StatusInternalServerError, payload.ResResult{Result: err.Error()})
		return
	}

	RespondWithJSON(w, http.StatusOK, payload.ResResult{Result: "Success"})
}

//GetPlayerRankEndPoint :
func GetPlayerRankEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var reqPayload payload.ReqGetRank
	if err := helper.DecodeReqBody(r.Body, &reqPayload); err != nil {
		RespondWithError(w, http.StatusBadRequest, payload.ResResult{Result: "Invalid request payload"})
		return
	}
	topRank, err := statusDAO.FindTopRank()
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, payload.ResResult{Result: err.Error()})
		return
	}
	currentRank, err := statusDAO.CalculateRankOf(reqPayload.DeviceID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, payload.ResResult{Result: err.Error()})
		return
	}
	lastestMatches, err := matchDAO.FindLastestMatchesOf(reqPayload.DeviceID, reqPayload.MatchLimit)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, payload.ResResult{Result: err.Error()})
		return
	}

	resPayload := payload.ResGetRank{}
	resPayload.TopRankPlayerID = topRank.DeviceID
	resPayload.TopRankPlayerName = topRank.PlayerName
	resPayload.TopRankPlayerNation = topRank.PlayerNation
	resPayload.YourCurrentRank = currentRank
	resPayload.LastestMatches = []payload.ResGetRankMatch{}

	var enemyIDs []string
	for _, match := range lastestMatches {
		resMatch := payload.ResGetRankMatch{}
		resMatch.MatchID = match.ID.Hex()
		resMatch.MatchDate = match.CreatedTime.Format(time.RFC3339)
		if resMatch.EnemyID = match.Device1ID; match.Device1ID == reqPayload.DeviceID {
			resMatch.EnemyID = match.Device2ID
		}
		resMatch.Win = match.WinnerID == reqPayload.DeviceID

		enemyIDs = append(enemyIDs, resMatch.EnemyID)
		resPayload.LastestMatches = append(resPayload.LastestMatches, resMatch)
	}

	enemyPlayers, err := statusDAO.FindByIDs(enemyIDs)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, payload.ResResult{Result: err.Error()})
		return
	}
	// for idx, enemy := range enemyPlayers {
	// 	resPayload.LastestMatches[idx].EnemyName = enemy.PlayerName
	// 	resPayload.LastestMatches[idx].EnemyNation = enemy.PlayerNation
	// }
	for idx, rankMatch := range resPayload.LastestMatches {
		enemyID := rankMatch.EnemyID
		for _, enemy := range enemyPlayers {
			if enemyID == enemy.DeviceID {
				resPayload.LastestMatches[idx].EnemyName = enemy.PlayerName
				resPayload.LastestMatches[idx].EnemyNation = enemy.PlayerNation
				break
			}
		}
	}

	RespondWithJSON(w, http.StatusOK, resPayload)
}

//GetWelcomeEndpoint :
func GetWelcomeEndpoint(w http.ResponseWriter, r *http.Request) {
	RespondWithJSON(w, http.StatusOK, payload.ResResult{Result: "Welcome to earthshaker APIs"})
}

//RespondWithError :
func RespondWithError(w http.ResponseWriter, code int, payload interface{}) {
	RespondWithJSON(w, code, payload)
}

//RespondWithJSON :
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

//AuthMiddleware :
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("x-earthshaker-token")
		if token == cfg.APIKey {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Forbidden", http.StatusForbidden)
		}
	})
}

//ContentTypeMiddleware :
func ContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Content-Type")
		if token == "application/json" {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Content-Type header is not application/json", http.StatusUnsupportedMediaType)
		}
	})
}

func init() {
	log.SetOutput(os.Stdout)
	cfg.Read()
	dao.Setup(cfg.Database)
	dao.Connect(cfg.AtlasURI)

	statusDAO.Setup()
	matchDAO.Setup()
}

func main() {
	defer dao.Disconnect()

	r := mux.NewRouter()
	r.HandleFunc("/earthshaker/v1/welcome", GetWelcomeEndpoint).Methods("GET")
	api := r.PathPrefix("/earthshaker/v1").Subrouter()
	api.Use(AuthMiddleware)
	api.Use(ContentTypeMiddleware)
	api.HandleFunc("/player/status", GetOnlinePlayersEndpoint).Methods("GET")
	api.HandleFunc("/player/status/upsert", UpsertStatusEndPoint).Methods("POST")
	api.HandleFunc("/player/rank", GetPlayerRankEndPoint).Methods("POST")
	api.HandleFunc("/match/info", GetMatchInfoEndPoint).Methods("POST")
	api.HandleFunc("/match/ready", GetMatchReadyEndPoint).Methods("POST")
	api.HandleFunc("/match/info/update", UpdateMatchResultEndPoint).Methods("PUT")
	// r.HandleFunc("/match/webrtc/send", SendWebRTCMsgEndPoint).Methods("POST")
	// r.HandleFunc("/match/webrtc/receive", ReceiveWebRTCMsgEndPoint).Methods("POST")
	api.HandleFunc("/match/sync/send", SendMoveEndPoint).Methods("POST")
	api.HandleFunc("/match/sync/receive", ReceiveMoveEndPoint).Methods("POST")

	log.Println("Try starting server at port 6526")
	err := http.ListenAndServe(":6526", r)
	if err != nil {
		log.Fatal(err)
	}
}
