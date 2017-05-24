package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

type moveResponse struct {
	Valid bool
}

type newGameResponse struct {
	Full   bool
	Player int
}

func newGame(w http.ResponseWriter, r *http.Request) {
	numOfPlayer := len(b.subscribers)
	var response = newGameResponse{
		true,
		0,
	}
	if numOfPlayer == 0 {
		response.Full = false
		response.Player = 1
	} else if numOfPlayer == 1 {
		response.Full = false
		response.Player = 2
		var eventData = event{
			"game",
			startData{
				1,
			},
		}
		defer publishEvent(eventData)
		setupGame()
	}
	printResponse(w, response)
}

func move(w http.ResponseWriter, r *http.Request) {
	var move moveData
	var response = moveResponse{
		true,
	}
	vars := mux.Vars(r)
	player, err := strconv.Atoi(vars["player"])
	if err != nil {
		log.Fatal(err)
	}
	if player > 0 && player == theGame.turn {

		move.Player = player
		err = r.ParseForm()
		if err != nil {
			log.Fatal(err)
		}

		decoder := schema.NewDecoder()
		decoder.Decode(&move, r.PostForm)

		if !movePiece(move) {
			response.Valid = false
		} else {
			move.Turn = theGame.turn
			var eventData = event{
				"move",
				move,
			}
			defer publishEvent(eventData)
		}

	} else {
		response.Valid = false
	}
	printResponse(w, response)
}

func events(w http.ResponseWriter, r *http.Request) {
	f := w.(http.Flusher)
	ch := Subscribe()
	defer Unsubscribe(ch)

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	f.Flush()

	cn := w.(http.CloseNotifier)

	for {
		select {
		case m := <-ch:
			msg := fmt.Sprintf("data: %s\n\n", m)
			fmt.Fprintln(w, msg)
			f.Flush()
		case <-cn.CloseNotify():
			fmt.Println("Connection Close")
			return
		}
	}
}

func printResponse(w http.ResponseWriter, response interface{}) {
	eventJSON, err := json.Marshal(response)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprint(w, string(eventJSON))
}
