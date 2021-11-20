package main

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

// The cat command.
// Sends a link to a cute cat.
func handleCat(s *discordgo.Session, i *discordgo.InteractionCreate) {
	cat := cat{}
	link, err := requestPet(&cat)
	var toSend string
	if err == nil {
		toSend = link
	} else {
		toSend = "Error getting cat."
		log.Println(err)
	}
	s.InteractionRespond(i.Interaction, respondText(toSend, i))
}

// The dog commands.
// Sends a link to a cute dog.
func handleDog(s *discordgo.Session, i *discordgo.InteractionCreate) {
	dog := dog{}
	link, err := requestPet(&dog)
	var toSend string
	if err == nil {
		toSend = link
	} else {
		toSend = "Error getting dog."
		log.Println(err)
	}
	s.InteractionRespond(i.Interaction, respondText(toSend, i))
}
