package fandom

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/antchfx/xpath"
	"github.com/ssouthcity/nogard"
)

type DescriptionPatcher struct {
	rootURL         string
	descriptionPath *xpath.Expr
	encyclopedia    nogard.DragonEncyclopedia
}

type parseBody struct {
	Parse struct {
		Text map[string]string `json:"text"`
	} `json:"parse"`
}

func NewDragonDescriptionPatcher(encyclopedia nogard.DragonEncyclopedia) *DescriptionPatcher {
	return &DescriptionPatcher{
		rootURL:         "https://dragonvale.fandom.com/api.php?action=parse&format=json",
		descriptionPath: xpath.MustCompile("//table[contains(@class, 'dragonbox')]/tbody/tr[th[contains(text(), 'Game Description')]]/following-sibling::tr[1]/td//text()"),
		encyclopedia:    encyclopedia,
	}
}

func (p *DescriptionPatcher) DragonNames() ([]string, error) {
	return p.encyclopedia.DragonNames()
}

func (p *DescriptionPatcher) Dragon(name string) (*nogard.Dragon, error) {
	d, err := p.encyclopedia.Dragon(name)
	if err != nil {
		return nil, err
	}

	var page string
	if d.Rarity >= nogard.Mythic {
		page = d.Name
	} else {
		page = fmt.Sprintf("%s_Dragon", d.Name)
	}
	page = strings.ReplaceAll(page, " ", "_")

	u, _ := url.Parse(p.rootURL)

	q := u.Query()
	q.Set("page", page)
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}

	var body parseBody

	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}

	doc, err := htmlquery.Parse(strings.NewReader(body.Parse.Text["*"]))
	if err != nil {
		return nil, err
	}

	descNode := htmlquery.QuerySelector(doc, p.descriptionPath)
	if descNode != nil {
		d.Description = descNode.Data
	}

	return d, nil
}
