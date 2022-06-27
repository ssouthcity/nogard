// Code generated by "stringer -type Element -trimprefix=El"; DO NOT EDIT.

package nogard

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ElPlant-0]
	_ = x[ElFire-1]
	_ = x[ElEarth-2]
	_ = x[ElCold-3]
	_ = x[ElLightning-4]
	_ = x[ElWater-5]
	_ = x[ElAir-6]
	_ = x[ElMetal-7]
	_ = x[ElLight-8]
	_ = x[ElDark-9]
	_ = x[ElGalaxy-10]
	_ = x[ElRift-11]
	_ = x[ElRainbow-12]
	_ = x[ElGemstone-13]
	_ = x[ElCrystalline-14]
	_ = x[ElSeasonal-15]
	_ = x[ElTreasure-16]
	_ = x[ElSun-17]
	_ = x[ElMoon-18]
	_ = x[ElOlympus-19]
	_ = x[ElApocalypse-20]
	_ = x[ElDream-21]
	_ = x[ElSnowflake-22]
	_ = x[ElMonolith-23]
	_ = x[ElOrnamental-24]
	_ = x[ElAura-25]
	_ = x[ElChrysalis-26]
	_ = x[ElHidden-27]
	_ = x[ElSurface-28]
	_ = x[ElMelody-29]
	_ = x[ElMythic-30]
	_ = x[ElLegendary-31]
}

const _Element_name = "PlantFireEarthColdLightningWaterAirMetalLightDarkGalaxyRiftRainbowGemstoneCrystallineSeasonalTreasureSunMoonOlympusApocalypseDreamSnowflakeMonolithOrnamentalAuraChrysalisHiddenSurfaceMelodyMythicLegendary"

var _Element_index = [...]uint8{0, 5, 9, 14, 18, 27, 32, 35, 40, 45, 49, 55, 59, 66, 74, 85, 93, 101, 104, 108, 115, 125, 130, 139, 147, 157, 161, 170, 176, 183, 189, 195, 204}

func (i Element) String() string {
	if i < 0 || i >= Element(len(_Element_index)-1) {
		return "Element(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Element_name[_Element_index[i]:_Element_index[i+1]]
}