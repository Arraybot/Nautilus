package commands

import (
	"log"
	"os"
	"strings"

	"github.com/arraybot/nautilus/database"
	"github.com/bwmarrin/discordgo"
)

// command is a wrapper for a slash command and its handler.
type command struct {
	appCommand *discordgo.ApplicationCommand
	handler    func(*discordgo.Session, *discordgo.InteractionCreate)
}

// commandType is a type alias.
type commandOption = discordgo.ApplicationCommandInteractionDataOption

var admins = strings.Split(os.Getenv("ADMINS"), ";")

// All commands are specified here.
var commands = []*command{
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
	// Fun commands.
	{
		appCommand: &discordgo.ApplicationCommand{
			Name:        "cat",
			Description: "Sends a random cat image/GIF/video.",
		},
		handler: handleCat,
	},
	{
		appCommand: &discordgo.ApplicationCommand{
			Name:        "dog",
			Description: "Sends a random dog image/GIF/video.",
		},
		handler: handleDog,
	},
	{
		appCommand: &discordgo.ApplicationCommand{
			Name:        "8ball",
			Description: "Decides your fate.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "question",
					Description: "A yes/no question you would like to ask.",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    true,
				},
			},
		},
		handler: handleEightBall,
	},
	{
		appCommand: &discordgo.ApplicationCommand{
			Name:        "urban",
			Description: "Looks up a phrase in the Urban Dictionary.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "phrase",
					Description: "The phrase to look up.",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    true,
				},
			},
		},
		handler: handleUrban,
	},
	// Moderator commands.
	{
		appCommand: &discordgo.ApplicationCommand{
			Name:        "clear",
			Description: "Deletes the last X messages where you specify X.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "amount",
					Description: "The amount to delete.",
					Type:        discordgo.ApplicationCommandOptionInteger,
					Required:    true,
				},
				{
					Name:        "user",
					Description: "Only delete messages by this user.",
					Type:        discordgo.ApplicationCommandOptionUser,
					Required:    false,
				},
			},
		},
		handler: handlerClear,
	},
	{
		appCommand: &discordgo.ApplicationCommand{
			Name:        "mute",
			Description: "Adds the muted role to a user.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "user",
					Description: "The user to mute.",
					Type:        discordgo.ApplicationCommandOptionUser,
					Required:    true,
				},
			},
		},
		handler: handlerMute,
	},
	{
		appCommand: &discordgo.ApplicationCommand{
			Name:        "unmute",
			Description: "Removes the muted role from a user.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "user",
					Description: "The user to unmute.",
					Type:        discordgo.ApplicationCommandOptionUser,
					Required:    true,
				},
			},
		},
		handler: handlerUnMute,
	},
	{
		appCommand: &discordgo.ApplicationCommand{
			Name:        "expire",
			Description: "Expires a punishment at a certain time in the future.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "case",
					Description: "The case ID.",
					Type:        discordgo.ApplicationCommandOptionInteger,
					Required:    true,
				},
				{
					Name:        "minutes",
					Description: "How many minutes from now.",
					Type:        discordgo.ApplicationCommandOptionInteger,
					Required:    true,
				},
				{
					Name:        "hours",
					Description: "How many hours from now.",
					Type:        discordgo.ApplicationCommandOptionInteger,
					Required:    false,
				},
				{
					Name:        "days",
					Description: "How many days from now.",
					Type:        discordgo.ApplicationCommandOptionInteger,
					Required:    false,
				},
				{
					Name:        "weeks",
					Description: "How many weeks from now.",
					Type:        discordgo.ApplicationCommandOptionInteger,
					Required:    false,
				},
				{
					Name:        "months",
					Description: "How many months from now.",
					Type:        discordgo.ApplicationCommandOptionInteger,
					Required:    false,
				},
			},
		},
		handler: handlerExpire,
	},
	{
		appCommand: &discordgo.ApplicationCommand{
			Name:        "revoke",
			Description: "Revokes a punishment immediately.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "case",
					Description: "The case ID.",
					Type:        discordgo.ApplicationCommandOptionInteger,
					Required:    true,
				},
			},
		},
		handler: handlerRevoke,
	},
	{
		appCommand: &discordgo.ApplicationCommand{
			Name:        "lookup",
			Description: "Looks up information on a certain case.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "case",
					Description: "The case ID.",
					Type:        discordgo.ApplicationCommandOptionInteger,
					Required:    true,
				},
			},
		},
		handler: handlerLookup,
	},
	{
		appCommand: &discordgo.ApplicationCommand{
			Name:        "history",
			Description: "Gets the punishment history of a user.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "user",
					Description: "The user to get the history for.",
					Type:        discordgo.ApplicationCommandOptionUser,
					Required:    true,
				},
			},
		},
		handler: handlerHistory,
	},
	// Server commands.
	{
		appCommand: &discordgo.ApplicationCommand{
			Name:        "guide",
			Description: "Shows instructions on how to use the server and/or bot.",
		},
		handler: handlerGuide,
	},
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
}

// Register registers all commands.
func Register(registrar func(*discordgo.ApplicationCommand) error) {
	for _, command := range commands {
		log.Printf("Registering command %s\n", command.appCommand.Name)
		if err := registrar(command.appCommand); err != nil {
			log.Println(err)
		}
	}
}

// Distributor is a common function between REST and WebSocket.
// This will call the correct command handler corresponding to the command name.
func Distributor(s *discordgo.Session, i *discordgo.InteractionCreate) {
	name := i.ApplicationCommandData().Name
	for _, command := range commands {
		if command.appCommand.Name == name {
			log.Printf("User %s invoked command %s in %s.\n", i.Interaction.Member.User.ID, name, i.Interaction.GuildID)
			if !disabled(i.GuildID, name) {
				command.handler(s, i)
			} else {
				s.InteractionRespond(i.Interaction, respondText("This command has been disabled by the/a server administrator(s).", i))
			}
		}
	}
}

// Whether a command is disabled.
func disabled(server, name string) bool {
	for _, disabled := range database.Disabled(server) {
		if disabled == name {
			return true
		}
	}
	return false
}

// Invokes a subcommand with the arguments if it matches the given name.
func sub(o []*commandOption, s string, do func([]*commandOption)) {
	for _, opt := range o {
		if opt.Name == s {
			do(opt.Options)
			return
		}
	}
}

// Gets an option value if it matches the given name.
func option(o []*commandOption, s string) *commandOption {
	for _, opt := range o {
		if opt.Name == s {
			return opt
		}
	}
	return nil
}

// Whether or not the command executor has a Discord permission.
func hasPermission(i *discordgo.InteractionCreate, p int) bool {
	perm := int64(p)
	return i.Member.Permissions&perm == perm
}

// Whether or not the command executor has permission to execute developer commands.
func hasDeveloper(i *discordgo.InteractionCreate) bool {
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

// Whether or not the command executor has any (kick, timeout, mute, ban) moderator permission.
func hasModerator(i *discordgo.InteractionCreate) bool {
	k := discordgo.PermissionKickMembers
	t := discordgo.PermissionModerateMembers
	m := discordgo.PermissionBanMembers
	return hasPermission(i, k) || hasPermission(i, t) || hasPermission(i, m) || hasMutePermission(i)
}

// Whether or not the command executor has the muted permission role.
func hasMutePermission(i *discordgo.InteractionCreate) bool {
	mod := database.MutePermission(i.GuildID)
	for _, role := range i.Member.Roles {
		if role == mod {
			return true
		}
	}
	return false
}

// Helper variables.
var permissionDenyPermission = "You do not have the required Discord permission to execute this command."
var permissionDenyDeveloper = "You need to be an Arraybot authorized developer to execute this command."
var permissionDenyModerator = "You need to have moderator permissions on this server to execute this command."
var permissionDenyMuteRole = "You need be given a moderative role to execute this command."
