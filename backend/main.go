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
	userDAO := NewUserDAO(session)
	userManager := NewUserManager(userDAO)
	loginManager := NewLoginManager(userManager)
	controller := NewController(roomManager)

	mux := http.NewServeMux()

	mux.Handle("/login", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(loginManager.Login)))
	mux.Handle("/register", handlers.LoggingHandler(os.Stdout, decorate(loginManager.Register)))

	mux.Handle("/room", handlers.LoggingHandler(os.Stdout, decorate(loginRequired(controller.Room))))
	mux.Handle("/shuffle", handlers.LoggingHandler(os.Stdout, decorate(loginRequired(controller.Shuffle))))
	mux.Handle("/openBuypack", handlers.LoggingHandler(os.Stdout, decorate(loginRequired(controller.OpenBuypack))))
	mux.Handle("/takeBuypack", handlers.LoggingHandler(os.Stdout, decorate(loginRequired(controller.TakeBuypack))))
	mux.Handle("/drop", handlers.LoggingHandler(os.Stdout, decorate(loginRequired(controller.Drop))))
	mux.Handle("/move", handlers.LoggingHandler(os.Stdout, decorate(loginRequired(controller.Move))))
	mux.Handle("/takeTrick", handlers.LoggingHandler(os.Stdout, decorate(loginRequired(controller.TakeTrick))))
	mux.Handle("/allPass", handlers.LoggingHandler(os.Stdout, decorate(loginRequired(controller.AllPass))))
	mux.Handle("/changeVisibility", handlers.LoggingHandler(os.Stdout, decorate(loginRequired(controller.ChangeVisibility))))

	mux.Handle("/playerIn", handlers.LoggingHandler(os.Stdout, decorate(loginRequired(controller.PlayerIn))))
	mux.Handle("/playerOut", handlers.LoggingHandler(os.Stdout, decorate(loginRequired(controller.PlayerOut))))
	mux.Handle("/roomReady", handlers.LoggingHandler(os.Stdout, decorate(loginRequired(controller.RoomReady))))
	mux.Handle("/createRoom", handlers.LoggingHandler(os.Stdout, decorate(loginRequired(controller.CreateRoom))))

	c := cors.New(cors.Options{
		AllowedOrigins:   Config.Hostnames,
		AllowCredentials: true,
		Debug:            false,
	})
	log.Fatal(http.ListenAndServe(":8090", c.Handler(mux)))
}

func decorate(f func(*http.Request) (interface{}, error)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		result, err := f(r)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		if err == nil {
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(result); err != nil {
				panic(err)
			}
		} else {
			// TODO: implement specific exceptions and mapping to status codes
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
	})
}
