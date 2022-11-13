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
	RoomDatabaseName   = "preference"
	RoomCollectionName = "rooms"
)

type RoomDAO struct {
	collection *mgo.Collection
}

func NewRoomDAO(session *mgo.Session) *RoomDAO {
	return &RoomDAO{
		collection: session.DB(RoomDatabaseName).C(RoomCollectionName),
	}
}

func (d *RoomDAO) FindOneByID(ctx context.Context, roomID string) (*Room, error) {
	var result Room
	if err := d.collection.FindId(roomID).One(&result); err != nil {
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

func (d *RoomDAO) TakeBuypack(
	ctx context.Context,
	roomID string,
	buypackIndex, playerIndex int,
	newPlayerCards []Card,
) error {
	return d.collection.Update(bson.M{
		"_id":    roomID,
		"status": RoomStatusBuypackOpened,
	}, bson.M{
		"$set": bson.M{
			"status": RoomStatusBuypackTaken,
			fmt.Sprintf("sides.%d.cards", buypackIndex): []Card{},
			fmt.Sprintf("sides.%d.cards", playerIndex):  newPlayerCards,
		},
	})
}

func (d *RoomDAO) Drop(
	ctx context.Context,
	roomID string,
	playerIndex int,
	newPlayerCards []Card,
) error {
	return d.collection.Update(bson.M{
		"_id":    roomID,
		"status": RoomStatusBuypackTaken,
	}, bson.M{
		"$set": bson.M{
			"status": RoomStatusPlaying,
			fmt.Sprintf("sides.%d.cards", playerIndex): newPlayerCards,
		},
	})
}

func (d *RoomDAO) Move(
	ctx context.Context,
	roomID string,
	playerIndex int,
	newCenterCard CenterCardInfo,
	newPlayerCards []Card,
) error {
	return d.collection.Update(bson.M{
		"_id": roomID,
		"status": bson.M{
			"$in": []RoomStatus{
				RoomStatusPlaying,
				RoomStatusAllPass,
			},
		},
	}, bson.M{
		"$set": bson.M{
			fmt.Sprintf("sides.%d.cards", playerIndex): newPlayerCards,
		},
		"$push": bson.M{
			"center": newCenterCard,
		},
	})
}

func (d *RoomDAO) TakeTrick(
	ctx context.Context,
	roomID string,
	buypackIndex int,
	playerIndex int,
	oldCenterCards []CenterCardInfo,
	newCenterCards []CenterCardInfo,
) error {
	return d.collection.Update(bson.M{
		"_id": roomID,
		"status": bson.M{
			"$in": []RoomStatus{
				RoomStatusPlaying,
				RoomStatusAllPass,
			},
		},
	}, bson.M{
		"$set": bson.M{
			"center":    newCenterCards,
			"lastTrick": oldCenterCards,
			fmt.Sprintf("sides.%d.cards", buypackIndex): []Card{},
		},
		"$inc": bson.M{
			fmt.Sprintf("sides.%d.tricks", playerIndex): 1,
		},
	})
}

func (d *RoomDAO) AllPass(
	ctx context.Context,
	roomID string,
	buypackIndex int,
	newBuypackCards []Card,
	newCenterCards []CenterCardInfo,
) error {
	return d.collection.Update(bson.M{
		"_id":    roomID,
		"status": RoomStatusCreated,
	}, bson.M{
		"$set": bson.M{
			"status": RoomStatusAllPass,
			"center": newCenterCards,
			fmt.Sprintf("sides.%d.cards", buypackIndex): newBuypackCards,
		},
	})
}

func (d *RoomDAO) ChangeVisibility(
	ctx context.Context,
	roomID string,
	playerIndex int,
	open bool,
) error {
	return d.collection.Update(bson.M{
		"_id": roomID,
	}, bson.M{
		"$set": bson.M{
			fmt.Sprintf("sides.%d.open", playerIndex): open,
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
	room.BuypackIndex = buypackIndex
	room.LastTrick = []CenterCardInfo{}
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

	return m.dao.OpenBuypack(ctx, roomID, room.BuypackIndex)
}

func (m *RoomManager) TakeBuypack(ctx context.Context, roomID, playerName string) error {
	room, err := m.dao.FindOneByID(ctx, roomID)
	if err != nil {
		return err
	}

	if room.Status != RoomStatusBuypackOpened {
		return errors.New("wrong room status")
	}

	playerIndex := -1
	for i, side := range room.Sides {
		if side.Name == playerName {
			playerIndex = i
		}
	}

	if playerIndex == -1 {
		return errors.New("wrong player name")
	}

	cards := append(room.Sides[playerIndex].Cards, room.Sides[room.BuypackIndex].Cards...)
	sort.Slice(cards, func(l, r int) bool {
		return cards[l].Less(cards[r])
	})

	return m.dao.TakeBuypack(ctx, roomID, room.BuypackIndex, playerIndex, cards)
}

func (m *RoomManager) Drop(ctx context.Context, roomID, playerName string, indexes []int) error {
	room, err := m.dao.FindOneByID(ctx, roomID)
	if err != nil {
		return err
	}

	if room.Status != RoomStatusBuypackTaken {
		return errors.New("wrong room status")
	}

	playerIndex := -1
	for i, side := range room.Sides {
		if side.Name == playerName {
			playerIndex = i
		}
	}

	if playerIndex == -1 {
		return errors.New("wrong player name")
	}

	if len(room.Sides[playerIndex].Cards) != 12 {
		return errors.New("wrong player cards length")
	}

	newCards := []Card{}
	for i, c := range room.Sides[playerIndex].Cards {
		good := true
		for _, index := range indexes {
			if i == index {
				good = false
			}
		}
		if good {
			newCards = append(newCards, c)
		}
	}

	return m.dao.Drop(ctx, roomID, playerIndex, newCards)
}

func (m *RoomManager) Move(ctx context.Context, roomID, playerName string, index int) error {
	room, err := m.dao.FindOneByID(ctx, roomID)
	if err != nil {
		return err
	}

	if room.Status != RoomStatusPlaying && room.Status != RoomStatusAllPass {
		return errors.New("wrong room status")
	}

	for _, centerCard := range room.Center {
		if playerName == centerCard.Player {
			return errors.New("player have made the move already")
		}
	}

	playerIndex := -1
	for i, side := range room.Sides {
		if side.Name == playerName {
			playerIndex = i
		}
	}

	if playerIndex == -1 {
		return errors.New("wrong player name")
	}

	if len(room.Sides[playerIndex].Cards) <= index {
		return errors.New("wrong player cards length")
	}

	newCenterCard := CenterCardInfo{
		Player: playerName,
	}
	newCards := []Card{}
	for i, c := range room.Sides[playerIndex].Cards {
		if i == index {
			newCenterCard.Card = c
		} else {
			newCards = append(newCards, c)
		}
	}

	return m.dao.Move(ctx, roomID, playerIndex, newCenterCard, newCards)
}

func (m *RoomManager) TakeTrick(ctx context.Context, roomID, playerName string) error {
	room, err := m.dao.FindOneByID(ctx, roomID)
	if err != nil {
		return err
	}

	if room.Status != RoomStatusPlaying && room.Status != RoomStatusAllPass {
		return errors.New("wrong room status")
	}

	playerIndex := -1
	for i, side := range room.Sides {
		if side.Name == playerName {
			playerIndex = i
		}
	}

	if playerIndex == -1 {
		return errors.New("wrong player name")
	}

	if len(room.Center) < 3 {
		return errors.New("unable to take trick")
	}

	newCenter := []CenterCardInfo{}
	if len(room.Sides[room.BuypackIndex].Cards) > 0 {
		newCenter = []CenterCardInfo{{
			Card:   room.Sides[room.BuypackIndex].Cards[0],
			Player: room.Sides[room.BuypackIndex].Name,
		}}
	}
	return m.dao.TakeTrick(ctx, roomID, room.BuypackIndex, playerIndex, room.Center, newCenter)
}

func (m *RoomManager) AllPass(ctx context.Context, roomID string) error {
	room, err := m.dao.FindOneByID(ctx, roomID)
	if err != nil {
		return err
	}

	if room.Status != RoomStatusCreated {
		return errors.New("wrong room status")
	}

	newCenterCards := []CenterCardInfo{{
		Card:   room.Sides[room.BuypackIndex].Cards[0],
		Player: room.Sides[room.BuypackIndex].Name,
	}}
	newBuypackCards := room.Sides[room.BuypackIndex].Cards[1:]
	return m.dao.AllPass(ctx, roomID, room.BuypackIndex, newBuypackCards, newCenterCards)
}

func (m *RoomManager) ChangeVisibility(ctx context.Context, roomID, playerName string) error {
	room, err := m.dao.FindOneByID(ctx, roomID)
	if err != nil {
		return err
	}

	playerIndex := -1
	for i, side := range room.Sides {
		if side.Name == playerName {
			playerIndex = i
		}
	}

	if playerIndex == -1 {
		return errors.New("wrong player name")
	}

	return m.dao.ChangeVisibility(ctx, roomID, playerIndex, !room.Sides[playerIndex].Open)
}
