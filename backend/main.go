package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/globalsign/mgo"
	"github.com/gorilla/handlers"
	"github.com/rs/cors"
)

var (
	CONFIGFILE = "conf.json"
	Config     Configuration
)

func main() {
	if err := Config.Init(CONFIGFILE); err != nil {
		log.Fatal(err)
	}

	rand.Seed(time.Now().UnixNano())

	session, err := mgo.Dial(Config.MongoURL)
	if err != nil {
		log.Fatal(err)
	}

	roomDAO := NewRoomDAO(session)
	roomManager := NewRoomManager(roomDAO)
	controller := NewController(roomManager)

	mux := http.NewServeMux()

	mux.Handle("/login", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(login)))
	mux.Handle("/room", handlers.LoggingHandler(os.Stdout, loginRequired(decorate(controller.Room))))
	mux.Handle("/shuffle", handlers.LoggingHandler(os.Stdout, loginRequired(decorate(controller.Shuffle))))
	mux.Handle("/openBuypack", handlers.LoggingHandler(os.Stdout, loginRequired(decorate(controller.OpenBuypack))))
	mux.Handle("/takeBuypack", handlers.LoggingHandler(os.Stdout, loginRequired(decorate(controller.TakeBuypack))))
	mux.Handle("/drop", handlers.LoggingHandler(os.Stdout, loginRequired(decorate(controller.Drop))))

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

func decorate(f func(*http.Request, string) (interface{}, error)) func(http.ResponseWriter, *http.Request, string) {
	return func(w http.ResponseWriter, r *http.Request, s string) {
		result, err := f(r, s)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		if err == nil {
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(result); err != nil {
				panic(err)
			}
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
	}
}
