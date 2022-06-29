package compendium

import (
	"errors"
	"fmt"
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

func (e *DragonEncyclopedia) setCategory(ctype string, category string) error {
	_, err := e.sheet.Spreadsheets.Values.Update(e.spreadsheetID, "Dragonarium!E1:E2", &sheets.ValueRange{
		Values: [][]interface{}{{ctype}, {category}},
	}).ValueInputOption("USER_ENTERED").Do()

	return err
}

func (e *DragonEncyclopedia) DragonNames() ([]string, error) {
	if err := e.setCategory("Compendium", "All"); err != nil {
		return nil, err
	}

	resp, err := e.sheet.Spreadsheets.Values.Get(e.spreadsheetID, "Dragonarium!B6:C").Do()
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
	if err := e.setCategory("Compendium", "All"); err != nil {
		return nil, err
	}

	resp, err := e.sheet.Spreadsheets.Get(e.spreadsheetID).Ranges("Dragonarium!B6:Q").IncludeGridData(true).Do()
	if err != nil {
		return nil, err
	}

	for _, row := range resp.Sheets[0].Data[0].RowData {
		n := row.Values[0].EffectiveValue.StringValue

		if strings.EqualFold(*n, name) {
			return e.mapRowToDragon(row), nil
		}
	}

	return nil, errors.New("dragon not found")
}

func (e *DragonEncyclopedia) cellToRarity(cell *sheets.CellData) nogard.Rarity {
	rarityMap := map[string]nogard.Rarity{"Primary Rift": nogard.Rare}
	for r := nogard.Primary; r <= nogard.Legendary; r++ {
		rarityMap[r.String()] = r
	}
	return rarityMap[cell.FormattedValue]
}

func (e *DragonEncyclopedia) cellToAvailability(cell *sheets.CellData) nogard.Availability {
	if cell.FormattedValue == "Y" {
		if cell.EffectiveFormat.TextFormat.ForegroundColor.Green == 1.0 {
			return nogard.Available
		} else {
			return nogard.Unavailable
		}
	}
	return nogard.Permanent
}

func (e *DragonEncyclopedia) cellToBreedingCombo(cell *sheets.CellData) []string {
	if cell.FormattedValue == "" {
		return make([]string, 0)
	}
	return strings.Split(cell.FormattedValue, "+")
}

func (e *DragonEncyclopedia) cellToIncubation(cell *sheets.CellData) time.Duration {
	segs := strings.Split(cell.FormattedValue, ":")

	durStr := fmt.Sprintf("%sh%sm%ss", segs[0], segs[1], segs[2])
	dur, _ := time.ParseDuration(durStr)

	return dur
}

func (e *DragonEncyclopedia) cellToCash(cell *sheets.CellData) float64 {
	v, _ := strconv.ParseFloat(cell.FormattedValue, 64)
	return v
}

func (e *DragonEncyclopedia) cellToHabitats(cell *sheets.CellData) []nogard.Element {
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

	segs := re.FindAllString(cell.FormattedValue, -1)

	for _, seg := range segs {
		elements = append(elements, elStrMap[seg])
	}

	return elements
}

func (e *DragonEncyclopedia) mapRowToDragon(row *sheets.RowData) *nogard.Dragon {
	d := &nogard.Dragon{}

	d.Name = row.Values[0].FormattedValue
	d.Rarity = e.cellToRarity(row.Values[1])
	d.Availability = e.cellToAvailability(row.Values[2])
	d.Combo = e.cellToBreedingCombo(row.Values[3])
	d.Incubation = e.cellToIncubation(row.Values[4])
	d.Habitats = e.cellToHabitats(row.Values[6])
	d.StartingCash = e.cellToCash(row.Values[12])

	return d
}
