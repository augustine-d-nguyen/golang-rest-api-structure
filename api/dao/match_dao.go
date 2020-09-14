package dao

import (
	"context"
	"earthshaker/api/models"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//MatchDAO : this is a comment
type MatchDAO struct {
	c       *mongo.Collection
	timeOut time.Duration
}

//Setup : Set collection name
func (m *MatchDAO) Setup() {
	m.c = mgoDB.Collection("match_info")
	m.timeOut = 3 * time.Second
}

//Exist : check if the player is exist or not.
func (m *MatchDAO) Exist(id string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeOut)
	defer cancel()
	opts := options.Count().SetMaxTime(2 * time.Second)
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return true, err
	}
	num, err := m.c.CountDocuments(ctx, bson.M{"_id": objID}, opts)
	return num != 0, err
}

//FindByID : find a player status by its id.
func (m *MatchDAO) FindByID(id string) (models.Match, error) {
	var mch models.Match
	ctx, cancel := context.WithTimeout(context.Background(), m.timeOut)
	defer cancel()
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return mch, err
	}
	err = m.c.FindOne(ctx, bson.M{"_id": objID}).Decode(&mch)
	return mch, err
}

//FindMatchOf : find a player status by its id.
func (m *MatchDAO) FindMatchOf(deviceID string) (models.Match, error) {
	var mch models.Match
	ctx, cancel := context.WithTimeout(context.Background(), m.timeOut)
	defer cancel()
	var conditions = bson.M{
		"$and": []bson.M{
			bson.M{"$or": []bson.M{bson.M{"device1_id": deviceID}, bson.M{"device2_id": deviceID}}},
			bson.M{"$or": []bson.M{bson.M{"match_status": models.INIT}, bson.M{"match_status": models.WAIT}}},
		},
	}
	opts := options.Count().SetMaxTime(2 * time.Second)
	num, err := m.c.CountDocuments(ctx, conditions, opts)
	if err != nil {
		return mch, err
	}
	if num == 0 {
		return mch, errors.New("NotFound")
	}
	err = m.c.FindOne(ctx, conditions).Decode(&mch)
	return mch, err
}

//IsReadyMatch :
func (m *MatchDAO) IsReadyMatch(deviceID string, matchID string) (models.Match, error) {
	// fmt.Println("Get ready match " + deviceID)
	mch, err := m.FindByID(matchID)
	if err != nil {
		return mch, err
	}
	if mch.MatchStatus == models.INIT {
		// fmt.Println("Set init match with " + deviceID)
		mch.MatchStatus = models.WAIT
		mch.FirstConnectID = deviceID
		err = m.Upsert(mch)
	} else if mch.MatchStatus == models.WAIT && deviceID != mch.FirstConnectID {
		// fmt.Println("Set start match with " + deviceID)
		mch.MatchStatus = models.START
		err = m.Upsert(mch)
	}
	if err != nil {
		return mch, err
	}
	if mch.MatchStatus != models.START {
		return mch, errors.New("NotReady")
	}
	return mch, err
}

//CleanMatchOf :
func (m *MatchDAO) CleanMatchOf(deviceID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeOut)
	defer cancel()
	conditions := bson.M{
		"$and": []bson.M{
			bson.M{"$or": []bson.M{bson.M{"device1_id": deviceID}, bson.M{"device2_id": deviceID}}},
			bson.M{"$or": []bson.M{bson.M{"match_status": models.INIT}, bson.M{"match_status": models.WAIT}}},
		},
	}
	updateFields := bson.M{
		"$set": bson.M{"match_status": models.ERR},
	}
	_, err := m.c.UpdateMany(ctx, conditions, updateFields)
	return err
}

//Upsert : this is a comment
func (m *MatchDAO) Upsert(mch models.Match) error {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeOut)
	defer cancel()
	exist, err := m.Exist(mch.ID.Hex())
	if err != nil {
		return err
	}
	if exist {
		var updateFields = bson.M{}
		updateFields["updated_time"] = time.Now()
		if len(mch.MatchStatus) > 0 {
			updateFields["match_status"] = mch.MatchStatus
		}
		if len(mch.FirstConnectID) > 0 {
			updateFields["first_connect_id"] = mch.FirstConnectID
		}
		if len(mch.WinnerID) > 0 {
			updateFields["winner_id"] = mch.WinnerID
		}
		if len(mch.LoserID) > 0 {
			updateFields["loser_id"] = mch.LoserID
		}
		if len(mch.WebRTCOffer) > 0 {
			updateFields["webrtc_offer"] = mch.WebRTCOffer
		}
		if len(mch.WebRTCCandidates) > 0 {
			updateFields["webrtc_candidates"] = mch.WebRTCCandidates
		}
		if len(mch.WebRTCAnswer) > 0 {
			updateFields["webrtc_answer"] = mch.WebRTCAnswer
		}
		_, err = m.c.UpdateOne(ctx, bson.M{"_id": primitive.ObjectID(mch.ID)}, bson.M{"$set": updateFields})
	} else {
		mch.UpdatedTime = time.Now()
		mch.CreatedTime = mch.UpdatedTime
		_, err = m.c.InsertOne(ctx, &mch)
	}
	return err
}

//FindAllActiveMatches :
func (m *MatchDAO) FindAllActiveMatches(duration time.Duration) ([]models.Match, error) {
	pivotTime := time.Now().Add(-duration)
	ctx, cancel := context.WithTimeout(context.Background(), m.timeOut)
	defer cancel()
	findOptions := options.Find()
	conditions := bson.M{}
	conditions["$or"] = []bson.M{
		bson.M{"match_status": models.INIT},
		bson.M{"match_status": models.WAIT},
		bson.M{"match_status": models.START},
	}
	conditions["created_time"] = bson.M{
		"$lt": pivotTime,
	}
	cur, err := m.c.Find(ctx, conditions, findOptions)
	if err != nil {
		return nil, err
	}

	var results []models.Match
	for cur.Next(ctx) {
		var elem models.Match
		err := cur.Decode(&elem)
		if err != nil {
			cur.Close(ctx)
			return nil, err
		}
		results = append(results, elem)
	}
	if err := cur.Err(); err != nil {
		cur.Close(ctx)
		return nil, err
	}
	cur.Close(ctx)
	return results, nil
}

//FindLastestMatchesOf :
func (m *MatchDAO) FindLastestMatchesOf(deviceID string, limit int64) ([]models.Match, error) {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeOut)
	defer cancel()
	findOptions := options.Find()
	findOptions.SetLimit(limit)
	findOptions.SetSort(bson.M{
		"created_time": -1,
	})
	var conditions = bson.M{
		"$or":          []bson.M{bson.M{"device1_id": deviceID}, bson.M{"device2_id": deviceID}},
		"match_status": models.END,
	}
	cur, err := m.c.Find(ctx, conditions, findOptions)
	if err != nil {
		return nil, err
	}

	var results []models.Match
	for cur.Next(ctx) {
		var elem models.Match
		err := cur.Decode(&elem)
		if err != nil {
			cur.Close(ctx)
			return nil, err
		}
		results = append(results, elem)
	}
	if err := cur.Err(); err != nil {
		cur.Close(ctx)
		return nil, err
	}
	cur.Close(ctx)
	return results, nil
}

//AppendMove :
func (m *MatchDAO) AppendMove(matchID string, mv models.Move) error {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeOut)
	defer cancel()
	objID, err := primitive.ObjectIDFromHex(matchID)
	if err != nil {
		return err
	}

	// mch, err := m.FindByID(matchID)
	// if err != nil {
	// 	return err
	// }
	// if mch.MatchStatus == models.INIT {
	// 	mch.MatchStatus = models.START
	// 	err = m.Upsert(mch)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	conditions := bson.M{}
	conditions["_id"] = objID
	conditions["match_status"] = models.START
	updateTerms := bson.M{}
	updateTerms["$push"] = bson.M{
		"moves": mv,
	}
	rs := m.c.FindOneAndUpdate(ctx, conditions, updateTerms)
	if rs.Err() != nil {
		return rs.Err()
	}
	return nil
}

//FindMove :
func (m *MatchDAO) FindMove(matchID string, seq int) (models.Match, error) {
	var mch models.Match
	ctx, cancel := context.WithTimeout(context.Background(), m.timeOut)
	defer cancel()
	objID, err := primitive.ObjectIDFromHex(matchID)
	if err != nil {
		return mch, err
	}

	// mtch, err := m.FindByID(matchID)
	// if err != nil {
	// 	return mch, err
	// }
	// if mtch.MatchStatus == models.INIT {
	// 	mtch.MatchStatus = models.START
	// 	err = m.Upsert(mtch)
	// 	if err != nil {
	// 		return mch, err
	// 	}
	// }

	conditions := bson.M{}
	conditions["_id"] = objID
	conditions["match_status"] = models.START
	conditions["moves.sequence"] = seq
	num, err := m.c.CountDocuments(ctx, conditions)
	if err != nil {
		return mch, err
	}
	if num == 0 {
		return mch, errors.New("NotFound")
	}
	opts := options.FindOne()
	opts.SetProjection(bson.M{
		"moves": bson.M{
			"$elemMatch": bson.M{
				"sequence": seq,
			},
		},
	})
	err = m.c.FindOne(ctx, conditions, opts).Decode(&mch)
	return mch, err
}
