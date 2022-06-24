package main

import (
	"flag"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"github.com/ssouthcity/nogard"
	"github.com/ssouthcity/nogard/compendium"
	"github.com/ssouthcity/nogard/fandom"
	"github.com/ssouthcity/nogard/interaction"
)

func main() {
	logger := logrus.New()

	token := flag.String("token", os.Getenv("NOGARD_TOKEN"), "")
	sheetID := flag.String("sheetID", os.Getenv("NOGARD_SHEET_ID"), "")
	sheetCreds := flag.String("credentialsPath", os.Getenv("NOGARD_SHEET_CREDENTIALS"), "")
	flag.Parse()

	s, err := discordgo.New("Bot " + *token)
	if err != nil {
		logger.WithError(err).Fatal("invalid session configuration")
	}

	var dragonSrv nogard.DragonEncyclopedia
	{
		dragonSrv, err = compendium.NewDragonEncyclopedia(*sheetID, *sheetCreds)
		if err != nil {
			logger.WithError(err).WithField("sheetID", *sheetID).Fatal("couldn't initialize encyclopedia service")
		}

		dragonSrv = fandom.NewDragonDescriptionPatcher(dragonSrv)
	}

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
