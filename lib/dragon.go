package nogard

import (
	"math"
	"strings"
	"time"
)

type DragonEncyclopedia interface {
	DragonNames() ([]string, error)
	Dragon(name string) (*Dragon, error)
}

type Dragon struct {
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	StartingCash float64       `json:"starting_cash"`
	Availability Availability  `json:"availability"`
	Incubation   time.Duration `json:"incubation"`
	Rarity       Rarity        `json:"rarity"`
	Habitats     []Element     `json:"habitats"`
	Elements     []Element     `json:"elements"`
	Combo        []string      `json:"combo"`
}

func (d *Dragon) BreedingTime(upgraded bool) time.Duration {
	modifier := 1.0
	if upgraded {
		modifier = 0.8
	}

	return time.Duration(float64(d.Incubation.Nanoseconds()) * modifier)
}

func (d *Dragon) CashPerMinute(level int) int {
	modifier := float64(level-1) * (0.6 * d.StartingCash)
	return (int)(math.Round(d.StartingCash + modifier))
}

func FilterDragonNames(names []string, query string) []string {
	filteredNames := make([]string, 0)

	for _, name := range names {
		if strings.Contains(strings.ToLower(name), strings.ToLower(query)) {
			filteredNames = append(filteredNames, name)
		}
	}

	return filteredNames
}
