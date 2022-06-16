package nogard

import "time"

type Availability uint8

const (
	Permanent Availability = iota
	Available
	Unavailable
)

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
