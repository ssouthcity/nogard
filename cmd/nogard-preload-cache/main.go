package main

import (
	"flag"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/ssouthcity/nogard"
	"github.com/ssouthcity/nogard/compendium"
	"github.com/ssouthcity/nogard/fandom"
	"github.com/ssouthcity/nogard/redis"
)

func main() {
	logger := logrus.New()

	sheetID := flag.String("sheetID", os.Getenv("NOGARD_COMPENDIUM_SHEET_ID"), "")
	sheetCreds := flag.String("sheetCredentials", os.Getenv("NOGARD_COMPENDIUM_SHEET_CREDENTIALS"), "")
	redisAddr := flag.String("redis", os.Getenv("NOGARD_REDIS_ADDRESS"), "")
	flag.Parse()

	var dragonSrv nogard.DragonEncyclopedia
	{
		dragonSrv, _ = compendium.NewDragonEncyclopedia(*sheetID, *sheetCreds)
		dragonSrv = fandom.NewDragonDescriptionPatcher(dragonSrv)
		dragonSrv = redis.NewDragonCache(*redisAddr, dragonSrv, logger)
	}

	names, err := dragonSrv.DragonNames()
	if err != nil {
		logger.WithError(err).Fatal("unable to get dragon names")
	}

	for _, name := range names {
		d, err := dragonSrv.Dragon(name)
		if err != nil {
			logger.WithError(err).WithField("name", name).Warn("unable to cache dragon")
		}

		logger.WithField("dragon", d.Name).Info("cached dragon")
	}

	logger.Info("success")
}
