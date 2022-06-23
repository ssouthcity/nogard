package nogard

import (
	"math"
	"time"
)

type DragonEncyclopedia interface {
	SearchDragons(query string) ([]string, error)
	Dragon(name string) (*Dragon, error)
}

type Dragon struct {
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	StartingCash float64       `json:"starting_cash"`
	Availability Availability  `json:"availability"`
	Incubation   time.Duration `json:"incubation"`
	Rarity       Rarity        `json:"rarity"`
	Habitats     []string      `json:"habitats"`
	Elements     []string      `json:"elements"`
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
