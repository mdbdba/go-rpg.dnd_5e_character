package schema

import (
	"encoding/json"
	"fmt"
	common "github.com/mdbdba/go_rpg_commonUtils"
)

// AbilityDescriptions returns a map[string]string holding all the ability
// descriptions
var AbilityDescriptions = func() map[string]string {
	return map[string]string{
		"Strength":     "measure of physical power",
		"Dexterity":    "measure of agility",
		"Constitution": "measure of endurance",
		"Intelligence": "measure of reasoning and memory",
		"Wisdom":       "measure of perception and insight",
		"Charisma":     "measure of personality force",
	}
}

// AbilityScoreModifier returns a map of modifiers for ability score rolls.
var AbilityScoreModifier = func() map[int]int {
	return map[int]int {
		1: -5, 2: -4, 3: -4,
		4: -3, 5: -3,
		6: -2, 7: -2,
		8: -1, 9: -1,
		10: 0, 11: 0,
		12: 1, 13: 1,
		14: 2, 15: 2,
		16: 3, 17: 3,
		18: 4, 19: 4,
		20: 5, 21: 5,
		22: 6, 23: 6,
		24: 7, 25: 7,
		26: 8, 27: 8,
		28: 9, 29: 9, 30: 10,
	}
}

// AbilityAssign returns a map with all the options for "rolling" the
// ability values and in the case of set ones, the values to be used.
var AbilityAssign = func() map[string][]int {
	return map[string][]int {
		"predefined": {},
		"strict": {}, // 3d6
		"common": {}, // 4d6 drop lowest
		"standard": {15, 14, 13, 12, 10, 8},
		"pointbuy_even": {13, 13, 13, 12, 12, 12},
		"pointbuy_onemax": {15, 12, 12, 12, 11, 11},
		"pointbuy_twomax": {15, 15, 11, 10, 10, 10},
		"pointbuy_threemax": {15, 15, 15, 8, 8, 8},
	}
}

// AbilityArrayTemplate is used to get a map with the abilities as keys
var AbilityArrayTemplate = func() map[string]int {
	return map[string]int {
		"Strength":     0,
		"Dexterity":    0,
		"Constitution": 0,
		"Intelligence": 0,
		"Wisdom":       0,
		"Charisma":     0,
	}
}

// GetAbilityRollingOptions returns a slice of strings getting the
// possible values to pass for "rolling" options.
func GetAbilityRollingOptions() (options []string) {
	a := AbilityAssign()
	for k := range a {
		options = append(options, k)
	}
	return
}


// rollRawAbilitySlice rolls up the slice of ints to be used in the
//  ability array generation for the "strict" and "common" roll options.
//  Where:
//    strict = 3d6
//    common = 4d6 drop lowest 1
//  The rest of the options are set values defined in AbilityAssign
func rollRawAbilitySlice(rollOption string) (rollSlice []int) {
	timesToRoll := 4
	options := []string{"drop lowest 1"}
	if rollOption == "strict" {
		timesToRoll = 3
		options = make([]string,0)
	}
	for i := 0; i < 6; i++ {
		r, err := common.Perform(6, timesToRoll, options...)
		if err != nil {
			panic("Roll attempt failed")
		}
		//log the roll results, then harvest roll results
		rollSlice = append(rollSlice, r.Result)
	}
	return
}

// PreGeneratedAbilityArray returns a Base Ability array based on a supplied
//  array that has an assumed order.  This will be used mostly for testing or
//  balance comparisons.  If a player has used this method we are expecting
//  this is an import of an existing player.  If not, it would be suspicious.
func PreGeneratedAbilityArray(pre []int)  map[string]int {
	retObj := AbilityArrayTemplate()
	for i:=0; i<len(pre); i++ {
		switch i {
		case 0:
			retObj["Strength"] = pre[i]
		case 1:
			retObj["Dexterity"] = pre[i]
		case 2:
			retObj["Constitution"] = pre[i]
		case 3:
			retObj["Intelligence"] = pre[i]
		case 4:
			retObj["Wisdom"] = pre[i]
		case 5:
			retObj["Charisma"] = pre[i]
		}
	}
	return retObj
}

// GetBaseAbilityArray returns a generated Base Ability array and the unsorted
//  values that went into it. The values are generated depending on the
//  rollingOption passed (see AbilityAssign). How they are assigned to the
//  abilities depends on a sorting order provided by the sortSlice and
//  a rolling option.
func GetBaseAbilityArray(sortOrder []string,
	rollingOption string) ( r map[string]int, rawValueSlice []int) {
	r = AbilityArrayTemplate()
	lu := AbilityAssign()
	switch rollingOption {
	case "common":
		rawValueSlice = rollRawAbilitySlice(rollingOption)
	case "strict":
		rawValueSlice = rollRawAbilitySlice(rollingOption)
	case "standard":
		rawValueSlice = lu["standard"]
	case "pointbuy_even":
		rawValueSlice = lu["pointbuy_even"]
	case "pointbuy_onemax":
		rawValueSlice = lu["pointbuy_onemax"]
	case "pointbuy_twomax":
		rawValueSlice = lu["pointbuy_twomax"]
	case "pointbuy_threemax":
		rawValueSlice = lu["pointbuy_threemax"]
	}
	for i := 0; i <len(sortOrder); i++ {
		switch sortOrder[i] {
		case "Strength":
			r["Strength"] = rawValueSlice[i]
		case "Dexterity":
			r["Dexterity"] = rawValueSlice[i]
		case "Constitution":
			r["Constitution"] = rawValueSlice[i]
		case "Intelligence":
			r["Intelligence"] = rawValueSlice[i]
		case "Wisdom":
			r["Wisdom"] = rawValueSlice[i]
		case "Charisma":
			r["Charisma"] = rawValueSlice[i]
		}
	}
	return r, rawValueSlice
}

// AbilityArray is the struct that completely defines the Ability Array and
// all the pieces that make it up.
//  Where:
//    Raw are the values as they were originally generated
//    RollingOption describes how the Raw values were generated
//    SortOrder is the order applied to the Raw values to make Base
//    Base is the base point for the Ability scores
//    ArchetypeBonus values that reflect racial/archetypal bonuses
//    ArchetypeBonusIgnored if true don't include any of the racial/archetypal bonuses
//    LevelChangeIncrease values added when levels achieved
//    AdditionalBonus any other values that influence ability values
//    Values the summation of Base + ArchetypeBonus (if used) +
//           LevelChangeIncrease + AdditionalBonus
type AbilityArray struct {
	Raw []int
	RollingOption string
	SortOrder []string
	Base map[string]int
	ArchetypeBonus map[string]int
	ArchetypeBonusIgnored bool
	LevelChangeIncrease map[string]int
	AdditionalBonus map[string]int
	Values map[string]int
}
// GetAbilityArray is the function that pulls all the pieces together.
func GetAbilityArray(RollingOption string,
	SortOrder []string, ArchetypeBonus map[string]int,
	ArchetypeBonusIgnored bool, LevelChangeIncrease map[string]int,
	AdditionalBonus map[string]int) *AbilityArray {
	b, raw := GetBaseAbilityArray(SortOrder, RollingOption)
	values := AbilityArrayTemplate()
	for k, _ := range values {
		atb := ArchetypeBonus[k]
		if ArchetypeBonusIgnored {
			atb = 0
		}
		values[k] = b[k] + atb + LevelChangeIncrease[k] + AdditionalBonus[k]
	}
	return &AbilityArray{
		Raw:                   raw,
		RollingOption:         RollingOption,
		SortOrder:             SortOrder,
		Base:                  b,
		ArchetypeBonus:        ArchetypeBonus,
		ArchetypeBonusIgnored: ArchetypeBonusIgnored,
		LevelChangeIncrease:   LevelChangeIncrease,
		AdditionalBonus:       AdditionalBonus,
		Values:                values,
	}

}

func (pa *AbilityArray) ToJson() string {
    j, err := json.Marshal(pa)
	if err != nil {
		panic("Issue converting Ability Array to JSON.")
	}
	return string(j)
}

func (pa *AbilityArray) ToPrettyString() string {
	return pa.ConvertToString(true)
}

func (pa *AbilityArray) ToString() string {
	return pa.ConvertToString(false)
}

func MapStringIntToString(src map[string]int) (tgt string) {
	tgt = "["
	firstLoop := true
	for k,v := range src {
		joinChr := ", "
		if firstLoop {
			joinChr = ""
			firstLoop = false
		}
		tgt = fmt.Sprintf("%s%s\"%s\": %d",tgt,joinChr,k,v)
	}
	tgt = fmt.Sprintf("%s]", tgt)
	return
}

func AbilityMapToString(src map[string]int) (tgt string) {
	tgt = fmt.Sprintf("[\"Strength\": %2d, \"Dexterity\": %2d, \"Constitution\": %2d, " +
		"\"Intelligence\": %2d, \"Wisdom\": %2d, \"Charisma\": %2d]",
		src["Strength"],src["Dexterity"],src["Constitution"],src["Intelligence"],
		src["Wisdom"],src["Charisma"])
	return
}

func StringSliceToString(src []string) (tgt string) {
	tgt ="["
	for i:=0; i<len(src);i++ {
		joinChr := ", "
		if i == 0 {
					joinChr = ""
			}
		tgt = fmt.Sprintf("%s%s\"%s\"", tgt,joinChr,src[i])
	}
	tgt = fmt.Sprintf("%s]", tgt)
	return
}
func (pa *AbilityArray) ConvertToString(p bool) (s string) {
	rawStr := common.IntSliceToString(pa.Raw)
	orderStr := StringSliceToString(pa.SortOrder)
	baseStr := AbilityMapToString(pa.Base)
	archStr := AbilityMapToString(pa.ArchetypeBonus)
	lvlStr := AbilityMapToString(pa.LevelChangeIncrease)
	addbStr := AbilityMapToString(pa.AdditionalBonus)
	valStr := AbilityMapToString(pa.Values)
	pStr := ""
	f := "AbilityArray -- %sRaw: %s, %sRollingOption: %s, " +
		"%sSortOrder: %s, %sBaseArray: %s, %sArchetypeBonus: %s, " +
		"%sArchetypeBonusIgnored: %v, %sLevelChangeIncreases: %s, " +
		"%sAdditionalBonus: %s, %sValues: %s\n"
	if p {
		pStr = "\n\t"
		f = "AbilityArray -- %sRaw: %s, %sRollingOption: %s, " +
			"%sSortOrder: %90s, %sBaseArray: %114s, %sArchetypeBonus: %109s, " +
			"%sArchetypeBonusIgnored: %v, %sLevelChangeIncreases: %s, " +
			"%sAdditionalBonus: %108s, %sValues: %117s\n"
	}
	s = fmt.Sprintf(f,
		pStr, rawStr,
		pStr, pa.RollingOption,
		pStr, orderStr,
		pStr, baseStr,
		pStr, archStr,
		pStr, pa.ArchetypeBonusIgnored,
		pStr, lvlStr,
		pStr, addbStr,
		pStr, valStr)
	return
}


