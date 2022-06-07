package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/ssouthcity/nogard"
	"github.com/ssouthcity/nogard/data"
)

func main() {
	token := flag.String("token", os.Getenv("NOGARD_TOKEN"), "")
	flag.Parse()

	s, err := discordgo.New("Bot " + *token)
	if err != nil {
		log.Fatalf("invalid session configuration '%s'", err)
	}

	dragonSrv := data.NewDragonService()

	s.AddHandler(NewInteractionHandler(dragonSrv))

	if err := s.Open(); err != nil {
		log.Fatalf("connection severed '%s'", err)
	}
	defer s.Close()

	select {}
}

func NewInteractionHandler(ds nogard.DragonService) func(*discordgo.Session, *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type == discordgo.InteractionApplicationCommandAutocomplete {
			data := i.ApplicationCommandData()
			query := data.Options[0].StringValue()

			suggestions, err := ds.SearchDragon(query)
			if err != nil {
				log.Fatalf("searching for dragon failed '%s'", err)
			}

			choices := make([]*discordgo.ApplicationCommandOptionChoice, len(suggestions))

			for i, s := range suggestions {
				choices[i] = &discordgo.ApplicationCommandOptionChoice{
					Name:  s,
					Value: s,
				}
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionApplicationCommandAutocompleteResult,
				Data: &discordgo.InteractionResponseData{
					Choices: choices,
				},
			})
		} else {
			data := i.ApplicationCommandData()
			name := data.Options[0].StringValue()

			dragon, err := ds.DragonInfo(name)
			if err != nil {
				log.Fatalf("querying dragon info failed '%s'", err)
			}

			embed := &discordgo.MessageEmbed{
				Title:       dragon.Name,
				URL:         fmt.Sprintf("https://dragonvale.fandom.com/wiki/%s_Dragon", dragon.Name),
				Description: dragon.Description,
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
}
