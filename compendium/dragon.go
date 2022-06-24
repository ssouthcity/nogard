package compendium

import (
	"errors"
	"fmt"
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

func (e *DragonEncyclopedia) SearchDragons(query string) ([]string, error) {
	var dragons []string

	if err := e.setCategory("Compendium", "All"); err != nil {
		return dragons, err
	}

	resp, err := e.sheet.Spreadsheets.Values.Get(e.spreadsheetID, "Dragonarium!B6:C").Do()
	if err != nil {
		return dragons, err
	}

	for _, row := range resp.Values {
		n := row[0].(string)

		if strings.Contains(strings.ToLower(n), strings.ToLower(query)) {
			dragons = append(dragons, n)
		}
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

func (e *DragonEncyclopedia) mapRowToDragon(row *sheets.RowData) *nogard.Dragon {
	d := &nogard.Dragon{}

	d.Name = row.Values[0].FormattedValue
	d.Rarity = e.cellToRarity(row.Values[1])
	d.Availability = e.cellToAvailability(row.Values[2])
	d.Combo = e.cellToBreedingCombo(row.Values[3])
	d.Incubation = e.cellToIncubation(row.Values[4])

	return d
}