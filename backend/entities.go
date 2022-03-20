package main

type Player struct {
	Name string
}

type RoomStatus int

const (
	RoomStatusCreated RoomStatus = 0
	RoomStatusPlaying RoomStatus = 1
)

type Suit string

const (
	SuitHearts   Suit = "H"
	SuitClubs    Suit = "C"
	SuitSpades   Suit = "S"
	SuitDiamonds Suit = "D"
)

type Card struct {
	Suit Suit   `json:"suit" bson:"suit"`
	Rank string `json:"rank" bson:"rank"`
}

var UnknownCard = Card{
	Suit: Suit("X"),
	Rank: "X",
}

type CenterCardInfo struct {
	Card   Card   `json:"card" bson:"card"`
	Player string `json:"player" bson:"player"`
}

type RoomSideInfo struct {
	Name   string `json:"name" bson:"name"`
	Cards  []Card `json:"cards" bson:"cards"`
	Tricks int    `json:"tricks" bson:"tricks"`
	Open   bool   `json:"-" bson:"open"`
}

type Room struct {
	ID     string           `json:"id" bson:"_id"`
	Sides  []RoomSideInfo   `json:"sides" bson:"sides"`
	Center []CenterCardInfo `json:"center" bson:"center"`
	Status RoomStatus       `json:"status" bson:"status"`
}
