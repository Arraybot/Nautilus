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
		appCommand: &discordgo.ApplicationCommand{
			Name:        "convert",
			Description: "Converts various values.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "hex",
					Description: "Convert an RGB colour to hexadecimal.",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "red",
							Description: "The red value.",
							Type:        discordgo.ApplicationCommandOptionInteger,
							Required:    true,
						},
						{
							Name:        "green",
							Description: "The green value.",
							Type:        discordgo.ApplicationCommandOptionInteger,
							Required:    true,
						},
						{
							Name:        "blue",
							Description: "The blue value.",
							Type:        discordgo.ApplicationCommandOptionInteger,
							Required:    true,
						},
					},
				},
				{
					Name:        "rgb",
					Description: "Convert a hexadecimal colour to RGB.",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "colour",
							Description: "The colour in hexadecimal.",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
						},
					},
				},
			},
		},
		handler: handleConvert,
	},
	{
		appCommand: &discordgo.ApplicationCommand{
			Name:        "help",
			Description: "Shows helpful information on the bot.",
		},
		handler: handleHelp,
	},
	{
		appCommand: &discordgo.ApplicationCommand{
			Name:        "invite",
			Description: "Shows an invite link for the server and bot.",
		},
		handler: handleInvite,
	},
	{
		appCommand: &discordgo.ApplicationCommand{
			Name:        "ping",
			Description: "Checks if the bot is online.",
		},
		handler: handlePing,
	},
	{
		appCommand: &discordgo.ApplicationCommand{
			Name:        "stats",
			Description: "Shows bot usage and technical statistics.",
		},
		handler: handleStats,
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
