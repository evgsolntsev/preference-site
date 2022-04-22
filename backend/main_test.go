package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/globalsign/mgo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRoomHandler(t *testing.T) {
	session, err := mgo.Dial("mongodb://localhost:27017")
	require.NoError(t, err)

	ctx := context.Background()
	dao := NewRoomDAO(session)
	defer dao.RemoveAll(ctx)

	r := &Room{
		Sides: []RoomSideInfo{{
			Name: "evgsol",
			Cards: []Card{{
				Suit: SuitDiamonds,
				Rank: "Q",
			}, {
				Suit: SuitDiamonds,
				Rank: "A",
			}, {
				Suit: SuitSpades,
				Rank: "7",
			}},
		}, {
			Name: "miracle",
			Cards: []Card{{
				Suit: SuitDiamonds,
				Rank: "J",
			}, {
				Suit: SuitDiamonds,
				Rank: "A",
			}, {
				Suit: SuitClubs,
				Rank: "J",
			}},
			Open: true,
		}, {
			Name: "solarka",
			Cards: []Card{{
				Suit: SuitDiamonds,
				Rank: "Q",
			}, {
				Suit: SuitDiamonds,
				Rank: "K",
			}, {
				Suit: SuitSpades,
				Rank: "7",
			}},
		}, {
			Name: "psmirnov",
			Cards: []Card{{
				Suit: SuitClubs,
				Rank: "Q",
			}, {
				Suit: SuitClubs,
				Rank: "K",
			}, {
				Suit: SuitClubs,
				Rank: "A",
			}},
		}},
		Center: []CenterCardInfo{{
			Card: Card{
				Suit: SuitSpades,
				Rank: "A",
			},
			Player: "evgsol",
		}, {
			Card: Card{
				Suit: SuitHearts,
				Rank: "10",
			},
			Player: "miracle",
		}},
		Status: RoomStatusPlaying,
	}
	stored, err := dao.Insert(ctx, r)
	require.NoError(t, err)

	manager := NewRoomManager(dao)
	handler := NewController(manager)

	req := httptest.NewRequest(http.MethodGet, "/room", nil)
	result, err := handler.Room(req, "evgsol")
	require.NoError(t, err)

	expected := fmt.Sprintf(`{
            "id": "%v",
            "sides": [
                {
                    "name": "evgsol",
                    "cards": [{
                        "suit": "D",
                        "rank": "Q"
                    }, {
                        "suit": "D",
                        "rank": "A"
                    }, {
                        "suit": "S",
                        "rank": "7"
                    }],
                    "tricks": 0
                },
                {
                    "name": "miracle",
                    "cards": [{
                        "suit": "D",
                        "rank": "J"
                    },
                    {
                        "suit": "D",
                        "rank": "A"
                    },
                    {
                        "suit": "C",
                        "rank": "J"
                    }],
                    "tricks": 0
                },
                {
                    "name": "solarka",
                    "cards": [{
                        "suit": "X",
                        "rank": "X"
                    },
                    {
                        "suit": "X",
                        "rank": "X"
                    },
                    {
                        "suit": "X",
                        "rank": "X"
                    }],
                    "tricks": 0
                },
                {
                    "name": "psmirnov",
                    "cards": [{
                        "suit": "X",
                        "rank": "X"
                    },
                    {
                        "suit": "X",
                        "rank": "X"
                    },
                    {
                        "suit": "X",
                        "rank": "X"
                    }],
                    "tricks": 0
                }
            ],
            "center": [{
                "player": "evgsol",
                "card": {"suit": "S", "rank": "A"}
            }, {
                "player": "miracle",
                "card": {"suit": "H", "rank": "10"}
            }],
            "status": 1,
            "lastTrick": [],
            "playersCount": 0
        }`, stored.ID)

	res, err := json.Marshal(result)
	require.NoError(t, err)

	assert.JSONEq(t, expected, string(res))
}
