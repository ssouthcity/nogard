package nogard

import (
	"testing"
	"time"
)

func TestBreedingTime(t *testing.T) {
	t.Parallel()

	expectedValues := map[string]string{
		"26h0m0s": "20h48m0s", // lullaby dragon
		"24h0m0s": "19h12m0s", // decay dragon
		"17h0m0s": "13h36m0s", // grove dragon
		"11h0m0s": "8h48m0s",  // smolder dragon
		"14h0m0s": "11h12m0s", // limulimu dragon
	}

	for r, u := range expectedValues {
		regularTime, err := time.ParseDuration(r)
		if err != nil {
			t.Errorf("invalid duration expression; got error '%s'", err)
		}

		d := &Dragon{Incubation: regularTime}

		if time := d.BreedingTime(false).String(); time != r {
			t.Errorf("regular breeding time incorrect; got %s, expected %s", time, r)
		}

		if time := d.BreedingTime(true).String(); time != u {
			t.Errorf("upgraded breeding time incorrect; got %s, expected %s", time, u)
		}
	}
}

func TestDragonCash(t *testing.T) {
	t.Parallel()

	expectedValues := map[float64][]int{
		5.91:  {6, 9, 13, 17, 20, 24, 27, 31, 34, 38, 41, 45, 48, 52, 56, 59, 63, 66, 70, 73},                  // plains dragon
		15.0:  {15, 24, 33, 42, 51, 60, 69, 78, 87, 96, 105, 114, 123, 132, 141, 150, 159, 168, 177, 186},      // cumberpatch dragon
		6.67:  {7, 11, 15, 19, 23, 27, 31, 35, 39, 43, 47, 51, 55, 59, 63, 67, 71, 75, 79, 83},                 // mirror dragon
		12:    {12, 19, 26, 34, 41, 48, 55, 62, 70, 77, 84, 91, 98, 106, 113, 120, 127, 134, 142, 149},         // mistletoe dragon
		27.27: {27, 44, 60, 76, 93, 109, 125, 142, 158, 175, 191, 207, 224, 240, 256, 273, 289, 305, 322, 338}, // fluffles dragon
	}

	for sc, levels := range expectedValues {
		d := &Dragon{StartingCash: sc}

		for l, expectedCPM := range levels {
			cpm := d.CashPerMinute(l + 1)

			if cpm != expectedCPM {
				t.Errorf("expected %d, got %d", expectedCPM, cpm)
			}
		}
	}
}
