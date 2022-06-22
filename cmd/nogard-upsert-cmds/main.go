package main

import (
	_ "embed"
	"encoding/json"
	"flag"
	"os"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

//go:embed commands.json
var cmdSpec []byte

func main() {
	var (
		token = flag.String("token", os.Getenv("NOGARD_TOKEN"), "bot token for authentication to discord")
		guild = flag.String("guild", os.Getenv("NOGARD_GUILD"), "guild to create commands in, leave empty for global")
	)
	flag.Parse()

	var cmds []*discordgo.ApplicationCommand

	if err := json.Unmarshal(cmdSpec, &cmds); err != nil {
		log.WithError(err).Fatal("failed to unmarshal commands")
	}

	s, err := discordgo.New("Bot " + *token)
	if err != nil {
		log.WithError(err).Fatal("invalid session configuration")
	}

	if err := s.Open(); err != nil {
		log.WithError(err).Fatal("gateway connection not established")
	}
	defer s.Close()

	if _, err := s.ApplicationCommandBulkOverwrite(s.State.User.ID, *guild, cmds); err != nil {
		log.WithError(err).Fatalf("upsert commands operation failed")
	}

	log.WithFields(log.Fields{
		"bot":   s.State.User.Username,
		"guild": *guild,
	}).Info("done, commands have been updated")
}