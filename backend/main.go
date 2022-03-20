package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/rs/cors"
	"github.com/globalsign/mgo"
)

type Handler struct {
	roomManager *RoomManager
}

func NewHandler(m *RoomManager) *Handler {
	return &Handler{
		roomManager: m,
	}
}

func (h *Handler) Room(w http.ResponseWriter, request *http.Request, playerName string) {
	result, err := h.roomManager.GetOneForPlayer(request.Context(), playerName)
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

var (
	CONFIGFILE = "conf.json"
	Config     Configuration
)

func main() {
	if err := Config.Init(CONFIGFILE); err != nil {
		log.Fatal(err)
	}

	session, err := mgo.Dial(Config.MongoURL)
	if err != nil {
		log.Fatal(err)
	}

	roomDAO := NewRoomDAO(session)
	roomManager := NewRoomManager(roomDAO)
	handler := NewHandler(roomManager)

	mux := http.NewServeMux()

	mux.Handle("/login", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(login)))
	mux.Handle("/room", handlers.LoggingHandler(os.Stdout, loginRequired(handler.Room)))

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
