package main

import (
	"earthshaker/api/config"
	"earthshaker/api/dao"
	"earthshaker/api/models"
	"log"
	"os"
	"time"
)

//Comment :
const (
	MinInterval = 3 * time.Second
	MaxInterval = 6 * time.Second
	MatchSize   = 2
)

var (
	logger    *log.Logger
	cfg       = config.Config{}
	statusDAO = dao.StatusDAO{}
	matchDAO  = dao.MatchDAO{}

	interval = MinInterval
)

func init() {
	logger = log.New(os.Stderr, "ERR: ", log.Ldate|log.Ltime|log.Lshortfile)

	cfg.Read()
	dao.Setup(cfg.Database)
	dao.Connect(cfg.AtlasURI)

	statusDAO.Setup()
	matchDAO.Setup()
}

func main() {
	defer dao.Disconnect()
	logger.Println("Start match maker service.")
	for {
		err := MakeMatch()
		if err != nil {
			logger.Println(err)
			break
		}
		// logger.Printf("Try to sleep %+v", interval)
		time.Sleep(interval)
	}
}

//MakeMatch :
func MakeMatch() error {
	// logger.Println("Start match maker")
	// - Get all players are waiting for a match
	players, err := statusDAO.FindAllWaitingPlayers()
	if err != nil {
		return err
	}

	var updatingPlayers []models.Status
	var creatingMatches []models.Match
	for len(players) >= MatchSize {
		player1, player2 := players[0], players[1]
		// logger.Printf("Player1 %+v", player1)
		// logger.Printf("Player2 %+v", player2)
		players = players[2:]
		newMatch := models.Match{
			MatchStatus: models.INIT,
			Device1ID:   player1.DeviceID,
			Device2ID:   player2.DeviceID,
			FirstTurnID: player2.DeviceID,
			CreatedTime: time.Now(),
			UpdatedTime: time.Now(),
		}
		// logger.Printf("New Match %+v", newMatch)
		creatingMatches = append(creatingMatches, newMatch)
		player1.PlayerStatus = models.INMATCH
		player2.PlayerStatus = models.INMATCH
		updatingPlayers = append(updatingPlayers, player1, player2)
	}
	if len(creatingMatches) > 0 {
		interval -= MinInterval
		if interval < MinInterval {
			interval = MinInterval
		}

		err := dao.CreateMatches(&updatingPlayers, &creatingMatches)
		if err != nil {
			return err
		}
	} else {
		interval += MinInterval
		if interval > MaxInterval {
			interval = MaxInterval
		}
	}
	return nil
}
