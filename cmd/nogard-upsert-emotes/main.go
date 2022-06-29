package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

func main() {
	token := os.Getenv("NOGARD_TOKEN")
	guild := os.Getenv("NOGARD_GUILD")

	s, err := discordgo.New("Bot " + token)
	if err != nil {
		log.WithError(err).Fatal("invalid session configuration")
	}

	defer s.Close()
	if err := s.Open(); err != nil {
		log.WithError(err).Fatal("discord gateway connection dropped")
	}

	dir, err := os.ReadDir("assets/icons")
	if err != nil {
		log.WithError(err).Fatal("unable to open directory")
	}

	emojis, err := s.GuildEmojis(guild)
	if err != nil {
		log.WithError(err).Fatal("could not fetch emotes")
	}

	log.Info("starting")

Main:
	for _, entry := range dir {
		name := entry.Name()
		if strings.HasPrefix(name, "el-") {
			emoteName := strings.TrimPrefix(name, "el-")
			emoteName = strings.TrimSuffix(emoteName, ".png")
			emoteName = fmt.Sprintf("%s_element", emoteName)

			for _, emoji := range emojis {
				if emoji.Name == emoteName {
					continue Main
				}
			}

			path := fmt.Sprintf("assets/icons/%s", name)

			f, err := os.Open(path)
			if err != nil {
				log.WithError(err).Fatal("unable to open image")
			}
			defer f.Close()

			buf, err := io.ReadAll(f)
			if err != nil {
				log.WithError(err).Fatal("unable to read file")
			}

			encStr := fmt.Sprintf("data:image/png;base64,%s", base64.StdEncoding.EncodeToString(buf))

			if _, err := s.GuildEmojiCreate(guild, emoteName, encStr, nil); err != nil {
				log.WithError(err).WithField("name", emoteName).Fatal("unable to create emote")
			}

			log.WithField("name", emoteName).Info("creating emoji")
		}
	}

	log.Info("success")
}
