package main

import "github.com/bwmarrin/discordgo"

type command struct {
	appCommand *discordgo.ApplicationCommand
	handler    func(*discordgo.Session, *discordgo.InteractionCreate)
}

var commands = []*command{
	{
		appCommand: &commandPing,
		handler:    handlePing,
	},
}
