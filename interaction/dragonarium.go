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

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	dragon, err := h.encylopedia.Dragon(name)
	if err != nil {
		h.logger.WithError(err).WithField("name", name).Error("dragon lookup failed")
		return
	}

	embed := &discordgo.MessageEmbed{
		URL:         fmt.Sprintf("https://dragonvale.fandom.com/%s_Dragon", strings.ReplaceAll(dragon.Name, " ", "_")),
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
		},
	}

	if len(dragon.Combo) != 0 {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:  "Breeding Combination",
			Value: strings.Join(dragon.Combo, " + "),
		})
	}

	if dragon.StartingCash != 0 {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name: "Cash Per Minute",
			Value: fmt.Sprintf("```lvl  %5d\t%5d\t%5d\t%5d\t%5d\ngold %5d\t%5d\t%5d\t%5d\t%5d```",
				1, 5, 10, 15, 20,
				dragon.CashPerMinute(1),
				dragon.CashPerMinute(5),
				dragon.CashPerMinute(10),
				dragon.CashPerMinute(15),
				dragon.CashPerMinute(20)),
		})
	}

	embeds := []*discordgo.MessageEmbed{embed}

	if _, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds: &embeds,
	}); err != nil {
		h.logger.WithError(err).WithFields(logrus.Fields{
			"command": data.Name,
			"dragon":  name,
		}).Error("unable to respond to interaction")
	}
}

func (h *dragonariumHandler) dragonAutocomplete(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	query := data.Options[0].StringValue()

	results, err := h.encylopedia.SearchDragons(query)
	if err != nil {
		h.logger.WithError(err).WithField("query", query).Error("dragon search failed")
		return
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

	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: choices,
		},
	}); err != nil {
		h.logger.WithError(err).WithFields(logrus.Fields{
			"command": data.Name,
			"query":   query,
		}).Error("unable to respond to interaction")
	}
}
