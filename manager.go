package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/NaeCqde/dca"
	"github.com/bwmarrin/discordgo"
)

var Instances = []*VoiceManager{}

type VoiceManager struct {
	Conn    *discordgo.VoiceConnection
	Channel *discordgo.Channel
	GuildID string
	UserID  string
	Stream  *dca.EncodeSession
}

func (m *VoiceManager) Play(audio Audio) {
	slog.Debug("Play", "Audio", audio)
	m.Skip()

	stream, err := GetAudio(audio)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	m.Stream = stream
	done := make(chan error)

	dca.NewStream(m.Stream, m.Conn, done)

	err = <-done
	if err != nil {
		slog.Error(err.Error())
	}
}

func (m *VoiceManager) Skip() {
	if m.Stream != nil {
		m.Stream.Cleanup()
		m.Stream = nil
	}
}

func (m *VoiceManager) Leave() {
	m.Conn.Disconnect(context.TODO())

	new := []*VoiceManager{}

	for _, vm := range Instances {
		if m.GuildID != vm.GuildID {
			new = append(new, vm)
		}
	}

	Instances = new
}

func NewManager(
	conn *discordgo.VoiceConnection,
	guildID string,
	userID string,
	channel *discordgo.Channel,
) (*VoiceManager, error) {
	instance := &VoiceManager{
		Conn:    conn,
		Channel: channel,
		GuildID: guildID,
		UserID:  userID,
	}

	Instances = append(Instances, instance)

	return instance, instance.Conn.Speaking(true)
}

func GetManager(guildID string, userID string) (*VoiceManager, error) {
	for _, instance := range Instances {
		if guildID == instance.GuildID && userID == instance.UserID {
			return instance, nil
		}
	}

	return &VoiceManager{}, fmt.Errorf("failed to find Voice Manager")
}
