package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sort"

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

func (d *RoomDAO) FindOneByID(ctx context.Context, roomID string) (*Room, error) {
	var result Room
	if err := d.collection.Find(bson.M{"_id": roomID}).One(&result); err != nil {
		return nil, err
	}

	return &result, nil
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

func (d *RoomDAO) Update(ctx context.Context, room *Room) error {
	return d.collection.UpdateId(room.ID, room)
}

func (d *RoomDAO) OpenBuypack(ctx context.Context, roomID string, buypackIndex int) error {
	return d.collection.Update(bson.M{
		"_id":    roomID,
		"status": RoomStatusCreated,
	}, bson.M{
		"$set": bson.M{
			"status": RoomStatusBuypackOpened,
			fmt.Sprintf("sides.%d.open", buypackIndex): true,
		},
	})
}

func (d *RoomDAO) RemoveAll(ctx context.Context) error {
	_, err := d.collection.RemoveAll(bson.M{})
	return err
}

type RoomManager struct {
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

func (m *RoomManager) Shuffle(ctx context.Context, roomID, playerName string) error {
	room, err := m.dao.FindOneByPlayer(ctx, playerName)
	if err != nil {
		return err
	}

	if room.PlayersCount < 3 || room.PlayersCount > 4 {
		return errors.New("wrong players count")
	}

	var allCards []Card
	for _, s := range AllSuits {
		for _, r := range AllRanks {
			allCards = append(allCards, Card{
				Suit: s,
				Rank: r,
			})
		}
	}

	rand.Shuffle(len(allCards), func(i, j int) { allCards[i], allCards[j] = allCards[j], allCards[i] })

	playerIndex := room.PlayerSideIndex(playerName)
	buypackIndex := 0
	var playersIndexes []int
	if room.PlayersCount == 3 {
		for i := 0; i < 4; i++ {
			index := (playerIndex + i) % 4
			if room.Sides[index].Name == EmptySidePlayerName {
				buypackIndex = index
			} else {
				playersIndexes = append(playersIndexes, index)
			}
		}
	} else {
		buypackIndex = playerIndex
		playersIndexes = []int{(playerIndex + 1) % 4, (playerIndex + 2) % 4, (playerIndex + 3) % 4}
	}

	room.Status = RoomStatusCreated
	room.Sides[buypackIndex].Cards = allCards[:2]
	room.Sides[buypackIndex].Tricks = 0
	room.Sides[buypackIndex].Open = false
	room.Center = nil
	for i := 0; i < 3; i++ {
		room.Sides[playersIndexes[i]].Cards = allCards[2+i*10 : 2+(i+1)*10]
		sort.Slice(room.Sides[playersIndexes[i]].Cards, func(l, r int) bool {
			return room.Sides[playersIndexes[i]].Cards[l].Less(room.Sides[playersIndexes[i]].Cards[r])
		})
		room.Sides[playersIndexes[i]].Tricks = 0
		room.Sides[playersIndexes[i]].Open = false
	}

	return m.dao.Update(ctx, room)
}

func (m *RoomManager) OpenBuypack(ctx context.Context, roomID string) error {
	room, err := m.dao.FindOneByID(ctx, roomID)
	if err != nil {
		return err
	}

	if room.Status != RoomStatusCreated {
		return errors.New("wrong room status")
	}

	buypackIndex := -1
	for i, side := range room.Sides {
		if len(side.Cards) == 2 {
			buypackIndex = i
		}
	}

	if buypackIndex == -1 {
		return errors.New("bad shuffling")
	}

	return m.dao.OpenBuypack(ctx, roomID, buypackIndex)
}