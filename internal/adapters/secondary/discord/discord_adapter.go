package discord

import (
	"log"

	"github.com/actions/internal/core/domain"
	"github.com/bwmarrin/discordgo"
)

type DiscordAdapter struct {
	token   string
	guildID string
	session *discordgo.Session

	developmentChannelID string
	testsChannelID       string
}

func NewDiscordAdapter(token, guildID string) (*DiscordAdapter, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	return &DiscordAdapter{
		token:   token,
		guildID: guildID,
		session: session,
	}, nil
}

func (a *DiscordAdapter) Initialize() error {
	// Crear los canales necesarios
	channels := []struct {
		name        string
		channelType discordgo.ChannelType
	}{
		{"desarrollo", discordgo.ChannelTypeGuildText},
		{"pruebas", discordgo.ChannelTypeGuildText},
		{"general", discordgo.ChannelTypeGuildText},
	}

	for _, ch := range channels {
		channel, err := a.session.GuildChannelCreate(a.guildID, ch.name, ch.channelType)
		if err != nil {
			log.Printf("Error al crear el canal %s: %v", ch.name, err)
			continue
		}

		// Guardar los IDs de los canales importantes
		switch ch.name {
		case "desarrollo":
			a.developmentChannelID = channel.ID
		case "pruebas":
			a.testsChannelID = channel.ID
		}
	}

	return nil
}

func (a *DiscordAdapter) Send(notification domain.Notification) error {
	_, err := a.session.ChannelMessageSend(notification.ChannelID, notification.Message)
	return err
}

func (a *DiscordAdapter) CreateChannels() error {
	return a.Initialize()
}

func (a *DiscordAdapter) GetDevelopmentChannelID() string {
	return a.developmentChannelID
}

func (a *DiscordAdapter) GetTestsChannelID() string {
	return a.testsChannelID
}
