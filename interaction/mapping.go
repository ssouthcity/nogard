package interaction

import (
	"github.com/bwmarrin/discordgo"
	"github.com/ssouthcity/nogard"
)

var RarityColors = map[nogard.Rarity]int{
	nogard.Primary:   15663085,
	nogard.Hybrid:    9043831,
	nogard.Rare:      39423,
	nogard.Epic:      9048814,
	nogard.Gemstone:  16742161,
	nogard.Galaxy:    122367,
	nogard.Mythic:    16720827,
	nogard.Legendary: 15606273,
}

var ElementEmotes = map[nogard.Element]*discordgo.Emoji{
	nogard.ElPlant:       {Name: "plant_element", ID: "990774338103955476"},
	nogard.ElAir:         {Name: "air_element", ID: "990774327521738822"},
	nogard.ElApocalypse:  {Name: "apocalypse_element", ID: "990778072078905374"},
	nogard.ElAura:        {Name: "aura_element", ID: "990779249826869279"},
	nogard.ElChrysalis:   {Name: "chrysalis_element", ID: "990779251022245928"},
	nogard.ElCold:        {Name: "cold_element", ID: "990774328561913956"},
	nogard.ElCrystalline: {Name: "crystalline_element", ID: "990778073140043858"},
	nogard.ElDark:        {Name: "dark_element", ID: "990774329564356608"},
	nogard.ElDream:       {Name: "dream_element", ID: "990778074293497866"},
	nogard.ElEarth:       {Name: "earth_element", ID: "990774330877169724"},
	nogard.ElMonolith:    {Name: "monolith_element", ID: "990779254331555901"},
	nogard.ElMetal:       {Name: "metal_element", ID: "990774336988274718"},
	nogard.ElMelody:      {Name: "melody_element", ID: "990779253668859994"},
	nogard.ElLightning:   {Name: "lightning_element", ID: "990774335973249074"},
	nogard.ElLight:       {Name: "light_element", ID: "990774334920482817"},
	nogard.ElHidden:      {Name: "hidden_element", ID: "990779252297330718"},
	nogard.ElGemstone:    {Name: "gemstone_element", ID: "990778075409182770"},
	nogard.ElGalaxy:      {Name: "galaxy_element", ID: "990774333716721715"},
	nogard.ElFire:        {Name: "fire_element", ID: "990774332198363216"},
	nogard.ElMoon:        {Name: "moon_element", ID: "990778076516474940"},
	nogard.ElOlympus:     {Name: "olympus_element", ID: "990778077527285760"},
	nogard.ElOrnamental:  {Name: "ornamental_element", ID: "990779255753408552"},
	nogard.ElRainbow:     {Name: "rainbow_element", ID: "990778078680739910"},
	nogard.ElRift:        {Name: "rift_element", ID: "990774339232202772"},
	nogard.ElSeasonal:    {Name: "seasonal_element", ID: "990778079796416612"},
	nogard.ElSnowflake:   {Name: "snowflake_element", ID: "990778080433950723"},
	nogard.ElSun:         {Name: "sun_element", ID: "990778081939685436"},
	nogard.ElWater:       {Name: "water_element", ID: "990774340196892672"},
	nogard.ElTreasure:    {Name: "treasure_element", ID: "990778083470610443"},
	nogard.ElSurface:     {Name: "surface_element", ID: "990779257074643016"},
}
