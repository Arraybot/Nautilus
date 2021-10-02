package main

import (
	"fmt"
	"regexp"
	"runtime"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

// The convert command.
var commandConvert = discordgo.ApplicationCommand{
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
}

// The help command.
var commandHelp = discordgo.ApplicationCommand{
	Name:        "help",
	Description: "Shows helpful information on the bot.",
}

// The invite command.
var commandInvite = discordgo.ApplicationCommand{
	Name:        "invite",
	Description: "Shows an invite link for the server and bot.",
}

// The ping command.
var commandPing = discordgo.ApplicationCommand{
	Name:        "ping",
	Description: "Checks if the bot is online.",
}

// The stats command.
var commandStats = discordgo.ApplicationCommand{
	Name:        "stats",
	Description: "Shows bot usage and technical statistics.",
}

// The convert command.
// Will respond with a fancy embed.
var handleConvert = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	o := i.ApplicationCommandData().Options
	commandWhen(o, "hex", func(o []*commandOption) {
		// These values will strategically underflow and cause values > 255.
		r := uint32(commandGet2(o, "red").IntValue())
		g := uint32(commandGet2(o, "green").IntValue())
		b := uint32(commandGet2(o, "blue").IntValue())
		// Check if the input range is valid.
		if r > 255 || g > 255 || b > 255 {
			s.InteractionRespond(i.Interaction, respondText("Invalid values for at least one colour (out of range, [0-255]).", true))
			return
		}
		result := fmt.Sprintf("#%x", rgbToHex(r, g, b))
		embed := embed()
		embed.field("Result", result, false)
		s.InteractionRespond(i.Interaction, respondEmbed(embed, true))
	})
	commandWhen(o, "rgb", func(o []*commandOption) {
		commandGet1(o, "colour", func(o *commandOption) {
			raw := o.StringValue()
			// Ensure it matches the regular expression.
			if !hexMatchRegex.MatchString(raw) {
				s.InteractionRespond(i.Interaction, respondText("Invalid hex code provided (format unknown).", true))
				return
			}
			// Extract all the non-hex characters.
			parsed := hexReplaceRegex.ReplaceAllString(raw, "")
			hex, err := strconv.ParseUint(parsed, 16, 32)
			if err != nil {
				s.InteractionRespond(i.Interaction, respondText("Invalid hex code provided (out of range).", true))
				return
			}
			r, g, b := hexToRgb(hex)
			// Format into nice output.
			result := fmt.Sprintf("R: %d, G: %d, B: %d", r, g, b)
			embed := embed()
			embed.field("Result", result, false)
			s.InteractionRespond(i.Interaction, respondEmbed(embed, true))
		})
	})
}

// The help command.
// Will display an embed with some useful information.
var handleHelp = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	embed := embed()
	description := "Arraybot is a multipurpose toolbox designed to run your guild. " +
		"Since 2016, Arraybot has powered some of Discord's biggest servers. " +
		"Arraybot's goals are reliability, flexibility, acccessibility.\n" + zwsp
	startUser := "For users, Arraybot is mainly command-based.\n" +
		"**1.** Execute the `/guide` command to show the current server's guide.\n" +
		"**2.** All commands can be found by opening the Discord command menu and selecting Arraybot.\n" + zwsp
	startAdmin := "Configuring Arraybot is done on the web panel.\n" +
		"**1.** Visit the [web panel](https://arraybot.xyz/panel/) and log in.\n" +
		"**2.** Configure the bot to your liking.\n" +
		"**3.** Consult the [documentation](https://arraybot.xyz/go/docs/) if you get stuck.\n" +
		"**4.** Discuss and provide feedback on the server, GitHub, etc.\n" + zwsp
	links := "Here are some helpful links.\n" +
		"**1.** Check out the [website](https://arraybot.xyz/).\n" +
		"**2.** Join the [server](https://arraybot.xyz/go/server/).\n" +
		"**3.** Invite the [bot](https://arraybot.xyz/go/invite/) to your server."
	embed.description(description)
	embed.field("User Guide", startUser, false)
	embed.field("Server Admin Guide", startAdmin, false)
	embed.field("Links", links, false)
	embed.Embed.Footer = &discordgo.MessageEmbedFooter{
		Text: "Arraybot is an open source project.",
	}
	s.InteractionRespond(i.Interaction, respondEmbed(embed, true))
}

// The invite command.
// Will send an embed with clickable links to invite the bot and join the server.
var handleInvite = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	embed := embed()
	embed.field("Invite Bot", "[https://arraybot.xyz/go/invite/](https://arraybot.xyz/go/invite)", false)
	embed.field("Join Server", "[https://arraybot.xyz/go/server/](https://arraybot.xyz/go/server/)", false)
	s.InteractionRespond(i.Interaction, respondEmbed(embed, true))
}

// The ping command.
// Will just respond with a message.
var handlePing = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, respondText("Pong! Command handler online and responsive.", true))
}

// The stats command.
// Will show a few statistics.
var handleStats = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	embed := embed()
	embed.description("Here are some statistics.")
	// TODO obtain actual statistics.
	statsGuilds := -1
	statsMessages := -1
	statsCommandsRun := -1
	embed.field("# of Servers", strconv.Itoa(statsGuilds), true)
	embed.field("# of Messages", strconv.Itoa(statsMessages), true)
	embed.field("# of Commands Run", strconv.Itoa(statsCommandsRun), true)
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	elapsed := time.Since(startTime).String()
	embed.field("Tot.Alloc.", fmt.Sprintf("%v MiB", bytesToMegabytes(mem.TotalAlloc)), true)
	embed.field("Sys.Alloc.", fmt.Sprintf("%v MiB", bytesToMegabytes(mem.Sys)), true)
	embed.field("Uptime", elapsed, true)
	s.InteractionRespond(i.Interaction, respondEmbed(embed, true))
}

// Helper variables.
var zwsp = "\u200b"
var hexMatchRegex = regexp.MustCompile("#?[0-9a-fA-F]{6}")
var hexReplaceRegex = regexp.MustCompile("[^0-9a-fA-F]")
