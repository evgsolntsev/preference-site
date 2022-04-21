package main

import (
	"encoding/json"
	"errors"
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

func (c *Controller) Room(request *http.Request, playerName string) (interface{}, error) {
	result, err := c.roomManager.GetOneForPlayer(request.Context(), playerName)
	if err != nil {
		return nil, err
	}

	for i, _ := range result.Sides {
		if result.Sides[i].Name == playerName || result.Sides[i].Open {
			continue
		}
		for j, _ := range result.Sides[i].Cards {
			result.Sides[i].Cards[j] = UnknownCard
		}
	}

	return result, nil
}

func (c *Controller) Shuffle(request *http.Request, playerName string) (interface{}, error) {
	room, err := c.roomManager.GetOneForPlayer(request.Context(), playerName)
	if err != nil {
		return nil, err
	}

	if err := c.roomManager.Shuffle(request.Context(), room.ID, playerName); err != nil {
		return nil, err
	}

	return nil, nil
}

func (c *Controller) OpenBuypack(request *http.Request, playerName string) (interface{}, error) {
	room, err := c.roomManager.GetOneForPlayer(request.Context(), playerName)
	if err != nil {
		return nil, err
	}

	if err := c.roomManager.OpenBuypack(request.Context(), room.ID); err != nil {
		return nil, err
	}

	return nil, nil
}

func (c *Controller) TakeBuypack(request *http.Request, playerName string) (interface{}, error) {
	room, err := c.roomManager.GetOneForPlayer(request.Context(), playerName)
	if err != nil {
		return nil, err
	}

	if err := c.roomManager.TakeBuypack(request.Context(), room.ID, playerName); err != nil {
		return nil, err
	}

	return nil, nil
}

type Indexes struct {
	Indexes []int `json:"indexes"`
}

func (c *Controller) Drop(request *http.Request, playerName string) (interface{}, error) {
	var indexes Indexes
	if err := json.NewDecoder(request.Body).Decode(&indexes); err != nil {
		return nil, errors.New("wrond number of indexes")
	}

	room, err := c.roomManager.GetOneForPlayer(request.Context(), playerName)
	if err != nil {
		return nil, err
	}

	if err := c.roomManager.Drop(request.Context(), room.ID, playerName, indexes.Indexes); err != nil {
		return nil, err
	}

	return nil, nil
}
