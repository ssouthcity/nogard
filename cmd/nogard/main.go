package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/ssouthcity/nogard"
	"github.com/ssouthcity/nogard/local"
)

func main() {
	token := flag.String("token", os.Getenv("NOGARD_TOKEN"), "")
	flag.Parse()

	s, err := discordgo.New("Bot " + *token)
	if err != nil {
		log.Fatalf("invalid session configuration '%s'", err)
	}

	dragonSrv := local.NewDragonService()

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {

		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			handler := NewDragonInfoHandler(dragonSrv)
			handler(s, i)
		case discordgo.InteractionApplicationCommandAutocomplete:
			handler := NewSearchHandler(dragonSrv)
			handler(s, i)
		}
	})

	if err := s.Open(); err != nil {
		log.Fatalf("connection severed '%s'", err)
	}
	defer s.Close()

	select {}
}

type InteractionHandlerFunc func(s *discordgo.Session, i *discordgo.InteractionCreate)

type DragonDetailser interface {
	DragonDetails(name string) (*nogard.Dragon, error)
}

func NewDragonInfoHandler(details DragonDetailser) func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		data := i.ApplicationCommandData()
		name := data.Options[0].StringValue()

		dragon, err := details.DragonDetails(name)
		if err != nil {
			log.Fatalf("unable to find info '%s'", err)
		}

		colors := map[nogard.Rarity]int{
			nogard.Primary:   15663085,
			nogard.Hybrid:    9043831,
			nogard.Rare:      39423,
			nogard.Epic:      9048814,
			nogard.Gemstone:  16742161,
			nogard.Galaxy:    122367,
			nogard.Mythic:    16720827,
			nogard.Legendary: 15606273,
		}

		embed := &discordgo.MessageEmbed{
			URL:         fmt.Sprintf("https://dragonvale.fandom.com/%s", strings.ReplaceAll(dragon.Name, " ", "_")),
			Title:       dragon.Name,
			Description: dragon.Description,
			Color:       colors[dragon.Rarity],
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Availability",
					Value:  fmt.Sprint(dragon.Availability),
					Inline: true,
				},
				{
					Name:   "Breeding",
					Value:  fmt.Sprintf("Regular: %s\nUpgraded: %s", dragon.BreedingTime(false), dragon.BreedingTime(true)),
					Inline: true,
				},
				{
					Name:  "Earnings",
					Value: fmt.Sprintf("%+v", dragon.Earnings),
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
}

type DragonSearcher interface {
	SearchDragon(query string) ([]string, error)
}

func NewSearchHandler(searcher DragonSearcher) func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		data := i.ApplicationCommandData()
		query := data.Options[0].StringValue()

		results, err := searcher.SearchDragon(query)
		if err != nil {
			log.Fatalf("unable to search '%s'", err)
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
}
