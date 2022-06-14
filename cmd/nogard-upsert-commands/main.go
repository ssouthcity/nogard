package main

import (
	_ "embed"
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

//go:embed commands.json
var commandSpec []byte

func main() {
	token := flag.String("token", os.Getenv("NOGARD_TOKEN"), "")
	guild := flag.String("guild", os.Getenv("NOGARD_GUILD"), "")
	flag.Parse()

	s, err := discordgo.New("Bot " + *token)
	if err != nil {
		log.Fatalf("invalid session configuration '%s'", err)
	}

	if err := s.Open(); err != nil {
		log.Fatalf("connection severed '%s'", err)
	}
	defer s.Close()

	cmds := []*discordgo.ApplicationCommand{}

	if err := json.Unmarshal(commandSpec, &cmds); err != nil {
		log.Fatalf("malformed file content '%s'", err)
	}

	if _, err := s.ApplicationCommandBulkOverwrite(s.State.User.ID, *guild, cmds); err != nil {
		log.Fatalf("updating commands failed '%s'", err)
	}

	log.Print("synchronized commands")
}
