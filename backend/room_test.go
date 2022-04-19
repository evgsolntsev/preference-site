package main

import (
	"context"
	"testing"

	"github.com/globalsign/mgo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type RoomSuite struct {
	suite.Suite

	Ctx     context.Context
	DAO     *RoomDAO
	Manager *RoomManager
}

func TestRoom(t *testing.T) {
	session, err := mgo.Dial("mongodb://localhost:27017")
	require.NoError(t, err)

	ctx := context.Background()
	dao := NewRoomDAO(session)
	defer dao.RemoveAll(ctx)

	manager := NewRoomManager(dao)

	suite.Run(t, &RoomSuite{
		Ctx:     ctx,
		DAO:     dao,
		Manager: manager,
	})
}

func (s *RoomSuite) TestRoomDAOFindByPlayer() {
	r := &Room{
		Sides: []RoomSideInfo{{
			Name: "evgsol",
		}, {
			Name: "solarka",
		}},
	}

	_, err := s.DAO.Insert(s.Ctx, r)
	require.NoError(s.T(), err)

	found, err := s.DAO.FindOneByPlayer(s.Ctx, "evgsol")
	require.NoError(s.T(), err)
	require.NotNil(s.T(), found)

	found, err = s.DAO.FindOneByPlayer(s.Ctx, "solarka")
	require.NoError(s.T(), err)
	require.NotNil(s.T(), found)

	found, err = s.DAO.FindOneByPlayer(s.Ctx, "miracle")
	require.Error(s.T(), err)
	require.Nil(s.T(), found)
}

func (s *RoomSuite) TestRoomManagerOpenBuypackOK() {
	room, err := s.DAO.Insert(s.Ctx, &Room{
		Sides: []RoomSideInfo{{
			Name: "evgsol",
		}, {
			Name: "solarka",
		}, {
			Name:  "lol",
			Cards: []Card{{}, {}},
		}, {
			Name: "kek",
		}},
		Status: RoomStatusCreated,
	})
	require.NoError(s.T(), err)

	err = s.Manager.OpenBuypack(s.Ctx, room.ID)
	require.NoError(s.T(), err)

	updatedRoom, err := s.DAO.FindOneByID(s.Ctx, room.ID)
	require.NoError(s.T(), err)

	assert.True(s.T(), updatedRoom.Sides[2].Open)
	assert.Equal(s.T(), RoomStatusBuypackOpened, updatedRoom.Status)
}

func (s *RoomSuite) TestRoomManagerOpenBuypackWrongStatus() {
	room, err := s.DAO.Insert(s.Ctx, &Room{
		Status: RoomStatusBuypackOpened,
	})
	require.NoError(s.T(), err)

	err = s.Manager.OpenBuypack(s.Ctx, room.ID)
	assert.Error(s.T(), err)
}

func (s *RoomSuite) TestRoomManagerTakeBuypackOK() {
	room, err := s.DAO.Insert(s.Ctx, &Room{
		Sides: []RoomSideInfo{{
			Name:  "evgsol",
			Cards: []Card{{SuitSpades, "A"}},
		}, {
			Name:  "solarka",
			Cards: []Card{{SuitDiamonds, "A"}},
		}, {
			Name:  "lol",
			Cards: []Card{{SuitClubs, "7"}, {SuitClubs, "8"}},
		}, {
			Name:  "kek",
			Cards: []Card{{SuitHearts, "A"}},
		}},
		Status: RoomStatusBuypackOpened,
	})
	require.NoError(s.T(), err)

	err = s.Manager.TakeBuypack(s.Ctx, room.ID, "evgsol")
	require.NoError(s.T(), err)

	updatedRoom, err := s.DAO.FindOneByID(s.Ctx, room.ID)
	require.NoError(s.T(), err)

	assert.Equal(s.T(), RoomStatusBuypackTaken, updatedRoom.Status)
	assert.Equal(s.T(), []Card{}, updatedRoom.Sides[2].Cards)
	assert.Equal(s.T(), []Card{
		{SuitSpades, "A"}, {SuitClubs, "7"}, {SuitClubs, "8"},
	}, updatedRoom.Sides[0].Cards)
}

func (s *RoomSuite) TestRoomManagerTakeBuypackWrongStatus() {
	room, err := s.DAO.Insert(s.Ctx, &Room{
		Status: RoomStatusBuypackTaken,
	})
	require.NoError(s.T(), err)

	err = s.Manager.TakeBuypack(s.Ctx, room.ID, "lol")
	assert.Error(s.T(), err)
}
