package main

type Player struct {
	Name string
}

type RoomStatus int

const (
	RoomStatusCreated       RoomStatus = 0
	RoomStatusPlaying       RoomStatus = 1
	RoomStatusBuypackOpened RoomStatus = 2
	RoomStatusBuypackTaken  RoomStatus = 3
	RoomStatusAllPass       RoomStatus = 4
)

type Suit string

const (
	SuitHearts   Suit = "H"
	SuitClubs    Suit = "C"
	SuitSpades   Suit = "S"
	SuitDiamonds Suit = "D"
)

var (
	AllSuits = []Suit{SuitSpades, SuitDiamonds, SuitClubs, SuitHearts}
	AllRanks = []string{"7", "8", "9", "10", "J", "Q", "K", "A"}
)

type Card struct {
	Suit Suit   `json:"suit" bson:"suit"`
	Rank string `json:"rank" bson:"rank"`
}

func (c Card) suitNumber() int {
	for i, s := range AllSuits {
		if c.Suit == s {
			return i
		}
	}

	return 4
}

func (c Card) rankNumber() int {
	for i, r := range AllRanks {
		if c.Rank == r {
			return i
		}
	}

	return 8
}

func (c Card) Less(other Card) bool {
	if c.suitNumber() < other.suitNumber() {
		return true
	}
	if (c.suitNumber() == other.suitNumber()) && (c.rankNumber() < other.rankNumber()) {
		return true
	}

	return false
}

var (
	UnknownCard = Card{
		Suit: Suit("X"),
		Rank: "X",
	}
	EmptySidePlayerName = "-"
)

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
	ID           string           `json:"id" bson:"_id"`
	Sides        []RoomSideInfo   `json:"sides" bson:"sides"`
	Center       []CenterCardInfo `json:"center" bson:"center"`
	LastTrick    []CenterCardInfo `json:"lastTrick" bson:"lastTrick"`
	Status       RoomStatus       `json:"status" bson:"status"`
	PlayersCount int              `json:"playersCount" bson:"playersCount"`
	BuypackIndex int              `json:"buypackIndex" bson:"buypackIndex"`
}

func (r *Room) PlayerSideIndex(playerName string) int {
	for i, s := range r.Sides {
		if s.Name == playerName {
			return i
		}
	}

	return -1
}
