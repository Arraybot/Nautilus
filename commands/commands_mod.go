package commands

import (
	"fmt"
	"log"
	"time"

	"github.com/arraybot/nautilus/database"
	"github.com/arraybot/nautilus/requests"
	"github.com/bwmarrin/discordgo"
)

const (
	punishmentTypeKick    = "KICK"
	punishmentTypeTimeout = "TIMEOUT"
	punishmentTypeMute    = "MUTE"
	punishmentTypeBan     = "BAN"
)

// The clear command.
// Will bulk delete messages in a channel and optionally by a specific user.
func handlerClear(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !hasPermission(i, discordgo.PermissionManageMessages) {
		s.InteractionRespond(i.Interaction, respondText(permissionDenyPermission, i))
		return
	}
	// Load all the command options.
	count := option(i.ApplicationCommandData().Options, "amount")
	c := count.IntValue()
	user := option(i.ApplicationCommandData().Options, "user")
	var u *discordgo.User = nil
	if user != nil {
		u = user.UserValue(nil)
	}
	// Validate that the user range is valid.
	if c < 2 || c > 100 {
		s.InteractionRespond(i.Interaction, respondText("You can only clear at least 2 and at most 100 messages.", i))
		return
	}
	// Try to get the IDs of messages to delete.
	var ids []string
	// Try to load all messages.
	messages, err := s.ChannelMessages(i.ChannelID, int(c), "", "", "")
	if err != nil {
		log.Println(err)
		s.InteractionRespond(i.Interaction, respondText("Error getting messages.", i))
		return
	}
	// Filter to only eligible messages.
	nt := time.Now()
	for _, message := range messages {
		ct := creationTime(message.ID)
		diff := nt.Sub(ct)
		hrs := diff.Hours()
		wks := hrs / (24 * 7)
		// Only include messages less than 2 weeks old.
		if wks <= 2.0 {
			// If user filtering is enabled, do that.
			if u != nil && message.Author.ID != u.ID {
				continue
			}
			ids = append(ids, message.ID)
		}
	}
	// Delete all the messages.
	err = s.ChannelMessagesBulkDelete(i.ChannelID, ids)
	if err != nil {
		log.Println(err)
		s.InteractionRespond(i.Interaction, respondText("An error ocurred deleting messages", i))
	} else {
		s.InteractionRespond(i.Interaction, respondText(fmt.Sprintf("Deleted %d messages.", len(ids)), i))
	}
}

// The mute command.
// This will add the muted role to the user.
func handlerMute(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !hasMutePermission(i) {
		s.InteractionRespond(i.Interaction, respondText(permissionDenyMuteRole, i))
		return
	}
	user := option(i.ApplicationCommandData().Options, "user").UserValue(nil).ID
	muteRole := database.MuteRole(i.GuildID)
	err := s.GuildMemberRoleAdd(i.GuildID, user, muteRole)
	if err != nil {
		log.Println(err)
		s.InteractionRespond(i.Interaction, respondText("An error occurred adding the role.", i))
	} else {
		s.InteractionRespond(i.Interaction, respondText("The person has been muted.", i))
	}
}

// The unmute command.
// This will remove the muted role from the user.
func handlerUnMute(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !hasMutePermission(i) {
		s.InteractionRespond(i.Interaction, respondText(permissionDenyMuteRole, i))
		return
	}
	user := option(i.ApplicationCommandData().Options, "user").UserValue(nil).ID
	muteRole := database.MuteRole(i.GuildID)
	err := s.GuildMemberRoleRemove(i.GuildID, user, muteRole)
	if err != nil {
		log.Println(err)
		s.InteractionRespond(i.Interaction, respondText("An error occurred removing the role.", i))
	} else {
		s.InteractionRespond(i.Interaction, respondText("The person has been unmuted.", i))
	}
}

// The revoke command.
// This will revoke the punishment either now or later.
func handlerRevoke(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Parse the relative time and add it to the current time.
	clock := time.Now()
	minutesOption := option(i.ApplicationCommandData().Options, "minutes")
	hoursOption := option(i.ApplicationCommandData().Options, "hours")
	daysOption := option(i.ApplicationCommandData().Options, "days")
	weeksOption := option(i.ApplicationCommandData().Options, "weeks")
	monthsOption := option(i.ApplicationCommandData().Options, "months")
	if minutesOption != nil && minutesOption.IntValue() > 0 {
		clock = clock.Add(time.Minute * time.Duration(minutesOption.IntValue()))
	}
	if hoursOption != nil && hoursOption.IntValue() > 0 {
		clock = clock.Add(time.Hour * time.Duration(hoursOption.IntValue()))
	}
	if daysOption != nil && daysOption.IntValue() > 0 {
		clock = clock.AddDate(0, 0, int(daysOption.IntValue()))
	}
	if weeksOption != nil && weeksOption.IntValue() > 0 {
		clock = clock.AddDate(0, 0, 7*int(daysOption.IntValue()))
	}
	if monthsOption != nil && monthsOption.IntValue() > 0 {
		clock = clock.AddDate(0, int(daysOption.IntValue()), 0)
	}
	// Attempt to load command.
	id := option(i.ApplicationCommandData().Options, "case").IntValue()
	punishment := database.GetPunishment(i.GuildID, fmt.Sprintf("%d", id))
	if punishment == nil {
		s.InteractionRespond(i.Interaction, respondText("This case could not be found.", i))
		return
	}
	// Determine permission.
	permission := false
	switch punishment.Type {
	case punishmentTypeKick:
		permission = hasPermission(i, discordgo.PermissionKickMembers)
	case punishmentTypeTimeout:
		permission = hasPermission(i, discordgo.PermissionModerateMembers)
	case punishmentTypeMute:
		permission = hasMutePermission(i)
	case punishmentTypeBan:
		permission = hasPermission(i, discordgo.PermissionBanMembers)
	}
	if !permission {
		s.InteractionRespond(i.Interaction, respondText("You do not have the permission to revoke this punishment.", i))
		return
	}
	// Schedule expiry and complete.
	epoch := clock.Unix()
	if err := requests.ListenerExpire(i.GuildID, id, epoch); err != nil {
		s.InteractionRespond(i.Interaction, respondText("An error occurred trying to revoke the punishment.", i))
	} else {
		s.InteractionRespond(i.Interaction, respondText(fmt.Sprintf("The punishment will now be revoked at <t:%d:f>. As long as the punishment is not yet expired, you may change its revocation time/date.", epoch), i))
	}
}

// The lookup command.
// Looks up a punishment by case.
func handlerLookup(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !hasModerator(i) {
		s.InteractionRespond(i.Interaction, respondText(permissionDenyModerator, i))
		return
	}
	s.InteractionRespond(i.Interaction, respondText("Lookup", i))
}

// The history command.
// Looks up all the punishment history by a specific user.
func handlerHistory(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !hasModerator(i) {
		s.InteractionRespond(i.Interaction, respondText(permissionDenyModerator, i))
		return
	}
	s.InteractionRespond(i.Interaction, respondText("History", i))
}
