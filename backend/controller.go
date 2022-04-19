package main

import (
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
