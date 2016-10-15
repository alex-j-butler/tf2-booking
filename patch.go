package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type PatchUser struct {
	*discordgo.User
}

func (u *PatchUser) GetMention() string {
	return fmt.Sprintf("<@%s>", u.ID)
}

// Helper function
func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
