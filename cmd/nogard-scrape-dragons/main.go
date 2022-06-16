package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/antchfx/htmlquery"
	"github.com/antchfx/xpath"
	"github.com/ssouthcity/nogard"
)

type DragonStore struct {
	Dragons []nogard.Dragon
	Mutex   sync.RWMutex
}

type FandomRenderResponse struct {
	Parsed struct {
		Title string            `json:"title"`
		Text  map[string]string `json:"text"`
	} `json:"parse"`
}

type FandomListResponse struct {
	Next *struct {
		CategoryMembers string `json:"cmcontinue"`
		Continue        string `json:"continue"`
	} `json:"continue"`
	Query struct {
		CategoryMembers []struct {
			PageID    int    `json:"pageid"`
			Namespace int    `json:"ns"`
			Title     string `json:"title"`
		} `json:"categorymembers"`
	} `json:"query"`
}

var (
	XPathDragonBox         = xpath.MustCompile("//table[contains(@class, 'dragonbox')]")
	XPathDragonDescription = xpath.MustCompile("tbody/tr[th[contains(text(), 'Game Description')]]/following-sibling::tr[1]/td//text()")
	XPathBreedingTime      = xpath.MustCompile("tbody/tr/th[contains(span/text(), 'Breeding Times')]/following-sibling::td[1]/text()")
	XPathAvailability      = xpath.MustCompile("tbody/tr/th[contains(span/text(), 'Limited')]/following-sibling::td[1]/span/b/text()")
	XPathRarity            = xpath.MustCompile("tbody/tr/th[contains(span/text(), 'Rarity Status')]/following-sibling::td[1]/a/span/b/text()")
	XPathEarningRates      = xpath.MustCompile("//h2[span[@id = 'Earning_Rates']]/following-sibling::table[1]//td/text()")
)

func scrapeDragon(pageID int) (d nogard.Dragon) {
	body := &FandomRenderResponse{}

	u, err := url.Parse("https://dragonvale.fandom.com/api.php?action=parse&format=json")
	if err != nil {
		log.Fatalf("malformed url '%s'", err)
	}

	q := u.Query()
	q.Set("pageid", fmt.Sprint(pageID))
	u.RawQuery = q.Encode()

	res, err := http.Get(u.String())
	if err != nil {
		log.Fatalf("http request failed '%s'", err)
	}

	defer res.Body.Close()
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		log.Fatalf("malformed body '%s'", err)
	}

	d.Name = body.Parsed.Title

	doc, err := htmlquery.Parse(strings.NewReader(body.Parsed.Text["*"]))
	if err != nil {
		log.Panic(err)
	}

	table := htmlquery.QuerySelector(doc, XPathDragonBox)

	descriptionNodes := htmlquery.QuerySelectorAll(table, XPathDragonDescription)
	if len(descriptionNodes) > 0 {
		d.Description = strings.TrimSpace(descriptionNodes[0].Data)
		for _, node := range descriptionNodes[1:] {
			d.Description += "\n" + strings.TrimSpace(node.Data)
		}
	}
	d.Description = strings.TrimSuffix(d.Description, "\n")

	earningNodes := htmlquery.QuerySelectorAll(doc, XPathEarningRates)
	if len(earningNodes) > 0 {
		for i, node := range earningNodes {
			if i > 20 {
				break
			}

			data := strings.TrimSuffix(node.Data, "\n")

			val, err := strconv.Atoi(data)
			if err != nil {
				continue
			}

			d.Earnings = append(d.Earnings, val)
		}
	}

	breedingNode := htmlquery.QuerySelector(table, XPathBreedingTime)
	if breedingNode != nil {
		replacer := strings.NewReplacer(" ", "", "\n", "", "sec", "s")

		durationString := replacer.Replace(strings.ToLower(breedingNode.Data))

		if durationString == "instant" {
			d.Incubation = time.Duration(0)
		} else {
			dur, err := time.ParseDuration(durationString)
			if err != nil {
				log.Print(d.Name)
				panic(err)
			}
			d.Incubation = dur
		}
	}

	availabilityNode := htmlquery.QuerySelector(table, XPathAvailability)
	if availabilityNode != nil {
		switch availabilityNode.Data {
		case "PERMANENT":
			d.Availability = nogard.Permanent
		case "AVAILABLE":
			d.Availability = nogard.Available
		case "EXPIRED":
			d.Availability = nogard.Unavailable
		}
	}

	rarityNode := htmlquery.QuerySelector(table, XPathRarity)
	if rarityNode != nil {
		table := map[string]nogard.Rarity{
			"Primary":   nogard.Primary,
			"Hybrid":    nogard.Hybrid,
			"Rare":      nogard.Rare,
			"Epic":      nogard.Epic,
			"Gemstone":  nogard.Gemstone,
			"Galaxy":    nogard.Galaxy,
			"Mythic":    nogard.Mythic,
			"Legendary": nogard.Legendary,
		}

		d.Rarity = table[rarityNode.Data]
	}

	return
}

func main() {
	ds := &DragonStore{
		Dragons: make([]nogard.Dragon, 0),
		Mutex:   sync.RWMutex{},
	}
	wg := sync.WaitGroup{}

	u, err := url.Parse("https://dragonvale.fandom.com/api.php?action=query&format=json&list=categorymembers&cmtitle=Category%3ADragons&cmlimit=500")
	if err != nil {
		log.Fatalf("malformed url '%s'", err)
	}

	log.Print("scrape begin")

	for {
		body := FandomListResponse{}

		res, err := http.Get(u.String())
		if err != nil {
			log.Fatalf("http request failed '%s'", err)
		}

		if res.StatusCode != http.StatusOK {
			time.Sleep(5 * time.Second)
			continue
		}

		defer res.Body.Close()
		if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
			log.Fatalf("invalid response body '%s'", err)
		}

		for _, page := range body.Query.CategoryMembers {
			if page.Namespace != 0 {
				continue
			}

			if !strings.HasSuffix(page.Title, "Dragon") {
				continue
			}

			wg.Add(1)

			go func(id int) {
				defer wg.Done()
				d := scrapeDragon(id)

				defer ds.Mutex.Unlock()
				ds.Mutex.Lock()
				ds.Dragons = append(ds.Dragons, d)
			}(page.PageID)
		}

		if body.Next != nil {
			q := u.Query()
			q.Set("continue", body.Next.Continue)
			q.Set("cmcontinue", body.Next.CategoryMembers)
			u.RawQuery = q.Encode()
		} else {
			break
		}
	}

	wg.Wait()

	log.Print("scrape end")

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "	")
	if err := enc.Encode(ds.Dragons); err != nil {
		panic(err)
	}
}
