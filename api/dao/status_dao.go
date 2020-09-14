package dao

import (
	"context"
	"earthshaker/api/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//StatusDAO : this is a comment
type StatusDAO struct {
	c       *mongo.Collection
	timeOut time.Duration
}

const (
	//StatusCollection : name
	StatusCollection = "player_status"
)

//Setup : Set collection name
func (m *StatusDAO) Setup() {
	m.c = mgoDB.Collection("player_status")
	m.timeOut = 3 * time.Second
}

//Exist : check if the player is exist or not.
func (m *StatusDAO) Exist(id string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeOut)
	defer cancel()
	opts := options.Count().SetMaxTime(2 * time.Second)
	num, err := m.c.CountDocuments(ctx, bson.M{"device_id": id}, opts)
	return num != 0, err
}

//FindByID : find a player status by its id.
func (m *StatusDAO) FindByID(id string) (models.Status, error) {
	var stt models.Status
	ctx, cancel := context.WithTimeout(context.Background(), m.timeOut)
	defer cancel()
	err := m.c.FindOne(ctx, bson.M{"device_id": id}).Decode(&stt)
	return stt, err
}

//FindByIDs :
func (m *StatusDAO) FindByIDs(ids []string) ([]models.Status, error) {
	var results []models.Status
	if len(ids) == 0 {
		return results, nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), m.timeOut)
	defer cancel()
	findOptions := options.Find()
	findOptions.Projection = bson.M{
		"device_id":     1,
		"player_name":   1,
		"player_nation": 1,
	}
	conditions := bson.M{
		"device_id": bson.M{"$in": ids},
	}
	cur, err := m.c.Find(ctx, conditions, findOptions)
	if err != nil {
		return nil, err
	}
	for cur.Next(ctx) {
		var elem models.Status
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

//Upsert : this is a comment
func (m *StatusDAO) Upsert(stt models.Status) error {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeOut)
	defer cancel()
	exist, err := m.Exist(stt.DeviceID)
	if err != nil {
		return err
	}
	if exist {
		var updateFields = bson.M{
			"updated_time": time.Now(),
		}
		if len(stt.PlayerName) > 0 {
			updateFields["player_name"] = stt.PlayerName
		}
		if len(stt.PlayerStatus) > 0 {
			updateFields["player_status"] = stt.PlayerStatus
		}
		if len(stt.PlayerNation) > 0 {
			updateFields["player_nation"] = stt.PlayerNation
		}
		_, err = m.c.UpdateOne(ctx, bson.M{"device_id": stt.DeviceID}, bson.M{"$set": updateFields})
	} else {
		stt.UpdatedTime = time.Now()
		stt.CreatedTime = stt.UpdatedTime
		_, err = m.c.InsertOne(ctx, stt)
	}
	return err
}

//CountOnlinePlayers : this is a comment
func (m *StatusDAO) CountOnlinePlayers() (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeOut)
	defer cancel()
	opts := options.Count().SetMaxTime(2 * time.Second)
	num, err := m.c.CountDocuments(ctx, bson.M{"player_status": models.ONLINE}, opts)
	return num, err
}

//CountInMatchPlayers : this is a comment
func (m *StatusDAO) CountInMatchPlayers() (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeOut)
	defer cancel()
	opts := options.Count().SetMaxTime(2 * time.Second)
	num, err := m.c.CountDocuments(ctx, bson.M{"player_status": models.INMATCH}, opts)
	return num, err
}

//FindAllWaitingPlayers :
func (m *StatusDAO) FindAllWaitingPlayers() ([]models.Status, error) {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeOut)
	defer cancel()
	findOptions := options.Find()
	conditions := bson.M{}
	conditions["player_status"] = models.WAITMATCH
	cur, err := m.c.Find(ctx, conditions, findOptions)
	if err != nil {
		return nil, err
	}

	var results []models.Status
	for cur.Next(ctx) {
		var elem models.Status
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

//CalculateRankOf :
func (m *StatusDAO) CalculateRankOf(deviceID string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeOut)
	defer cancel()
	stt, err := m.FindByID(deviceID)
	if err != nil {
		return 0, nil
	}
	opts := options.Count().SetMaxTime(2 * time.Second)
	conditions := bson.M{}
	conditions["player_mmr"] = bson.M{
		"$gt": stt.PlayerMMR,
	}
	num, err := m.c.CountDocuments(ctx, conditions, opts)
	return num + 1, err
}

//FindTopRank :
func (m *StatusDAO) FindTopRank() (models.Status, error) {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeOut)
	defer cancel()
	findOptions := options.Find()
	findOptions.SetLimit(1)
	findOptions.SetSort(bson.M{
		"player_mmr": -1,
	})
	var elem models.Status
	cur, err := m.c.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return elem, err
	}
	if !cur.Next(ctx) {
		return elem, nil
	}
	err = cur.Decode(&elem)
	if err != nil {
		cur.Close(ctx)
		return elem, err
	}
	if err := cur.Err(); err != nil {
		cur.Close(ctx)
		return elem, err
	}
	cur.Close(ctx)
	return elem, nil
}
