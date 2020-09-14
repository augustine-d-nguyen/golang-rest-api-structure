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
	DurationBeforeNow = 10 * time.Minute
)

var (
	logger    *log.Logger
	cfg       = config.Config{}
	statusDAO = dao.StatusDAO{}
	matchDAO  = dao.MatchDAO{}
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
	logger.Println("Start match cleaner service.")
	for {
		err := CleanMatchUpdateMMR()
		if err != nil {
			logger.Println(err)
			break
		}
		time.Sleep(DurationBeforeNow)
	}
}

//CleanMatchUpdateMMR :
func CleanMatchUpdateMMR() error {
	matches, err := matchDAO.FindAllActiveMatches(DurationBeforeNow)
	if err != nil {
		return err
	}
	var playerMMRs map[string]int64
	var updatingMatches []models.Match
	for _, match := range matches {
		match.UpdatedTime = time.Now()
		if match.MatchStatus == models.WAIT || match.MatchStatus == models.INIT {
			match.MatchStatus = models.ERR
			updatingMatches = append(updatingMatches, match)
		} else {
			if len(match.WinnerID) > 0 && len(match.LoserID) > 0 && match.WinnerID != match.LoserID {
				match.MatchStatus = models.END
				updatingMatches = append(updatingMatches, match)
				IncreaseMMR(&playerMMRs, match.WinnerID)
			} else if len(match.WinnerID) == 0 && len(match.LoserID) > 0 {
				if match.LoserID == match.Device1ID {
					match.WinnerID = match.Device2ID
					match.MatchStatus = models.END
					updatingMatches = append(updatingMatches, match)
					IncreaseMMR(&playerMMRs, match.WinnerID)
				} else if match.LoserID == match.Device2ID {
					match.WinnerID = match.Device1ID
					match.MatchStatus = models.END
					updatingMatches = append(updatingMatches, match)
					IncreaseMMR(&playerMMRs, match.WinnerID)
				} else {
					match.MatchStatus = models.INV
					updatingMatches = append(updatingMatches, match)
				}
			} else if len(match.WinnerID) > 0 && len(match.LoserID) == 0 {
				if match.WinnerID == match.Device1ID {
					match.LoserID = match.Device2ID
					match.MatchStatus = models.END
					updatingMatches = append(updatingMatches, match)
					IncreaseMMR(&playerMMRs, match.WinnerID)
				} else if match.WinnerID == match.Device2ID {
					match.LoserID = match.Device1ID
					match.MatchStatus = models.END
					updatingMatches = append(updatingMatches, match)
					IncreaseMMR(&playerMMRs, match.WinnerID)
				} else {
					match.MatchStatus = models.INV
					updatingMatches = append(updatingMatches, match)
				}
			} else {
				match.MatchStatus = models.INV
				updatingMatches = append(updatingMatches, match)
			}
		}
	}

	return dao.VerifyAndUpdateMMR(updatingMatches, playerMMRs)
}

//IncreaseMMR :
func IncreaseMMR(m *map[string]int64, key string) {
	AppendValue(m, key, 1)
}

//AppendValue :
func AppendValue(m *map[string]int64, key string, val int64) {
	if *m == nil {
		*m = map[string]int64{}
	}
	if _, exist := (*m)[key]; exist {
		(*m)[key] += val
	} else {
		(*m)[key] = val
	}
}
