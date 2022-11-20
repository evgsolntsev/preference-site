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

func (s *RoomSuite) TearDownTest() {
	s.DAO.RemoveAll(s.Ctx)
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
		BuypackIndex: 2,
		Status:       RoomStatusReady,
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
		BuypackIndex: 2,
		Status:       RoomStatusBuypackOpened,
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

func (s *RoomSuite) TestRoomManagerDrop() {
	room, err := s.DAO.Insert(s.Ctx, &Room{
		Sides: []RoomSideInfo{{
			Name:  "evgsol",
			Cards: []Card{{SuitSpades, "A"}},
		}, {
			Name:  "solarka",
			Cards: []Card{{SuitDiamonds, "A"}},
		}, {
			Name: "lol",
			Cards: []Card{
				{SuitClubs, "7"},
				{SuitClubs, "8"},
				{SuitClubs, "9"},
				{SuitClubs, "10"},
				{SuitClubs, "J"},
				{SuitClubs, "Q"},
				{SuitClubs, "K"},
				{SuitClubs, "A"},
				{SuitHearts, "7"},
				{SuitHearts, "8"},
				{SuitHearts, "9"},
				{SuitHearts, "10"},
			},
		}, {
			Name:  "kek",
			Cards: []Card{{SuitHearts, "A"}},
		}},
		Status: RoomStatusBuypackTaken,
	})
	require.NoError(s.T(), err)

	err = s.Manager.Drop(s.Ctx, room.ID, "lol", []int{8, 9})
	require.NoError(s.T(), err)

	updatedRoom, err := s.DAO.FindOneByID(s.Ctx, room.ID)
	require.NoError(s.T(), err)

	assert.Equal(s.T(), RoomStatusPlaying, updatedRoom.Status)
	assert.Equal(s.T(), []Card{
		{SuitClubs, "7"},
		{SuitClubs, "8"},
		{SuitClubs, "9"},
		{SuitClubs, "10"},
		{SuitClubs, "J"},
		{SuitClubs, "Q"},
		{SuitClubs, "K"},
		{SuitClubs, "A"},
		{SuitHearts, "9"},
		{SuitHearts, "10"},
	}, updatedRoom.Sides[2].Cards)
}

func (s *RoomSuite) TestRoomManagerMove() {
	room, err := s.DAO.Insert(s.Ctx, &Room{
		Sides: []RoomSideInfo{{
			Name:  "evgsol",
			Cards: []Card{{SuitSpades, "A"}},
		}, {
			Name:  "solarka",
			Cards: []Card{{SuitDiamonds, "A"}},
		}, {
			Name: "lol",
			Cards: []Card{
				{SuitClubs, "7"},
				{SuitClubs, "8"},
			},
		}, {
			Name: "kek",
		}},
		Center: []CenterCardInfo{{
			Card:   Card{SuitSpades, "K"},
			Player: "evgsol",
		}, {
			Card:   Card{SuitSpades, "Q"},
			Player: "solarka",
		}},
		Status: RoomStatusPlaying,
	})
	require.NoError(s.T(), err)

	err = s.Manager.Move(s.Ctx, room.ID, "lol", 0)
	require.NoError(s.T(), err)

	updatedRoom, err := s.DAO.FindOneByID(s.Ctx, room.ID)
	require.NoError(s.T(), err)

	assert.Equal(s.T(), []Card{
		{SuitClubs, "8"},
	}, updatedRoom.Sides[2].Cards)
	assert.Equal(s.T(), []CenterCardInfo{{
		Card:   Card{SuitSpades, "K"},
		Player: "evgsol",
	}, {
		Card:   Card{SuitSpades, "Q"},
		Player: "solarka",
	}, {
		Card:   Card{SuitClubs, "7"},
		Player: "lol",
	}}, updatedRoom.Center)

	err = s.Manager.Move(s.Ctx, room.ID, "lol", 0)
	require.Error(s.T(), err)
}

func (s *RoomSuite) TestRoomManagerTakeTrick() {
	room, err := s.DAO.Insert(s.Ctx, &Room{
		Sides: []RoomSideInfo{{
			Name:  "evgsol",
			Cards: []Card{{SuitSpades, "A"}},
		}, {
			Name:  "solarka",
			Cards: []Card{{SuitDiamonds, "A"}},
		}, {
			Name:   "lol",
			Cards:  []Card{{SuitClubs, "7"}, {SuitClubs, "8"}},
			Tricks: 5,
		}, {
			Name:  "kek",
			Cards: []Card{{SuitHearts, "A"}},
		}},
		Center: []CenterCardInfo{{
			Card:   Card{SuitSpades, "K"},
			Player: "evgsol",
		}, {
			Card:   Card{SuitSpades, "Q"},
			Player: "solarka",
		}, {
			Card:   Card{SuitClubs, "8"},
			Player: "lol",
		}},
		Status:       RoomStatusAllPass,
		BuypackIndex: 0,
	})
	require.NoError(s.T(), err)

	err = s.Manager.TakeTrick(s.Ctx, room.ID, "lol")
	require.NoError(s.T(), err)

	updatedRoom, err := s.DAO.FindOneByID(s.Ctx, room.ID)
	require.NoError(s.T(), err)

	assert.Equal(s.T(), []CenterCardInfo{{
		Card:   Card{SuitSpades, "A"},
		Player: "evgsol",
	}}, updatedRoom.Center)
	assert.Equal(s.T(), []Card{}, updatedRoom.Sides[0].Cards)
	assert.Equal(s.T(), []CenterCardInfo{{
		Card:   Card{SuitSpades, "K"},
		Player: "evgsol",
	}, {
		Card:   Card{SuitSpades, "Q"},
		Player: "solarka",
	}, {
		Card:   Card{SuitClubs, "8"},
		Player: "lol",
	}}, updatedRoom.LastTrick)
	assert.Equal(s.T(), 6, updatedRoom.Sides[2].Tricks)
}

func (s *RoomSuite) TestRoomManagerTakeTrickWithTwoCenterCards() {
	room, err := s.DAO.Insert(s.Ctx, &Room{
		Sides: []RoomSideInfo{{
			Name:  "evgsol",
			Cards: []Card{{SuitSpades, "A"}},
		}, {
			Name:  "solarka",
			Cards: []Card{{SuitDiamonds, "A"}},
		}, {
			Name:   "lol",
			Cards:  []Card{{SuitClubs, "7"}, {SuitClubs, "8"}},
			Tricks: 5,
		}, {
			Name:  "kek",
			Cards: []Card{{SuitHearts, "A"}},
		}},
		Center: []CenterCardInfo{{
			Card:   Card{SuitSpades, "K"},
			Player: "evgsol",
		}, {
			Card:   Card{SuitSpades, "Q"},
			Player: "solarka",
		}},
		Status:       RoomStatusAllPass,
		BuypackIndex: 0,
	})
	require.NoError(s.T(), err)

	err = s.Manager.TakeTrick(s.Ctx, room.ID, "lol")
	require.Error(s.T(), err)
}

func (s *RoomSuite) TestRoomManagerAllPass() {
	room, err := s.DAO.Insert(s.Ctx, &Room{
		Sides: []RoomSideInfo{{
			Name: "evgsol",
		}, {
			Name: "solarka",
		}, {
			Name:  "lol",
			Cards: []Card{{SuitSpades, "Q"}, {SuitSpades, "K"}},
		}, {
			Name: "kek",
		}},
		BuypackIndex: 2,
		Status:       RoomStatusReady,
	})
	require.NoError(s.T(), err)

	err = s.Manager.AllPass(s.Ctx, room.ID)
	require.NoError(s.T(), err)

	updatedRoom, err := s.DAO.FindOneByID(s.Ctx, room.ID)
	require.NoError(s.T(), err)

	assert.False(s.T(), updatedRoom.Sides[2].Open)
	assert.Equal(s.T(), RoomStatusAllPass, updatedRoom.Status)
	assert.Equal(s.T(), []CenterCardInfo{{
		Card:   Card{SuitSpades, "Q"},
		Player: "lol",
	}}, updatedRoom.Center)
	assert.Equal(s.T(), []Card{{SuitSpades, "K"}}, updatedRoom.Sides[2].Cards)
}

func (s *RoomSuite) TestRoomManagerChangeVisibility() {
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
			Open:  true,
		}},
		BuypackIndex: 2,
		Status:       RoomStatusBuypackOpened,
	})
	require.NoError(s.T(), err)

	err = s.Manager.ChangeVisibility(s.Ctx, room.ID, "kek")
	require.NoError(s.T(), err)

	updatedRoom, err := s.DAO.FindOneByID(s.Ctx, room.ID)
	require.NoError(s.T(), err)

	assert.Equal(s.T(), RoomStatusBuypackOpened, updatedRoom.Status)
	assert.False(s.T(), updatedRoom.Sides[3].Open)
}

func (s *RoomSuite) TestPlayerIn() {
	s.Run("Smoke", func() {
		room, err := s.DAO.Insert(s.Ctx, &Room{
			Sides: []RoomSideInfo{{
				Name: "evgsol",
			}, {
				Name: "solarka",
			}, {
				Name: "lol",
			}, {
				Name: EMPTY_SIDE,
			}},
			Status:       RoomStatusCreated,
			PlayersCount: 3,
		})
		require.NoError(s.T(), err)

		require.NoError(s.T(), s.Manager.PlayerIn(s.Ctx, room.ID, "kek"))

		updatedRoom, err := s.DAO.FindOneByID(s.Ctx, room.ID)
		require.NoError(s.T(), err)

		require.Equal(s.T(), "kek", updatedRoom.Sides[3].Name)
		require.Equal(s.T(), 4, updatedRoom.PlayersCount)
	})

	s.Run("Player already in room", func() {
		_, err := s.DAO.Insert(s.Ctx, &Room{
			Sides: []RoomSideInfo{{
				Name: "evgsol",
			}},
			Status: RoomStatusCreated,
		})
		require.NoError(s.T(), err)

		room, err := s.DAO.Insert(s.Ctx, &Room{
			Status: RoomStatusCreated,
		})
		require.NoError(s.T(), err)

		err = s.Manager.PlayerIn(s.Ctx, room.ID, "evgsol")
		require.Error(s.T(), err)
		require.Equal(s.T(), "player is already in room", err.Error())
	})

	s.Run("Wrong room status", func() {
		room, err := s.DAO.Insert(s.Ctx, &Room{
			Status: RoomStatusReady,
		})
		require.NoError(s.T(), err)

		err = s.Manager.PlayerIn(s.Ctx, room.ID, "donald trump")
		require.Error(s.T(), err)
		require.Equal(s.T(), "wrong room status", err.Error())
	})

	s.Run("No empty sides", func() {
		room, err := s.DAO.Insert(s.Ctx, &Room{
			Sides: []RoomSideInfo{{
				Name: "evgsol",
			}, {
				Name: "solarka",
			}, {
				Name: "lol",
			}, {
				Name: "kek",
			}},
			Status:       RoomStatusCreated,
			PlayersCount: 4,
		})
		require.NoError(s.T(), err)

		err = s.Manager.PlayerIn(s.Ctx, room.ID, "joe biden")
		require.Error(s.T(), err)
		require.Equal(s.T(), "no empty sides", err.Error())
	})
}

func (s *RoomSuite) TestRoomReady() {
	s.Run("Smoke", func() {
		room, err := s.DAO.Insert(s.Ctx, &Room{
			Sides: []RoomSideInfo{{
				Name: "evgsol",
			}, {
				Name: "solarka",
			}, {
				Name: "lol",
			}, {
				Name: "kek",
			}},
			PlayersCount: 4,
			Status:       RoomStatusCreated,
		})

		err = s.Manager.RoomReady(s.Ctx, "evgsol")
		require.NoError(s.T(), err)

		updatedRoom, err := s.DAO.FindOneByID(s.Ctx, room.ID)
		require.NoError(s.T(), err)

		require.Equal(s.T(), RoomStatusReady, updatedRoom.Status)

		err = s.Manager.RoomReady(s.Ctx, "evgsol")
		require.NoError(s.T(), err)
	})

	s.Run("WrongStatus", func() {
		_, err := s.DAO.Insert(s.Ctx, &Room{
			Sides: []RoomSideInfo{{
				Name: "evgsol",
			}, {
				Name: "solarka",
			}, {
				Name: "lol",
			}, {
				Name: "elon mask",
			}},
			PlayersCount: 4,
			Status:       RoomStatusAllPass,
		})

		err = s.Manager.RoomReady(s.Ctx, "elon mask")
		require.Error(s.T(), err)
		require.Equal(s.T(), "wrong room status", err.Error())
	})

	s.Run("WrongPlayersCount", func() {
		_, err := s.DAO.Insert(s.Ctx, &Room{
			Sides: []RoomSideInfo{{
				Name: "joe biden",
			}, {
				Name: "solarka",
			}},
			PlayersCount: 2,
			Status:       RoomStatusCreated,
		})

		err = s.Manager.RoomReady(s.Ctx, "joe biden")
		require.Error(s.T(), err)
		require.Equal(s.T(), "wrong players count", err.Error())
	})
}

func (s *RoomSuite) TestPlayerOut() {
	s.Run("Smoke", func() {
		room, err := s.DAO.Insert(s.Ctx, &Room{
			Sides: []RoomSideInfo{{
				Name: "evgsol",
			}, {
				Name: "solarka",
			}, {
				Name: "lol",
			}, {
				Name: "kek",
			}},
			Status:       RoomStatusAllPass,
			PlayersCount: 4,
		})
		require.NoError(s.T(), err)

		require.NoError(s.T(), s.Manager.PlayerOut(s.Ctx, "solarka"))

		updatedRoom, err := s.DAO.FindOneByID(s.Ctx, room.ID)
		require.NoError(s.T(), err)

		require.Equal(s.T(), EMPTY_SIDE, updatedRoom.Sides[1].Name)
		require.Equal(s.T(), 3, updatedRoom.PlayersCount)
		require.Equal(s.T(), RoomStatusCreated, updatedRoom.Status)
	})

	s.Run("Last player in room", func() {
		room, err := s.DAO.Insert(s.Ctx, &Room{
			Sides: []RoomSideInfo{{
				Name: "brad pitt",
			}},
			PlayersCount: 1,
			Status:       RoomStatusCreated,
		})
		require.NoError(s.T(), err)

		err = s.Manager.PlayerOut(s.Ctx, "brad pitt")
		require.NoError(s.T(), err)

		_, err = s.DAO.FindOneByID(s.Ctx, room.ID)
		require.Error(s.T(), err)
		require.IsType(s.T(), mgo.ErrNotFound, err)
	})

	s.Run("Player is not in room", func() {
		err := s.Manager.PlayerOut(s.Ctx, "elon musk")
		require.Error(s.T(), err)
		require.Equal(s.T(), "player is not in room", err.Error())
	})
}

func (s *RoomSuite) TestCreateRoom() {
	s.Require().NoError(s.Manager.CreateRoom(s.Ctx, "evgsol"))

	room, err := s.Manager.GetOneForPlayer(s.Ctx, "evgsol")
	s.Require().NoError(err)
	s.NotNil(room)
}

func (s *RoomSuite) TestGetAll() {
	room1, err := s.DAO.Insert(s.Ctx, &Room{})
	s.Require().NoError(err)
	room2, err := s.DAO.Insert(s.Ctx, &Room{})
	s.Require().NoError(err)

	rooms, err := s.Manager.GetAll(s.Ctx)
	s.Require().NoError(err)

	if s.Len(rooms, 2) {
		s.Equal(room1.ID.String(), rooms[0].ID)
		s.Equal(room2.ID.String(), rooms[1].ID)
	}
}
