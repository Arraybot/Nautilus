package main

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/arraybot/nautilus/commands"
	"github.com/bsdlp/discord-interactions-go/interactions"
	"github.com/bwmarrin/discordgo"
)

// Handles incoming (HTTP) interactions.
func slashHandler(w http.ResponseWriter, req *http.Request) {
	// First, get the public key.
	pubKey := os.Getenv("APP_PUBKEY")
	keyBytes, err := hex.DecodeString(pubKey)
	if err != nil {
		endpointError(w, err)
		return
	}
	// Check if the signature matches.
	verified := interactions.Verify(req, ed25519.PublicKey(keyBytes))
	if !verified {
		log.Println("Received invalid command request")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// Read the body.
	defer req.Body.Close()
	var data discordgo.InteractionCreate
	err = json.NewDecoder(req.Body).Decode(&data)
	if err != nil {
		endpointError(w, err)
		return
	}
	// Check if it's a ping, if so, just respond.
	if data.Type == discordgo.InteractionPing {
		_, err := w.Write([]byte(`{"type":1}`))
		if err != nil {
			endpointError(w, err)
		}
		log.Println("Responded to Discord ping.")
		return
	}
	// If not a ping, it is a command. Handle appropriately.
	commands.Distributor(client, &data)
}

// Helper method to write a server error in the HTTP response.
func endpointError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	log.Println(err)
}
