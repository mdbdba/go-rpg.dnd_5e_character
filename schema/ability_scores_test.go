package schema

import (
	"fmt"
	"github.com/stretchr/testify/assert"
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

func TestPreGeneratedAbilityArray(t *testing.T) {
	actual := PreGeneratedAbilityArray([]int{18,17,16,15,14,13})
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
}

func TestGetBaseAbilityArray(t *testing.T) {
	sortOrder := []string{"Dexterity","Constitution","Strength",
		"Charisma","Wisdom","Intelligence"}
	rollingOption := "standard"
	actual, r := GetBaseAbilityArray(sortOrder,rollingOption)
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
}

func TestGetAbilityArray(t *testing.T) {
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

	a := GetAbilityArray(rollingOption,sortOrder,ArchetypeBonus,
		ArchetypeBonusIgnored,LevelChangeIncrease, AdditionalBonus)
	fmt.Println(a.ToPrettyString())
	assert.Equal(t, 15 , a.Values["Strength"])
	assert.Equal(t, 14, a.Values["Charisma"])
	assert.Equal(t, 17, a.Values["Dexterity"])
	assert.Equal(t, 9, a.Values["Intelligence"])

}

func TestAdjustValues(t *testing.T) {
	rollingOption := "standard"
	sortOrder := []string{"Dexterity","Constitution","Strength",
		"Charisma","Wisdom","Intelligence"}
	ArchetypeBonus := AbilityArrayTemplate()
	ArchetypeBonusIgnored := false
	LevelChangeIncrease := AbilityArrayTemplate()
	AdditionalBonus := AbilityArrayTemplate()

	a := GetAbilityArray(rollingOption,sortOrder,ArchetypeBonus,
		ArchetypeBonusIgnored,LevelChangeIncrease, AdditionalBonus)
	a.AdjustValues("ArchetypeBonus", "Charisma", 2)
	a.AdjustValues("ArchetypeBonus", "Intelligence", 1)
	assert.Equal(t, 14, a.Values["Charisma"])
	assert.Equal(t, 9, a.Values["Intelligence"])
	a.AdjustValues("LevelChangeIncrease", "Dexterity", 2)
	assert.Equal(t, 17, a.Values["Dexterity"])
	a.AdjustValues("AdditionalBonus", "Strength",2)
	assert.Equal(t, 15 , a.Values["Strength"])

	fmt.Println(a.ToPrettyString())
}