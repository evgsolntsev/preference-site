package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRoomHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/room", nil)
	w := httptest.NewRecorder()
	room(w, req)
	res := w.Result()
	defer res.Body.Close()
	expected := `{
            "sides": [
                {
                    "name": "evgsol",
                    "cards": [
                        {
                            "card": {
                                "suit": 9830,
                                "rank": "Q"
                            },
                            "open": true
                        },
                        {
                            "card": {
                                "suit": 9830,
                                "rank": "K"
                            },
                            "open": true
                        },
                        {
                            "card": {
                                "suit": 9824,
                                "rank": "7"
                            },
                            "open": true
                        }
                    ],
                    "tricks": 0
                },
                {
                    "name": "miracle",
                    "cards": [
                        {
                            "card": {
                                "suit": 9830,
                                "rank": "J"
                            },
                            "open": true
                        },
                        {
                            "card": {
                                "suit": 9830,
                                "rank": "A"
                            },
                            "open": true
                        },
                        {
                            "card": {
                                "suit": 9827,
                                "rank": "J"
                            },
                            "open": true
                        }
                    ],
                    "tricks": 0
                },
                {
                    "name": "solarka",
                    "cards": null,
                    "tricks": 0
                },
                {
                    "name": "psmirnov",
                    "cards": [
                        {
                            "card": {
                                "suit": 9827,
                                "rank": "Q"
                            },
                            "open": false
                        },
                        {
                            "card": {
                                "suit": 9827,
                                "rank": "K"
                            },
                            "open": false
                        },
                        {
                            "card": {
                                "suit": 9827,
                                "rank": "A"
                            },
                            "open": false
                        }
                    ],
                    "tricks": 0
                }
            ],
            "center": [
                {},
                {}
            ],
            "status": 1
        }`
	data, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	assert.JSONEq(t, expected, string(data))
}
