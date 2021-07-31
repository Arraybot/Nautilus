package main

import "github.com/bwmarrin/discordgo"

var commandPing = discordgo.ApplicationCommand{
	Name:        "ping",
	Description: "Checks if the bot is online.",
}

var handlePing = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, respondText("Pong! Command handler online and responsive.", true))
}
