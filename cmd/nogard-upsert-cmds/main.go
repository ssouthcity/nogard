package main

import (
	_ "embed"
	"encoding/json"
	"os"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

//go:embed commands.json
var cmdSpec []byte

func main() {
	var (
		token = os.Getenv("NOGARD_TOKEN")
		guild = os.Getenv("NOGARD_GUILD")
	)

	var cmds []*discordgo.ApplicationCommand

	if err := json.Unmarshal(cmdSpec, &cmds); err != nil {
		log.WithError(err).Fatal("failed to unmarshal commands")
	}

	s, err := discordgo.New("Bot " + token)
	if err != nil {
		log.WithError(err).Fatal("invalid session configuration")
	}

	if err := s.Open(); err != nil {
		log.WithError(err).Fatal("gateway connection not established")
	}
	defer s.Close()

	if _, err := s.ApplicationCommandBulkOverwrite(s.State.User.ID, guild, cmds); err != nil {
		log.WithError(err).Fatalf("upsert commands operation failed")
	}

	log.WithFields(log.Fields{
		"bot":   s.State.User.Username,
		"guild": guild,
	}).Info("done, commands have been updated")
}
