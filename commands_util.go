package main

import (
	"fmt"
	"regexp"
	"strconv"

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

// The ping command.
var commandPing = discordgo.ApplicationCommand{
	Name:        "ping",
	Description: "Checks if the bot is online.",
}

// The convert command.
// Will respond with a fancy embed.
var handleConvert = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	o := i.ApplicationCommandData().Options
	commandWhen(o, "rgb", func(o []*discordgo.ApplicationCommandInteractionDataOption) {
		// These values will strategically underflow and cause values > 255.
		r := uint32(commandGet2(o, "red").IntValue())
		g := uint32(commandGet2(o, "green").IntValue())
		b := uint32(commandGet2(o, "blue").IntValue())
		// Check if the input range is valid.
		if r > 255 || g > 255 || b > 255 {
			s.InteractionRespond(i.Interaction, respondText("Invalid values for at least one colour (out of range, [0-255]).", true))
			return
		}
		// Red is MS 8 bits.
		val := (r << 16)
		// Green is middle 8 bits.
		val |= (g << 8)
		// Blue is LS 8 bits.
		val |= b
		// Format into a hex string, trims leading 0s.
		result := fmt.Sprintf("#%x", val)
		embed := embed()
		embed.field("Result", result, false)
		s.InteractionRespond(i.Interaction, respondEmbed(embed, true))
	})
	commandWhen(o, "hex", func(o []*discordgo.ApplicationCommandInteractionDataOption) {
		commandGet1(o, "colour", func(o *discordgo.ApplicationCommandInteractionDataOption) {
			raw := o.StringValue()
			// Ensure it matches the regular expression.
			if !hexMatchRegex.MatchString(raw) {
				s.InteractionRespond(i.Interaction, respondText("Invalid hex code provided (format unknown).", true))
				return
			}
			// Extract all the non-hex characters.
			parsed := hexReplaceRegex.ReplaceAllString(raw, "")
			values, err := strconv.ParseUint(string(parsed), 16, 32)
			if err != nil {
				s.InteractionRespond(i.Interaction, respondText("Invalid hex code provided (out of range).", true))
				return
			}
			// Red is MS 8 bits.
			red := uint8(values >> 16)
			// Green is middle 8 bits
			green := uint8((values >> 8) & 0xFF)
			// Blue is LS 8 bits.
			blue := uint8(values & 0xFF)
			// Format into nice output.
			result := fmt.Sprintf("R: %d, G: %d, B: %d", red, green, blue)
			embed := embed()
			embed.field("Result", result, false)
			s.InteractionRespond(i.Interaction, respondEmbed(embed, true))
		})
	})
}

// The ping command.
// Will just respond with a message.
var handlePing = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, respondText("Pong! Command handler online and responsive.", true))
}

// Helper variables.
var hexMatchRegex = regexp.MustCompile("#?[0-9a-fA-F]{6}")
var hexReplaceRegex = regexp.MustCompile("[^0-9a-fA-F]")
