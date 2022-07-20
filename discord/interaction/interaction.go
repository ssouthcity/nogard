package interaction

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"github.com/ssouthcity/dgimux"
	nogard "github.com/ssouthcity/nogard/lib"
)

type InteractionRouter struct {
	Encyclopedia nogard.DragonEncyclopedia
	Logger       logrus.FieldLogger
}

func (r *InteractionRouter) HandleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	sr := dgimux.NewRouter()

	dragonariumSrv := &dragonariumHandler{
		encylopedia: r.Encyclopedia,
		logger:      r.Logger,
	}

	sr.ApplicationCommand("dragonarium", dgimux.InteractionHandlerFunc(dragonariumSrv.dragonariumCommand))
	sr.ApplicationCommandAutoComplete("dragonarium", dgimux.InteractionHandlerFunc(dragonariumSrv.dragonAutocomplete))

	sr.HandleInteraction(s, i)
}
