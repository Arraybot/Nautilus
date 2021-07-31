package main

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/bsdlp/discord-interactions-go/interactions"
	"github.com/bwmarrin/discordgo"
)

func slashHandler(w http.ResponseWriter, req *http.Request) {
	pubKey := os.Getenv("PUBKEY")
	keyBytes, err := hex.DecodeString(pubKey)
	if err != nil {
		slashEndpointError(w, err)
		return
	}
	verified := interactions.Verify(req, ed25519.PublicKey(keyBytes))
	if !verified {
		log.Println("Received invalid command request")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	defer req.Body.Close()
	var data discordgo.InteractionCreate
	err = json.NewDecoder(req.Body).Decode(&data)
	if err != nil {
		slashEndpointError(w, err)
		return
	}
	if data.Type == discordgo.InteractionPing {
		_, err := w.Write([]byte(`{"type":1}`))
		if err != nil {
			slashEndpointError(w, err)
		}
		log.Println("Responded to Discord ping.")
		return
	}
	slashDistributor(client, &data)
}

func slashDistributor(s *discordgo.Session, i *discordgo.InteractionCreate) {
	name := i.ApplicationCommandData().Name
	for _, command := range commands {
		if command.appCommand.Name == name {
			log.Printf("User %s invoked command %s in %s.\n", i.Interaction.Member.User.ID, name, i.Interaction.GuildID)
			command.handler(s, i)
		}
	}
}

func slashEndpointError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	log.Println(err)
}
