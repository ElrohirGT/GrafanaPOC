package main

import (
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

func rolldice(w http.ResponseWriter, r *http.Request) {
	roll := 1 + rand.Intn(6)

	var msg string
	if player := r.PathValue("player"); player != "" {
		msg = player + " is rolling the dice"
	} else {
		msg = "Anonymous player is rolling the dice"
	}
	log.Println(msg)

	log.Println("Getting first user...")
	err := GetUser(r.Context())
	if err != nil {
		log.Println("Failed to get user!", "error", err)
		http.Error(w, "Failed to get user! "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Posting a custom obj...")
	err = PostObj(r.Context())
	if err != nil {
		log.Println("Failed to post obj!", "error", err)
		http.Error(w, "Failed to post obj! "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := strconv.Itoa(roll) + "\n"
	if _, err := io.WriteString(w, resp); err != nil {
		log.Println("Write failed", "error", err)
	}
}
