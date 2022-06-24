// Code generated by "stringer -type Rarity"; DO NOT EDIT.

package nogard

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Primary-0]
	_ = x[Hybrid-1]
	_ = x[Rare-2]
	_ = x[Epic-3]
	_ = x[Gemstone-4]
	_ = x[Galaxy-5]
	_ = x[Mythic-6]
	_ = x[Legendary-7]
}

const _Rarity_name = "PrimaryHybridRareEpicGemstoneGalaxyMythicLegendary"

var _Rarity_index = [...]uint8{0, 7, 13, 17, 21, 29, 35, 41, 50}

func (i Rarity) String() string {
	if i >= Rarity(len(_Rarity_index)-1) {
		return "Rarity(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Rarity_name[_Rarity_index[i]:_Rarity_index[i+1]]
}
