package local

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ssouthcity/nogard"
)

//go:embed dragons.json
var dragons []byte

type DragonService struct {
	dragons map[string]*nogard.Dragon
}

func NewDragonService() *DragonService {
	ds := &DragonService{}
	json.Unmarshal(dragons, &ds.dragons)
	return ds
}

func (s *DragonService) SearchDragon(query string) ([]string, error) {
	suggestions := []string{}

	for name := range s.dragons {
		if strings.HasPrefix(strings.ToLower(name), strings.ToLower(query)) {
			suggestions = append(suggestions, name)
		}
	}

	return suggestions, nil
}

func (s *DragonService) DragonDetails(name string) (*nogard.Dragon, error) {
	if dragon, ok := s.dragons[name]; ok {
		return dragon, nil
	}
	return nil, fmt.Errorf("dragon with name '%s' not found", name)
}
