package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/rs/cors"
)

func room(w http.ResponseWriter, request *http.Request, playerName string) {
	result := Room{
		Sides: []RoomSideInfo{{
			Name: "evgsol",
			Cards: []Card{{
				Suit: SuitDiamonds,
				Rank: "Q",
			}, {

				Suit: SuitSpades,
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
			}, {

				Suit: SuitDiamonds,
				Rank: "Q",
			}, {

				Suit: SuitDiamonds,
				Rank: "K",
			}, {

				Suit: SuitSpades,
				Rank: "7",
			}, {

				Suit: SuitDiamonds,
				Rank: "Q",
			}, {

				Suit: SuitDiamonds,
				Rank: "K",
			}, {

				Suit: SuitSpades,
				Rank: "7",
			}, {

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
			}, {

				Suit: SuitClubs,
				Rank: "Q",
			}, {

				Suit: SuitClubs,
				Rank: "K",
			}, {

				Suit: SuitClubs,
				Rank: "A",
			}, {

				Suit: SuitClubs,
				Rank: "Q",
			}, {

				Suit: SuitClubs,
				Rank: "K",
			}, {

				Suit: SuitClubs,
				Rank: "A",
			}, {

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

func main() {
	mux := http.NewServeMux()

	mux.Handle("/login", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(login)))
	mux.Handle("/room", handlers.LoggingHandler(os.Stdout, loginRequired(room)))

	c := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://0.0.0.0", "http://0.0.0.0:8080",
			"http://52.91.188.222", "https://52.91.188.222",
		},
		AllowCredentials: true,
		Debug:            false,
	})
	log.Fatal(http.ListenAndServe(":8090", c.Handler(mux)))
}
