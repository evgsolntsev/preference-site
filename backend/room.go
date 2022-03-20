package main

import (
	"context"

	"github.com/globalsign/mgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	DatabaseName   = "preference"
	CollectionName = "rooms"
)

type RoomDAO struct {
	collection *mgo.Collection
}

func NewRoomDAO(session *mgo.Session) *RoomDAO {
	return &RoomDAO{
		collection: session.DB(DatabaseName).C(CollectionName),
	}
}

func (d *RoomDAO) FindOneByPlayer(ctx context.Context, playerName string) (*Room, error) {
	var result Room
	if err := d.collection.Find(bson.M{"sides.name": playerName}).One(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *RoomDAO) Insert(ctx context.Context, room *Room) (*Room, error) {
	if len(room.ID) == 0 {
		room.ID = primitive.NewObjectID().Hex()
	}

	if err := d.collection.Insert(room); err != nil {
		return nil, err
	}

	return room, nil
}

func (d *RoomDAO) RemoveAll(ctx context.Context) error {
	_, err := d.collection.RemoveAll(bson.M{})
	return err
}

type RoomManager struct{
	dao *RoomDAO
}

func NewRoomManager(dao *RoomDAO) *RoomManager {
	return &RoomManager{
		dao: dao,
	}
}

func (m *RoomManager) GetOneForPlayer(ctx context.Context, playerName string) (*Room, error) {
	return m.dao.FindOneByPlayer(ctx, playerName)
}
