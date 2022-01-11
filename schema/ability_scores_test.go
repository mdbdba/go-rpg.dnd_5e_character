package schema

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	"testing"
)

func TestAbilityDescriptions(t *testing.T) {
	actual := AbilityDescriptions()
	assert.Equal(t,6, len(actual))
}

func TestAbilityScoreModifier(t *testing.T) {
	actual := AbilityScoreModifier()
	assert.Equal(t, 30, len(actual))
	assert.Equal(t, 0, actual[0])
	assert.Equal(t, 10, actual[30])
}

func TestAbilityAssign(t *testing.T) {
	actual := AbilityAssign()
	actualKeys := GetAbilityRollingOptions()
	assert.Equal(t, 8, len(actual))
	assert.Equal(t, 8, len(actualKeys))
}

func TestAbilityArrayTemplate(t *testing.T) {
	actual := AbilityArrayTemplate()
	assert.Equal(t, 6, len(actual))
	assert.Equal(t, 0, actual["Charisma"])
}

func TestGetPreGeneratedBaseAbilityArray(t *testing.T) {
	actual, sortOrder := GetPreGeneratedBaseAbilityArray([]int{18,17,16,15,14,13})
	assert.Equal(t, 6, len(actual))
	expected := map[string]int{
		"Strength":     18,
		"Dexterity":    17,
		"Constitution": 16,
		"Intelligence": 15,
		"Wisdom":       14,
		"Charisma":     13,
	}
	assert.Equal(t, expected, actual)
	expectedSortOrder := []string {
		"Strength",
		"Dexterity",
		"Constitution",
		"Intelligence",
		"Wisdom",
		"Charisma",
	}
	assert.Equal(t, expectedSortOrder, sortOrder)
}

func TestGetBaseAbilityArray(t *testing.T) {
	// Given
	observedZapCore, observedLogs := observer.New(zap.InfoLevel)
	observedLoggerSugared := zap.New(observedZapCore).Sugar()
	sortOrder := []string{"Dexterity","Constitution","Strength",
		"Charisma","Wisdom","Intelligence"}
	rollingOption := "standard"

	// When
	actual, r := GetBaseAbilityArray(sortOrder,rollingOption, observedLoggerSugared)

	// Then
	assert.Equal(t, []int{15, 14, 13, 12, 10, 8}, r)
	expected := map[string]int{
		"Strength":     13,
		"Dexterity":    15,
		"Constitution": 14,
		"Intelligence": 8,
		"Wisdom":       10,
		"Charisma":     12,
	}
	assert.Equal(t, expected, actual)
	require.Equal(t, 1, observedLogs.Len())
}

func TestGetPreGeneratedAbilityArray(t *testing.T) {
	Raw := []int{18,17,16,15,14,13}
	ArchetypeBonus := AbilityArrayTemplate()
	ArchetypeBonus["Charisma"] = 2
	ArchetypeBonus["Intelligence"] = 1
	ArchetypeBonusIgnored := false
	LevelChangeIncrease := AbilityArrayTemplate()
	LevelChangeIncrease["Dexterity"] = 2
	AdditionalBonus := AbilityArrayTemplate()
	AdditionalBonus["Strength"] = 2
	ctxRef := "TestGetPreGeneratedAbilityArray"
	isMonsterOrGod := false
	a := GetPreGeneratedAbilityArray( Raw, ArchetypeBonus,
		ArchetypeBonusIgnored, LevelChangeIncrease,
		AdditionalBonus, ctxRef, isMonsterOrGod)
	fmt.Println(a.ToPrettyString())
	assert.Equal(t, 20 , a.Values["Strength"])
	assert.Equal(t, 15, a.Values["Charisma"])
	assert.Equal(t, 19, a.Values["Dexterity"])
	assert.Equal(t, 16, a.Values["Intelligence"])
	actual, _ := a.GetModifier("Strength")
	assert.Equal(t, 5, actual )
	actual, _ = a.GetModifier("Intelligence")
	assert.Equal(t, 3, actual )
}

func TestGetAbilityArray(t *testing.T) {
	// Given
	observedZapCore, observedLogs := observer.New(zap.InfoLevel)
	observedLoggerSugared := zap.New(observedZapCore).Sugar()
	rollingOption := "standard"
	sortOrder := []string{"Dexterity","Constitution","Strength",
		"Charisma","Wisdom","Intelligence"}
	ArchetypeBonus := AbilityArrayTemplate()
	ArchetypeBonus["Charisma"] = 2
	ArchetypeBonus["Intelligence"] = 1
	ArchetypeBonusIgnored := false
	LevelChangeIncrease := AbilityArrayTemplate()
	LevelChangeIncrease["Dexterity"] = 2
	AdditionalBonus := AbilityArrayTemplate()
	AdditionalBonus["Strength"] = 2
	ctxRef := "TestGetAbilityArray"
	isMonsterOrGod := false

	// When
	a := GetAbilityArray(rollingOption,sortOrder,ArchetypeBonus,
		ArchetypeBonusIgnored,LevelChangeIncrease, AdditionalBonus,
		ctxRef, isMonsterOrGod, observedLoggerSugared)

	// Then
	fmt.Println(a.ToPrettyString())
	assert.Equal(t, 15 , a.Values["Strength"])
	assert.Equal(t, 14, a.Values["Charisma"])
	assert.Equal(t, 17, a.Values["Dexterity"])
	assert.Equal(t, 9, a.Values["Intelligence"])
	actual, _ := a.GetModifier("Strength")
	assert.Equal(t, 2, actual )
	actual, _ = a.GetModifier("Intelligence")
	assert.Equal(t, -1, actual )
	allLogs := observedLogs.All()

	ctxMap := allLogs[0].ContextMap()
	tStr, _ := ctxMap["rawValues"].(string)
	assert.Equal(t, "[15, 14, 13, 12, 10, 8]", tStr)
	tStr, _ = ctxMap["sortedValues"].(string)
	assert.Equal(t, "{\"Strength\": 13, \"Dexterity\": 15, \"Constitution\": 14, "+
		"\"Intelligence\":  8, \"Wisdom\": 10, \"Charisma\": 12}", tStr)

	assert.Equal(t, "GetAbilityArray", allLogs[len(allLogs)-1].Message )
}

func TestAdjustValues(t *testing.T) {
	// Given
	observedZapCore, observedLogs := observer.New(zap.InfoLevel)
	observedLoggerSugared := zap.New(observedZapCore).Sugar()
	rollingOption := "standard"
	sortOrder := []string{"Dexterity","Constitution","Strength",
		"Charisma","Wisdom","Intelligence"}
	ArchetypeBonus := AbilityArrayTemplate()
	ArchetypeBonusIgnored := false
	LevelChangeIncrease := AbilityArrayTemplate()
	AdditionalBonus := AbilityArrayTemplate()
	ctxRef := "TestAdjustValues"
	isMonsterOrGod := false

	a := GetAbilityArray(rollingOption,sortOrder,ArchetypeBonus,
		ArchetypeBonusIgnored,LevelChangeIncrease, AdditionalBonus,
		ctxRef, isMonsterOrGod, observedLoggerSugared)
	a.AdjustValues("ArchetypeBonus", "Charisma",
		2, observedLoggerSugared)
	a.AdjustValues("ArchetypeBonus", "Intelligence",
		1, observedLoggerSugared)
	assert.Equal(t, 14, a.Values["Charisma"])
	assert.Equal(t, 9, a.Values["Intelligence"])
	a.AdjustValues("LevelChangeIncrease", "Dexterity",
		2, observedLoggerSugared)
	assert.Equal(t, 17, a.Values["Dexterity"])
	a.AdjustValues("AdditionalBonus", "Strength",
		2, observedLoggerSugared)
	assert.Equal(t, 15 , a.Values["Strength"])
	actual, _ := a.GetModifier("Strength")
	assert.Equal(t, 2, actual )
	actual, _ = a.GetModifier("Intelligence")
	assert.Equal(t, -1, actual )
	fmt.Println(a.ToPrettyString())

	allLogs := observedLogs.All()
	assert.Equal(t, "AdjustValues", allLogs[len(allLogs)-1].Message )
}