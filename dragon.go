package nogard

import "time"

type DragonEncyclopedia interface {
	SearchDragons(query string) ([]string, error)
	Dragon(name string) (*Dragon, error)
}

type Dragon struct {
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	Earnings     []int         `json:"earnings"`
	Availability Availability  `json:"availability"`
	Incubation   time.Duration `json:"incubation"`
	Rarity       Rarity        `json:"rarity"`
	Elements     []string      `json:"elements"`
	Latent       []string      `json:"latent"`
	Breeding     [][]string    `json:"breeding"`
}

func (d *Dragon) BreedingTime(upgraded bool) time.Duration {
	modifier := 1.0
	if upgraded {
		modifier = 0.8
	}

	return time.Duration(float64(d.Incubation.Nanoseconds()) * modifier)
}
