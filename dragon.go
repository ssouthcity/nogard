package nogard

type Dragon struct {
	Name         string     `json:"name"`
	Description  string     `json:"description"`
	Type         string     `json:"type"`
	Elements     []string   `json:"elements"`
	Latent       []string   `json:"latent"`
	Time         int        `json:"time"`
	Availability string     `json:"availability"`
	Breeding     [][]string `json:"breeding"`
}

type DragonService interface {
	SearchDragon(query string) ([]string, error)
	DragonInfo(name string) (*Dragon, error)
}
