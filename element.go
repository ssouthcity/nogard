package nogard

type Element int

const (
	ElPlant Element = iota
	ElFire
	ElEarth
	ElCold
	ElLightning
	ElWater
	ElAir
	ElMetal
	ElLight
	ElDark

	ElGalaxy
	ElRift
	ElRainbow
	ElGemstone
	ElCrystalline
	ElSeasonal
	ElTreasure
	ElSun
	ElMoon
	ElOlympus
	ElApocalypse
	ElDream
	ElSnowflake
	ElMonolith
	ElOrnamental
	ElAura
	ElChrysalis
	ElHidden
	ElSurface
	ElMelody
	ElMythic
	ElLegendary
)

//go:generate stringer -type Element -trimprefix=El
