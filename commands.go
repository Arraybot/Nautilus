package main

import "github.com/bwmarrin/discordgo"

// command is a wrapper for a slash command and its handler.
type command struct {
	appCommand *discordgo.ApplicationCommand
	handler    func(*discordgo.Session, *discordgo.InteractionCreate)
}

// All commands are specified here.
var commands = []*command{
	{
		appCommand: &commandPing,
		handler:    handlePing,
	},
}
