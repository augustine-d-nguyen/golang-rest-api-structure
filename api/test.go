package main

import (
	"fmt"

	"earthshaker/api/config"
	"earthshaker/api/dao"
	"earthshaker/api/helper"
	"earthshaker/api/models"
)

var cfg = config.Config{}
var statusDAO = dao.StatusDAO{}
var matchDAO = dao.MatchDAO{}

func init() {
	cfg.Read()
	dao.Setup(cfg.Database)
	dao.Connect(cfg.AtlasURI)

	statusDAO.Setup()
	matchDAO.Setup()
}

func main() {
	// testStatusDAOUpsert()
	// testStatusDAOCountOnline()
	// testMatchDAOFindOne()
	// fmt.Println(helper.IsValidHexID("fsadsa"))
	// testMatchDAOCreateOne()
	// testMatchDAOFindMatchOf()
	testStatusDAOFindAllWaitingPlayers()
}

func testStatusDAOUpsert() {
	var stt = models.Status{
		DeviceID:     "abcd12346",
		PlayerName:   "Augustine D. Nguyen",
		PlayerStatus: models.ONLINE,
	}
	err := statusDAO.Upsert(stt)
	if err != nil {
		fmt.Println(err)
	}
}

func testStatusDAOCountOnline() {
	num, err := statusDAO.CountOnlinePlayers()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(num)
	}
}

func testMatchDAOFindOne() {
	mch, err := matchDAO.FindByID("5e4b58bacb7fb7da7586d8d0")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(mch)
	}
}
func testMatchDAOCreateOne() {
	objID, err := helper.HexToObjID("5e4b696cda7a6270ffe5b5f1")
	var mch = models.Match{
		ID:          *objID,
		MatchStatus: models.WAIT,
	}
	err = matchDAO.Upsert(mch)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Done")
	}
}

func testMatchDAOFindMatchOf() {

	mch, err := matchDAO.FindMatchOf("abcd12346")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(mch)
	}
}

func testStatusDAOFindAllWaitingPlayers() {
	stts, err := statusDAO.FindAllWaitingPlayers()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%+v\n", stts)
	}
}
