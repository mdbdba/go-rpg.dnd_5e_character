package schema

import (
	"encoding/json"
	"fmt"
	"github.com/itchyny/timefmt-go"
	common "github.com/mdbdba/go_rpg_commonUtils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strconv"
	"time"
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
	return map[int]int{
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
	return map[string][]int{
		"predefined":        {},
		"strict":            {}, // 3d6
		"common":            {}, // 4d6 drop lowest
		"standard":          {15, 14, 13, 12, 10, 8},
		"pointbuy_even":     {13, 13, 13, 12, 12, 12},
		"pointbuy_onemax":   {15, 12, 12, 12, 11, 11},
		"pointbuy_twomax":   {15, 15, 11, 10, 10, 10},
		"pointbuy_threemax": {15, 15, 15, 8, 8, 8},
	}
}

// AbilityArrayTemplate is used to get a map with the abilities as keys
var AbilityArrayTemplate = func() map[string]int {
	return map[string]int{
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
func rollRawAbilitySlice(rollOption string,
	logger *zap.SugaredLogger) (rollSlice []int, err error) {
	// %s is The number of seconds since the Epoch
	nowStr := timefmt.Format(time.Now(), "%s")
	var rnd string
	rnd, err = common.GenerateRandomString(5)
	if err != nil {
		return
	}
	timesToRoll := 4
	options := []string{"drop lowest 1"}
	if rollOption == "strict" {
		timesToRoll = 3
		options = make([]string, 0)
	}
	for i := 0; i < 6; i++ {
		msg := fmt.Sprintf("{\"RawAbilitySlice\": \"%s-%s-%s-%d/6\"", nowStr,
			rnd, strconv.FormatInt(time.Now().UnixNano(), 10), i+1)
		r, err := common.Perform(6, timesToRoll, msg, options...)
		if err != nil {
			panic("Roll attempt failed")
		}
		//log the roll results, then harvest roll results
		rollSlice = append(rollSlice, r.Result)
		//Log the results
		logger.Infow("Roll", "Sides", r.Sides,
			"TimesToRoll", r.TimesToRoll,
			"RollsGenerated", common.IntSliceToString(r.RollsGenerated),
			"RollsUsed", common.IntSliceToString(r.RollsUsed),
			"Options", r.Options,
			"AdditiveValue", r.AdditiveValue)
	}
	rollSlice = common.SortDescendingIntSlice(rollSlice)
	return
}

// GetPreGeneratedBaseAbilityArray returns a Base Ability array based on a supplied
//  array that has an assumed order.  This will be used mostly for testing or
//  balance comparisons.  If a player has used this method we are expecting
//  this is an import of an existing player.  If not, it would be suspicious.
func GetPreGeneratedBaseAbilityArray(pre []int) (map[string]int, []string) {
	sortOrder := []string{
		"Strength",
		"Dexterity",
		"Constitution",
		"Intelligence",
		"Wisdom",
		"Charisma",
	}
	retObj := AbilityArrayTemplate()
	for i := 0; i < len(pre); i++ {
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
	return retObj, sortOrder
}

// GetBaseAbilityArray returns a generated Base Ability array and the unsorted
//  values that went into it. The values are generated depending on the
//  rollingOption passed (see AbilityAssign). How they are assigned to the
//  abilities depends on a sorting order provided by the sortSlice and
//  a rolling option.
func GetBaseAbilityArray(sortOrder []string, rollingOption string,
	logger *zap.SugaredLogger) (r map[string]int, rawValueSlice []int, err error) {
	r = AbilityArrayTemplate()
	lu := AbilityAssign()
	switch rollingOption {
	case "common":
		rawValueSlice, err = rollRawAbilitySlice(rollingOption, logger)
		fmt.Println(rawValueSlice)
		if err != nil {
			return
		}
	case "strict":
		rawValueSlice, err = rollRawAbilitySlice(rollingOption, logger)
		if err != nil {
			return
		}
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
	for i := 0; i < len(sortOrder); i++ {
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
	logger.Infow("Base Array Info",
		"rawValues", common.IntSliceToString(rawValueSlice),
		"sortedValues", AbilityMapToString(r),
	)
	return r, rawValueSlice, err
}

// AbilityArray is the struct that completely defines the Ability Array and
// all the pieces that make it up.
//  Where:
//    Raw are the values as they were originally generated
//    RollingOption describes how the Raw values were generated
//    SortOrder is the order applied to the Raw values to make Base
//    Base is the base point for the Ability scores
//    ArchetypeBonus are values that reflect racial/archetypal bonuses
//    ArchetypeBonusIgnored if true don't include any of the racial/archetypal bonuses
//    LevelChangeIncrease are values added when levels achieved
//    AdditionalBonus any other values that influence ability values
//    Values are the summation of Base + ArchetypeBonus (if used) +
//           LevelChangeIncrease + AdditionalBonus
//    Modifiers are the modifiers based on Values
//    CtxRef is the context reference for the assignment
//
type AbilityArray struct {
	Raw                   []int
	RollingOption         string
	SortOrder             []string
	Base                  map[string]int
	ArchetypeBonus        map[string]int
	ArchetypeBonusIgnored bool
	LevelChangeIncrease   map[string]int
	AdditionalBonus       map[string]int
	Values                map[string]int
	Modifiers             map[string]int
	CtxRef                string
	IsMonsterOrGod        bool
}

func GetPreGeneratedAbilityArray(Raw []int, ArchetypeBonus map[string]int,
	ArchetypeBonusIgnored bool, LevelChangeIncrease map[string]int,
	AdditionalBonus map[string]int, CtxRef string, IsMonsterOrGod bool) *AbilityArray {
	b, sortOrder := GetPreGeneratedBaseAbilityArray(Raw)
	values := AbilityArrayTemplate()
	mods := AbilityArrayTemplate()
	a := AbilityArray{
		Raw:                   Raw,
		RollingOption:         "pregenerated",
		SortOrder:             sortOrder,
		Base:                  b,
		ArchetypeBonus:        ArchetypeBonus,
		ArchetypeBonusIgnored: ArchetypeBonusIgnored,
		LevelChangeIncrease:   LevelChangeIncrease,
		AdditionalBonus:       AdditionalBonus,
		Values:                values,
		Modifiers:             mods,
		CtxRef:                CtxRef,
		IsMonsterOrGod:        IsMonsterOrGod,
	}
	a.setValuesAndModifiers()
	return &a
}

// GetAbilityArray is the function to use to get a Fully populated ability array for a
// character. The Ability Array struct will contain everything you need to build a
// character and all the info to know how it was all put together. It returns a pointer
// to an AbilityArray
//   Parameters:
//    RollingOption describes how the Raw values were generated
//    SortOrder is the order applied to the Raw values
//    ArchetypeBonus are values that reflect racial/archetypal bonuses
//    ArchetypeBonusIgnored if true don't include any of the racial/archetypal bonuses
//    LevelChangeIncrease are values added when levels achieved
//    AdditionalBonus any other values that influence ability values
//    CtxRef is the context reference for the assignment. A freetext
//      string that you can use to keep track of it in the logs.
//
func GetAbilityArray(RollingOption string,
	SortOrder []string, ArchetypeBonus map[string]int,
	ArchetypeBonusIgnored bool, LevelChangeIncrease map[string]int,
	AdditionalBonus map[string]int, CtxRef string, IsMonsterOrGod bool,
	logger *zap.SugaredLogger) (*AbilityArray, error) {
	b, raw, err := GetBaseAbilityArray(SortOrder, RollingOption, logger)
	if err != nil {
		return &AbilityArray{}, err
	}
	values := AbilityArrayTemplate()
	mods := AbilityArrayTemplate()
	a := AbilityArray{
		Raw:                   raw,
		RollingOption:         RollingOption,
		SortOrder:             SortOrder,
		Base:                  b,
		ArchetypeBonus:        ArchetypeBonus,
		ArchetypeBonusIgnored: ArchetypeBonusIgnored,
		LevelChangeIncrease:   LevelChangeIncrease,
		AdditionalBonus:       AdditionalBonus,
		Values:                values,
		Modifiers:             mods,
		CtxRef:                CtxRef,
		IsMonsterOrGod:        IsMonsterOrGod,
	}
	a.setValuesAndModifiers()
	logger.Infow("GetAbilityArray", zap.Object("AbilityArray", &a))
	return &a, nil
}

func (pa *AbilityArray) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("Raw", common.IntSliceToString(pa.Raw))
	enc.AddString("RollingOption", pa.RollingOption)
	enc.AddString("SortOrder", common.StringSliceToString(pa.SortOrder))
	enc.AddString("Base", AbilityMapToString(pa.Base))
	enc.AddString("ArchetypeBonus", AbilityMapToString(pa.ArchetypeBonus))
	enc.AddBool("ArchetypeBonusIgnored", pa.ArchetypeBonusIgnored)
	enc.AddString("LevelChangeIncrease", AbilityMapToString(pa.LevelChangeIncrease))
	enc.AddString("AdditionalBonus", AbilityMapToString(pa.AdditionalBonus))
	enc.AddString("Values", AbilityMapToString(pa.Values))
	enc.AddString("Modifiers", AbilityMapToString(pa.Modifiers))
	enc.AddString("CtxRef", pa.CtxRef)
	enc.AddBool("IsMonsterOrGod", pa.IsMonsterOrGod)
	return nil
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

// GetScore gets a single ability score from the array. Valid abilities are:
// Strength, Dexterity, Constitution, Intelligence, Wisdom, and Charisma
func (pa *AbilityArray) GetScore(ability string) (int, error) {
	if ValidateAbilityName(ability) {
		return pa.Values[ability], nil
	}
	return -10, fmt.Errorf("value passed for ability, %s, is not correct", ability)
}

// GetModifier gets a single modifier from the array. Valid abilities are:
// Strength, Dexterity, Constitution, Intelligence, Wisdom, and Charisma
func (pa *AbilityArray) GetModifier(ability string) (int, error) {
	if ValidateAbilityName(ability) {
		return pa.Modifiers[ability], nil
	}
	return -10, fmt.Errorf("value passed for ability, %s, is not correct", ability)
}

func (pa *AbilityArray) setValuesAndModifiers() {
	maxVal := 20
	if pa.IsMonsterOrGod {
		maxVal = 30 // Gods and Monsters can have ability scores up to 30
	}
	for k := range pa.Base {
		atb := pa.ArchetypeBonus[k]
		if pa.ArchetypeBonusIgnored {
			atb = 0
		}
		tVal := pa.Base[k] + atb + pa.LevelChangeIncrease[k] +
			pa.AdditionalBonus[k]
		// Values cannot exceed 20 or 30. Set that as max.
		if tVal > maxVal {
			tVal = maxVal
		}
		pa.Values[k] = tVal
	}

	lu := AbilityScoreModifier()
	for k, v := range pa.Values {
		pa.Modifiers[k] = lu[v]
	}
}

func ValidateAbilityName(ability string) bool {
	lu := AbilityArrayTemplate()
	if _, ok := lu[ability]; ok {
		return true
	}
	return false
}

// AdjustValues changes the totals in the maps within an AbilityArray
// and recalculates the total values.
//  Where:
//    ValueType is "ArchetypeBonus", "LevelChangeIncrease", or "AdditionalBonus".
//    Values is a map containing the adjustments to make
func (pa *AbilityArray) AdjustValues(ValueType string, Ability string,
	ChangeValue int, logger *zap.SugaredLogger) {
	switch ValueType {
	case "ArchetypeBonus":
		pa.ArchetypeBonus[Ability] += ChangeValue
	case "LevelChangeIncrease":
		pa.LevelChangeIncrease[Ability] += ChangeValue
	case "AdditionalBonus":
		pa.AdditionalBonus[Ability] += ChangeValue
	}
	pa.setValuesAndModifiers()
	logger.Infow("AdjustValues", zap.Object("AbilityArray", pa))
}

// AbilityMapToString converts a map keyed with the abilities to a single string.
func AbilityMapToString(src map[string]int) (tgt string) {
	tgt = fmt.Sprintf("{\"Strength\": %2d, \"Dexterity\": %2d, \"Constitution\": %2d, "+
		"\"Intelligence\": %2d, \"Wisdom\": %2d, \"Charisma\": %2d}",
		src["Strength"], src["Dexterity"], src["Constitution"], src["Intelligence"],
		src["Wisdom"], src["Charisma"])
	return
}

// ConvertToString converts an AbilityArray to a single string. The p argument
// is for pretty print where everything lines up.
func (pa *AbilityArray) ConvertToString(p bool) (s string) {
	rawStr := common.IntSliceToString(pa.Raw)
	orderStr := common.StringSliceToString(pa.SortOrder)
	baseStr := AbilityMapToString(pa.Base)
	archStr := AbilityMapToString(pa.ArchetypeBonus)
	lvlStr := AbilityMapToString(pa.LevelChangeIncrease)
	addbStr := AbilityMapToString(pa.AdditionalBonus)
	valStr := AbilityMapToString(pa.Values)
	modStr := AbilityMapToString(pa.Modifiers)
	pStr := ""
	f := "AbilityArray -- %sRaw: %s, %sRollingOption: %s, " +
		"%sSortOrder: %s, %sBaseArray: %s, %sArchetypeBonus: %s, " +
		"%sArchetypeBonusIgnored: %v, %sLevelChangeIncreases: %s, " +
		"%sAdditionalBonus: %s, %sValues: %s, %sModifiers: %s, %sCtxRef: %s, " +
		"%sIsMonsterOrGod: %v\n"
	if p {
		pStr = "\n\t"
		f = "AbilityArray -- %sRaw:                   %s, %sRollingOption:         %s, " +
			"%sSortOrder: %91s, %sBaseArray: %115s, %sArchetypeBonus: %110s, " +
			"%sArchetypeBonusIgnored: %v, %sLevelChangeIncreases:  %s, " +
			"%sAdditionalBonus: %109s, %sValues: %118s, %sModifiers: %115s, " +
			"%sCtxRef:                %s, %sIsMonsterOrGod:        %v\n"
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
		pStr, valStr,
		pStr, modStr,
		pStr, pa.CtxRef,
		pStr, pa.IsMonsterOrGod)
	return
}
