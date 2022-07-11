package compendium

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ssouthcity/nogard"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type DragonEncyclopedia struct {
	sheet         *sheets.Service
	spreadsheetID string
}

func NewDragonEncyclopedia(sheetID string, credentials string) (*DragonEncyclopedia, error) {
	sheetSrv, err := sheets.NewService(nil, option.WithCredentialsJSON([]byte(credentials)))
	if err != nil {
		return nil, err
	}

	return &DragonEncyclopedia{
		sheet:         sheetSrv,
		spreadsheetID: sheetID,
	}, nil
}

func (e *DragonEncyclopedia) DragonNames() ([]string, error) {
	resp, err := e.sheet.Spreadsheets.Values.Get(e.spreadsheetID, "Data!A2:A").Do()
	if err != nil {
		return nil, err
	}

	dragons := make([]string, len(resp.Values))

	for i, row := range resp.Values {
		dragons[i] = row[0].(string)
	}

	return dragons, nil
}

func (e *DragonEncyclopedia) Dragon(name string) (*nogard.Dragon, error) {
	resp, err := e.sheet.Spreadsheets.Values.Get(e.spreadsheetID, "Data!A2:CA").Do()
	if err != nil {
		return nil, err
	}

	for _, row := range resp.Values {
		n := row[0].(string)
		if strings.EqualFold(n, name) {
			return e.mapRowToDragon(row), nil
		}
	}

	return nil, errors.New("dragon not found")
}

func (e *DragonEncyclopedia) rowToRarity(row []interface{}) nogard.Rarity {
	rarity := row[3].(string)

	rarityMap := map[string]nogard.Rarity{"Primary Rift": nogard.Rare}
	for r := nogard.Primary; r <= nogard.Legendary; r++ {
		rarityMap[r.String()] = r
	}

	return rarityMap[rarity]
}

func (e *DragonEncyclopedia) rowToAvailability(row []interface{}) nogard.Availability {
	isLimited, _ := strconv.ParseBool(row[2].(string))
	isLegacy, _ := strconv.ParseBool(row[7].(string))

	if isLegacy {
		if isLimited {
			return nogard.Available
		} else {
			return nogard.Unavailable
		}
	}
	return nogard.Permanent
}

func (e *DragonEncyclopedia) rowToBreedingCombo(row []interface{}) []string {
	combo := row[75].(string)
	if combo == "" {
		return make([]string, 0)
	}
	return strings.Split(combo, "+")
}

func (e *DragonEncyclopedia) rowToIncubation(row []interface{}) time.Duration {
	seconds, _ := strconv.ParseInt(row[1].(string), 10, 64)
	return time.Duration(seconds) * time.Second
}

func (e *DragonEncyclopedia) rowToCash(row []interface{}) float64 {
	earn, _ := strconv.ParseFloat(row[34].(string), 64)
	return 60.0 / earn * 100
}

func (e *DragonEncyclopedia) rowToHabitats(row []interface{}) []nogard.Element {
	elStrMap := map[string]nogard.Element{
		"P":  nogard.ElPlant,
		"F":  nogard.ElFire,
		"E":  nogard.ElEarth,
		"C":  nogard.ElCold,
		"L":  nogard.ElLightning,
		"W":  nogard.ElWater,
		"A":  nogard.ElAir,
		"M":  nogard.ElMetal,
		"I":  nogard.ElLight,
		"D":  nogard.ElDark,
		"Ga": nogard.ElGalaxy,
		"R":  nogard.ElRift,
		"Rb": nogard.ElRainbow,
		"Ge": nogard.ElGemstone,
		"Cr": nogard.ElCrystalline,
		"Se": nogard.ElSeasonal,
		"Tr": nogard.ElTreasure,
		"Su": nogard.ElSun,
		"Mo": nogard.ElMoon,
		"Ol": nogard.ElOlympus,
		"Ap": nogard.ElApocalypse,
		"Dr": nogard.ElDream,
		"Sn": nogard.ElSnowflake,
		"Mh": nogard.ElMonolith,
		"Or": nogard.ElOrnamental,
		"Au": nogard.ElAura,
		"Ch": nogard.ElChrysalis,
		"Hi": nogard.ElHidden,
		"Sf": nogard.ElSurface,
		"Me": nogard.ElMelody,
	}

	elements := make([]nogard.Element, 0)

	re := regexp.MustCompile("[A-Z][^A-Z]*")

	segs := re.FindAllString(row[76].(string), -1)

	for _, seg := range segs {
		elements = append(elements, elStrMap[seg])
	}

	return elements
}

func (e *DragonEncyclopedia) mapRowToDragon(row []interface{}) *nogard.Dragon {
	return &nogard.Dragon{
		Name:         row[0].(string),
		Incubation:   e.rowToIncubation(row),
		Availability: e.rowToAvailability(row),
		Rarity:       e.rowToRarity(row),
		Combo:        e.rowToBreedingCombo(row),
		StartingCash: e.rowToCash(row),
		Habitats:     e.rowToHabitats(row),
	}
}
