package main

import (
	"encoding/json"
	"net/http"
)

type Controller struct {
	roomManager *RoomManager
}

func NewController(m *RoomManager) *Controller {
	return &Controller{
		roomManager: m,
	}
}

func (c *Controller) Room(w http.ResponseWriter, request *http.Request, playerName string) {
	result, err := c.roomManager.GetOneForPlayer(request.Context(), playerName)
	if err != nil {
		panic(err)
	}

	for i, _ := range result.Sides {
		if result.Sides[i].Name == playerName || result.Sides[i].Open {
			continue
		}
		for j, _ := range result.Sides[i].Cards {
			result.Sides[i].Cards[j] = UnknownCard
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		panic(err)
	}
}

func (c *Controller) Shuffle(w http.ResponseWriter, request *http.Request, playerName string) {
	if err := c.roomManager.Shuffle(request.Context(), playerName); err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
}
