package interaction

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"github.com/ssouthcity/nogard"
)

type dragonariumHandler struct {
	encylopedia nogard.DragonEncyclopedia
	logger      logrus.FieldLogger
}

func (h *dragonariumHandler) dragonariumCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	name := data.Options[0].StringValue()

	dragon, err := h.encylopedia.Dragon(name)
	if err != nil {
		h.logger.WithError(err).WithField("name", name).Error("dragon lookup failed")
	}

	embed := &discordgo.MessageEmbed{
		URL:         fmt.Sprintf("https://dragonvale.fandom.com/%s", strings.ReplaceAll(dragon.Name, " ", "_")),
		Title:       dragon.Name,
		Description: dragon.Description,
		Color:       RarityColors[dragon.Rarity],
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Availability",
				Value:  fmt.Sprint(dragon.Availability),
				Inline: true,
			},
			{
				Name:   "Breeding Time",
				Value:  fmt.Sprintf("Regular: %s\nUpgraded: %s", dragon.BreedingTime(false), dragon.BreedingTime(true)),
				Inline: true,
			},
			{
				Name: "Cash Per Minute",
				Value: fmt.Sprintf("```lvl  %4d\t%4d\t%4d\t%4d\t%4d\ngold %4d\t%4d\t%4d\t%4d\t%4d```",
					1, 5, 10, 15, 20,
					dragon.CashPerMinute(1),
					dragon.CashPerMinute(5),
					dragon.CashPerMinute(10),
					dragon.CashPerMinute(15),
					dragon.CashPerMinute(20)),
			},
		},
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "",
			Embeds:  []*discordgo.MessageEmbed{embed},
		},
	})
}

func (h *dragonariumHandler) dragonAutocomplete(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	query := data.Options[0].StringValue()

	results, err := h.encylopedia.SearchDragons(query)
	if err != nil {
		h.logger.WithError(err).WithField("query", query).Error("dragon search failed")
	}

	if len(results) > 25 {
		results = results[:25]
	}

	choices := make([]*discordgo.ApplicationCommandOptionChoice, len(results))

	for i, res := range results {
		choices[i] = &discordgo.ApplicationCommandOptionChoice{
			Name:  strings.TrimSuffix(res, " Dragon"),
			Value: res,
		}
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: choices,
		},
	})
}
