package main

import (
	"flag"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"github.com/ssouthcity/nogard/interaction"
	"github.com/ssouthcity/nogard/local"
)

func main() {
	logger := logrus.New()

	token := flag.String("token", os.Getenv("NOGARD_TOKEN"), "")
	flag.Parse()

	s, err := discordgo.New("Bot " + *token)
	if err != nil {
		logger.WithError(err).Fatal("invalid session configuration")
	}

	dragonSrv := local.NewDragonService()

	r := &interaction.InteractionRouter{
		Encyclopedia: dragonSrv,
		Logger:       logger,
	}

	s.AddHandler(r.HandleInteraction)

	if err := s.Open(); err != nil {
		logger.WithError(err).Fatal("discord gateway connection severed")
	}
	defer s.Close()

	select {}
}
