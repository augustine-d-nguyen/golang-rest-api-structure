package dao

import (
	"context"
	"earthshaker/api/models"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

var mgoDatabase string
var mgoClient *mongo.Client
var mgoDB *mongo.Database

//Setup : setup database name from config
func Setup(database string) {
	mgoDatabase = database
}

//Connect : this is a comment
func Connect(atlasURI string) {
	aClient, err := mongo.NewClient(options.Client().ApplyURI(atlasURI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*60*time.Second)
	defer cancel()
	err = aClient.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	err = aClient.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	mgoClient = aClient
	mgoDB = mgoClient.Database(mgoDatabase)
}

//Disconnect : this is a comment
func Disconnect() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mgoClient.Disconnect(ctx)
}

//CreateMatches :
func CreateMatches(players *[]models.Status, matches *[]models.Match) error {
	sttDAO := mgoDB.Collection("player_status")
	mchDAO := mgoDB.Collection("match_info")
	return mgoClient.UseSession(context.Background(), func(sctx mongo.SessionContext) error {
		err := sctx.StartTransaction(options.Transaction().
			SetReadConcern(readconcern.Snapshot()).
			SetWriteConcern(writeconcern.New(writeconcern.WMajority())),
		)
		if err != nil {
			return err
		}

		for i := 0; i < len(*players); i++ {
			player := (*players)[i]
			_, err := sttDAO.UpdateOne(sctx, bson.M{"device_id": player.DeviceID}, bson.M{"$set": player})
			if err != nil {
				sctx.AbortTransaction(sctx)
				return err
			}
		}
		var mches []interface{}
		for _, m := range *matches {
			mches = append(mches, m)
		}
		_, err = mchDAO.InsertMany(sctx, mches)
		if err != nil {
			sctx.AbortTransaction(sctx)
			return err
		}

		err = sctx.CommitTransaction(sctx)
		if err != nil {
			sctx.AbortTransaction(sctx)
			return err
		}
		return nil
	})
}

//VerifyAndUpdateMMR :
func VerifyAndUpdateMMR(matches []models.Match, playerMMRs map[string]int64) error {
	statusDAO := mgoDB.Collection("player_status")
	matchDAO := mgoDB.Collection("match_info")
	return mgoClient.UseSession(context.Background(), func(sctx mongo.SessionContext) error {
		err := sctx.StartTransaction(options.Transaction().
			SetReadConcern(readconcern.Snapshot()).
			SetWriteConcern(writeconcern.New(writeconcern.WMajority())),
		)
		if err != nil {
			return err
		}
		for _, match := range matches {
			_, err := matchDAO.UpdateOne(sctx, bson.M{"_id": match.ID}, bson.M{"$set": match})
			if err != nil {
				sctx.AbortTransaction(sctx)
				return err
			}
		}

		for k, v := range playerMMRs {
			rs := statusDAO.FindOneAndUpdate(sctx, bson.M{"device_id": k}, bson.M{"$inc": bson.M{"player_mmr": v}})
			if rs.Err() != nil {
				sctx.AbortTransaction(sctx)
				return rs.Err()
			}
		}

		err = sctx.CommitTransaction(sctx)
		if err != nil {
			sctx.AbortTransaction(sctx)
			return err
		}
		return nil
	})
}
