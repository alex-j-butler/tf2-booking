package commands

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

type CommandInformation interface {
	ReplyChannel(formatStr string, args ...interface{})
	ReplyUser(formatStr string, args ...interface{})
	GetChannelID() string
	GetUserID() string
}

type DiscordCommandInformation struct {
	*discordgo.MessageCreate

	session *discordgo.Session
}

func (c *DiscordCommandInformation) ReplyChannel(formatStr string, args ...interface{}) {
	c.session.ChannelMessageSend(c.ChannelID, fmt.Sprintf(formatStr, args...))
}

func (c *DiscordCommandInformation) ReplyUser(formatStr string, args ...interface{}) {
	UserChannel, _ := c.session.UserChannelCreate(c.Author.ID)
	c.session.ChannelMessageSend(UserChannel.ID, fmt.Sprintf(formatStr, args...))
}

func (c *DiscordCommandInformation) GetChannelID() string {
	return c.ChannelID
}

func (c *DiscordCommandInformation) GetUserID() string {
	return c.Author.ID
}

type TF2CommandInformation struct {
	steamID    string
	serverName string
}

func (c *TF2CommandInformation) ReplyChannel(formatStr string, args ...interface{}) {
	log.Println("TF2Command, ReplyChannel:", fmt.Sprintf(formatStr, args...))
}

func (c *TF2CommandInformation) ReplyUser(formatStr string, args ...interface{}) {
	log.Println("TF2Command, ReplyUser:", fmt.Sprintf(formatStr, args...))
}

func (c *TF2CommandInformation) GetChannelID() string {
	return c.serverName
}

func (c *TF2CommandInformation) GetUserID() string {
	return c.steamID
}
