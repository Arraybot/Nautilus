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
		appCommand: &commandConvert,
		handler:    handleConvert,
	},
	{
		appCommand: &commandPing,
		handler:    handlePing,
	},
}

func commandWhen(o []*discordgo.ApplicationCommandInteractionDataOption, s string, do func([]*discordgo.ApplicationCommandInteractionDataOption)) {
	for _, opt := range o {
		if opt.Name == s {
			do(opt.Options)
			return
		}
	}
}

func commandGet1(o []*discordgo.ApplicationCommandInteractionDataOption, s string, do func(*discordgo.ApplicationCommandInteractionDataOption)) {
	do(commandGet2(o, s))
}

func commandGet2(o []*discordgo.ApplicationCommandInteractionDataOption, s string) *discordgo.ApplicationCommandInteractionDataOption {
	for _, opt := range o {
		if opt.Name == s {
			return opt
		}
	}
	return nil
}
