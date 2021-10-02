package main

import "github.com/bwmarrin/discordgo"

// command is a wrapper for a slash command and its handler.
type command struct {
	appCommand *discordgo.ApplicationCommand
	handler    func(*discordgo.Session, *discordgo.InteractionCreate)
}

// commandType is a type alias.
type commandOption = discordgo.ApplicationCommandInteractionDataOption

// All commands are specified here.
var commands = []*command{
	{
		appCommand: &commandConvert,
		handler:    handleConvert,
	},
	{
		appCommand: &commandHelp,
		handler:    handleHelp,
	},
	{
		appCommand: &commandInvite,
		handler:    handleInvite,
	},
	{
		appCommand: &commandPing,
		handler:    handlePing,
	},
	{
		appCommand: &commandStats,
		handler:    handleStats,
	},
}

// Invokes a subcommand with the arguments if it matches the given name.
func commandWhen(o []*commandOption, s string, do func([]*commandOption)) {
	for _, opt := range o {
		if opt.Name == s {
			do(opt.Options)
			return
		}
	}
}

// Invokes a function with an option if it matches the given name.
func commandGet1(o []*commandOption, s string, do func(*commandOption)) {
	do(commandGet2(o, s))
}

// Gets an option value if it matches the given name.
func commandGet2(o []*commandOption, s string) *commandOption {
	for _, opt := range o {
		if opt.Name == s {
			return opt
		}
	}
	return nil
}
