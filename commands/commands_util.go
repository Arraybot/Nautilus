package commands

import (
	"fmt"
	"log"
	"regexp"
	"runtime"
	"strconv"
	"time"

	"github.com/arraybot/nautilus/config"
	"github.com/arraybot/nautilus/requests"
	"github.com/bwmarrin/discordgo"
)

// The convert command.
// Will respond with a fancy embed.
func handleConvert(s *discordgo.Session, i *discordgo.InteractionCreate) {
	o := i.ApplicationCommandData().Options
	sub(o, "hex", func(o []*commandOption) {
		// These values will strategically underflow and cause values > 255.
		r := uint32(option(o, "red").IntValue())
		g := uint32(option(o, "green").IntValue())
		b := uint32(option(o, "blue").IntValue())
		// Check if the input range is valid.
		if r > 255 || g > 255 || b > 255 {
			s.InteractionRespond(i.Interaction, respondText("Invalid values for at least one colour (out of range, [0-255]).", i))
			return
		}
		result := fmt.Sprintf("#%x", rgbToHex(r, g, b))
		embed := embed()
		embed.field("Result", result, false)
		s.InteractionRespond(i.Interaction, respondEmbed(embed, i))
	})
	sub(o, "rgb", func(o []*commandOption) {
		raw := option(o, "colour").StringValue()
		// Ensure it matches the regular expression.
		if !hexMatchRegex.MatchString(raw) {
			s.InteractionRespond(i.Interaction, respondText("Invalid hex code provided (format unknown).", i))
			return
		}
		// Extract all the non-hex characters.
		parsed := hexReplaceRegex.ReplaceAllString(raw, "")
		hex, err := strconv.ParseUint(parsed, 16, 32)
		if err != nil {
			s.InteractionRespond(i.Interaction, respondText("Invalid hex code provided (out of range).", i))
			return
		}
		r, g, b := hexToRgb(hex)
		// Format into nice output.
		result := fmt.Sprintf("R: %d, G: %d, B: %d", r, g, b)
		embed := embed()
		embed.field("Result", result, false)
		s.InteractionRespond(i.Interaction, respondEmbed(embed, i))
	})
}

// The help command.
// Will display an embed with some useful information.
func handleHelp(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
	s.InteractionRespond(i.Interaction, respondEmbed(embed, i))
}

// The invite command.
// Will send an embed with clickable links to invite the bot and join the server.
func handleInvite(s *discordgo.Session, i *discordgo.InteractionCreate) {
	embed := embed()
	embed.field("Invite Bot", "[https://arraybot.xyz/go/invite/](https://arraybot.xyz/go/invite)", false)
	embed.field("Join Server", "[https://arraybot.xyz/go/server/](https://arraybot.xyz/go/server/)", false)
	s.InteractionRespond(i.Interaction, respondEmbedRaw(embed, 64))
}

// The ping command.
// Will just respond with a message.
func handlePing(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var healthString string
	if hasDeveloper(i) {
		health, err := requests.PanelHealthcheck()
		if err == nil {
			healthString = fmt.Sprintf("Panel handling %d concurrent connections.", health.Connections)
		} else {
			healthString = "Panel unreachable."
			log.Println(err)
		}
	} else {
		healthString = ""
	}
	s.InteractionRespond(i.Interaction, respondText("Pong! Command handler online and responsive. "+healthString, i))
}

// The stats command.
// Will show a few statistics.
func handleStats(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
	elapsed := time.Since(config.StartTime).String()
	embed.field("Tot.Alloc.", fmt.Sprintf("%v MiB", bytesToMegabytes(mem.TotalAlloc)), true)
	embed.field("Sys.Alloc.", fmt.Sprintf("%v MiB", bytesToMegabytes(mem.Sys)), true)
	embed.field("Uptime", elapsed, true)
	s.InteractionRespond(i.Interaction, respondEmbed(embed, i))
}

// Helper variables.
var zwsp = "\u200b"
var hexMatchRegex = regexp.MustCompile("#?[0-9a-fA-F]{6}")
var hexReplaceRegex = regexp.MustCompile("[^0-9a-fA-F]")
