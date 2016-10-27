package util

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
