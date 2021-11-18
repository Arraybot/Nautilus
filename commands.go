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
	// Utility commands.
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
	// Developer commands.
	{
		appCommand: &discordgo.ApplicationCommand{
			Name:        "kill",
			Description: "Kills and restarts specific services/shards.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "services",
					Description: "Kills and restarts a specific Arraybot service.",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "name",
							Description: "The name of the service.",
							Type:        discordgo.ApplicationCommandOptionString,
							Choices: []*discordgo.ApplicationCommandOptionChoice{
								{
									Name:  "Carbon (Web Panel)",
									Value: "carbon",
								},
								{
									Name:  "Mantis (Gateway Handler)",
									Value: "mantis",
								},
								{
									Name:  "Nautilus (Command Engine)",
									Value: "nautilus",
								},
							},
							Required: true,
						},
					},
				},
				{
					Name:        "shard",
					Description: "Kills and restarts a specific gateway shard.",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "shard",
							Description: "The ID of the shard.",
							Type:        discordgo.ApplicationCommandOptionInteger,
							Required:    true,
						},
					},
				},
			},
		},
		handler: handleKill,
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

// Whether or not the command executor has permission to execute developer commands.
func commandPermissionDeveloper(i *discordgo.InteractionCreate) bool {
	for _, a := range admins {
		if i.User != nil && i.User.ID == a {
			return true
		}
		if i.Member != nil && i.Member.User.ID == a {
			return true
		}
	}
	return false
}

// Helper variables.
var permissionDenyDeveloper = "You need to be an Arraybot authorized developer to execute this command."
var permissionDenyModerator = "You need to be set as a server moderator to execute this command."
var permissionDenyPermission = "You do not have the required Discord permission to execute this command."
