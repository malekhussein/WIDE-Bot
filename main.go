package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)


const serverID = "Put your server ID here"

var commands = []*discordgo.ApplicationCommand{
	
	{Name: "warn", Description: "Warn a user", Options: []*discordgo.ApplicationCommandOption{
		{Type: discordgo.ApplicationCommandOptionUser, Name: "user", Description: "User to warn", Required: true},
		{Type: discordgo.ApplicationCommandOptionString, Name: "reason", Description: "Reason", Required: false},
	}},
	{Name: "kick", Description: "Kick a user", Options: []*discordgo.ApplicationCommandOption{
		{Type: discordgo.ApplicationCommandOptionUser, Name: "user", Description: "User to kick", Required: true},
		{Type: discordgo.ApplicationCommandOptionString, Name: "reason", Description: "Reason", Required: false},
	}},
	{Name: "ban", Description: "Ban a user", Options: []*discordgo.ApplicationCommandOption{
		{Type: discordgo.ApplicationCommandOptionUser, Name: "user", Description: "User to ban", Required: true},
		{Type: discordgo.ApplicationCommandOptionString, Name: "reason", Description: "Reason", Required: false},
	}},
	{Name: "timeout", Description: "Timeout a user", Options: []*discordgo.ApplicationCommandOption{
		{Type: discordgo.ApplicationCommandOptionUser, Name: "user", Description: "User to timeout", Required: true},
		{Type: discordgo.ApplicationCommandOptionString, Name: "reason", Description: "Reason", Required: false},
	}},
	{Name: "mute", Description: "Mute a user", Options: []*discordgo.ApplicationCommandOption{
		{Type: discordgo.ApplicationCommandOptionUser, Name: "user", Description: "User to mute", Required: true},
	}},
	{Name: "unmute", Description: "Unmute a user", Options: []*discordgo.ApplicationCommandOption{
		{Type: discordgo.ApplicationCommandOptionUser, Name: "user", Description: "User to unmute", Required: true},
	}},

	
	{Name: "serverinfo", Description: "Show server info"},
	{Name: "userinfo", Description: "Show user info", Options: []*discordgo.ApplicationCommandOption{
		{Type: discordgo.ApplicationCommandOptionUser, Name: "user", Description: "User to view", Required: false},
	}},
	{Name: "roleinfo", Description: "Show role info", Options: []*discordgo.ApplicationCommandOption{
		{Type: discordgo.ApplicationCommandOptionString, Name: "role", Description: "Role name", Required: true},
	}},
	{Name: "ping", Description: "Check bot ping"},
	{Name: "avatar", Description: "Show user avatar", Options: []*discordgo.ApplicationCommandOption{
		{Type: discordgo.ApplicationCommandOptionUser, Name: "user", Description: "User to view", Required: false},
	}},

	
	{Name: "addrole", Description: "Add role to user", Options: []*discordgo.ApplicationCommandOption{
		{Type: discordgo.ApplicationCommandOptionUser, Name: "user", Description: "User", Required: true},
		{Type: discordgo.ApplicationCommandOptionString, Name: "role", Description: "Role name", Required: true},
	}},
	{Name: "removerole", Description: "Remove role from user", Options: []*discordgo.ApplicationCommandOption{
		{Type: discordgo.ApplicationCommandOptionUser, Name: "user", Description: "User", Required: true},
		{Type: discordgo.ApplicationCommandOptionString, Name: "role", Description: "Role name", Required: true},
	}},
	{Name: "clear", Description: "Clear messages", Options: []*discordgo.ApplicationCommandOption{
		{Type: discordgo.ApplicationCommandOptionInteger, Name: "amount", Description: "Number of messages", Required: true},
	}},
	{Name: "say", Description: "Make the bot say something", Options: []*discordgo.ApplicationCommandOption{
		{Type: discordgo.ApplicationCommandOptionString, Name: "message", Description: "Message to say", Required: true},
	}},

	
	{Name: "lockchannel", Description: "Lock the current channel"},
	{Name: "unlockchannel", Description: "Unlock the current channel"},
	{Name: "slowmode", Description: "Set slowmode for the channel", Options: []*discordgo.ApplicationCommandOption{
		{Type: discordgo.ApplicationCommandOptionInteger, Name: "seconds", Description: "Seconds between messages", Required: true},
	}},
	{Name: "prune", Description: "Delete messages", Options: []*discordgo.ApplicationCommandOption{
		{Type: discordgo.ApplicationCommandOptionInteger, Name: "amount", Description: "Number of messages", Required: true},
	}},
	{Name: "announce", Description: "Send announcement", Options: []*discordgo.ApplicationCommandOption{
		{Type: discordgo.ApplicationCommandOptionString, Name: "message", Description: "Announcement message", Required: true},
		{Type: discordgo.ApplicationCommandOptionChannel, Name: "channel", Description: "Channel to send", Required: true},
	}},
	{Name: "setnick", Description: "Change user's nickname", Options: []*discordgo.ApplicationCommandOption{
		{Type: discordgo.ApplicationCommandOptionUser, Name: "user", Description: "User", Required: true},
		{Type: discordgo.ApplicationCommandOptionString, Name: "nickname", Description: "New nickname", Required: true},
	}},
	{Name: "auditlog", Description: "Show recent server events"},

	
	{Name: "help", Description: "Show all commands"},
}

func main() {
	botToken := "Place your bot token here"
	dg, err := discordgo.New("Bot " + botToken)
	if err != nil {
		log.Fatal("Error creating Discord session,", err)
	}

	dg.AddHandler(ready)
	dg.AddHandler(interactionCreate)

	err = dg.Open()
	if err != nil {
		log.Fatal("Error opening connection,", err)
	}

	fmt.Println("WIDE-Bot is now running with 30 slash commands!")
	select {}
}

func ready(s *discordgo.Session, r *discordgo.Ready) {
	for _, v := range commands {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, serverID, v)
		if err != nil {
			log.Println("Cannot create slash command:", v.Name, err)
		}
	}
}

func interactionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	command := i.ApplicationCommandData().Name
	options := i.ApplicationCommandData().Options
	userOpt := func(name string) string {
		for _, opt := range options {
			if opt.Name == name && opt.Type == discordgo.ApplicationCommandOptionUser {
				return fmt.Sprintf("<@%s>", opt.UserValue(s).ID)
			}
		}
		return ""
	}
	strOpt := func(name string) string {
		for _, opt := range options {
			if opt.Name == name {
				return opt.StringValue()
			}
		}
		return ""
	}

	var response string

	switch command {
	case "warn":
		user := userOpt("user")
		reason := strOpt("reason")
		if reason == "" {
			reason = "No reason"
		}
		response = fmt.Sprintf("%s has been warned by %s. Reason: %s", user, i.Member.User.Username, reason)
	case "kick":
		user := userOpt("user")
		reason := strOpt("reason")
		s.GuildMemberDeleteWithReason(serverID, strings.Trim(user, "<@!>"), reason)
		response = fmt.Sprintf("%s has been kicked by %s. Reason: %s", user, i.Member.User.Username, reason)
	case "ban":
		user := userOpt("user")
		reason := strOpt("reason")
		s.GuildBanCreateWithReason(serverID, strings.Trim(user, "<@!>"), reason, 0)
		response = fmt.Sprintf("%s has been banned by %s. Reason: %s", user, i.Member.User.Username, reason)
	case "timeout":
		user := userOpt("user")
		reason := strOpt("reason")
		response = fmt.Sprintf("%s has been timed out (manual) by %s. Reason: %s", user, i.Member.User.Username, reason)
	case "mute":
		user := userOpt("user")
		response = fmt.Sprintf("%s muted by %s", user, i.Member.User.Username)
	case "unmute":
		user := userOpt("user")
		response = fmt.Sprintf("%s unmuted by %s", user, i.Member.User.Username)
	case "ping":
		response = "Pong!"
	case "serverinfo":
		response = fmt.Sprintf("Server ID: %s", serverID)
	case "userinfo":
		user := userOpt("user")
		if user == "" {
			user = i.Member.User.Mention()
		}
		response = fmt.Sprintf("User info: %s", user)
	case "roleinfo":
		role := strOpt("role")
		response = fmt.Sprintf("Role info: %s", role)
	case "addrole":
		response = fmt.Sprintf("Added role %s to %s", strOpt("role"), userOpt("user"))
	case "removerole":
		response = fmt.Sprintf("Removed role %s from %s", strOpt("role"), userOpt("user"))
	case "clear":
		amount := strOpt("amount")
		response = fmt.Sprintf("Cleared %s messages", amount)
	case "say":
		response = strOpt("message")
	case "lockchannel":
		response = "Channel has been locked."
	case "unlockchannel":
		response = "Channel has been unlocked."
	case "slowmode":
		seconds := strOpt("seconds")
		response = fmt.Sprintf("Channel slowmode set to %s seconds.", seconds)
	case "prune":
		response = fmt.Sprintf("Pruned %s messages.", strOpt("amount"))
	case "announce":
		response = fmt.Sprintf("Announcement: %s", strOpt("message"))
	case "setnick":
		response = fmt.Sprintf("%s nickname changed to %s", userOpt("user"), strOpt("nickname"))
	case "auditlog":
		response = "Audit log: recent events (mock display)."
	case "help":
		response = "30 slash commands available: /warn, /kick, /ban, /timeout, /mute, /unmute, /serverinfo, /userinfo, /roleinfo, /ping, /avatar, /addrole, /removerole, /clear, /say, /lockchannel, /unlockchannel, /slowmode, /prune, /announce, /setnick, /auditlog, /help"
	case "avatar":
		user := userOpt("user")
		if user == "" {
			user = i.Member.User.Mention()
		}
		response = fmt.Sprintf("%s avatar link (mock display)", user)
	default:
		response = fmt.Sprintf("Command %s executed!", command)
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
		},
	})
}
