package main

import "github.com/bwmarrin/discordgo"

// The ping command.
var commandPing = discordgo.ApplicationCommand{
	Name:        "ping",
	Description: "Checks if the bot is online.",
}

// The ping command.
// Will just respond with a message.
var handlePing = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, respondText("Pong! Command handler online and responsive.", true))
}
