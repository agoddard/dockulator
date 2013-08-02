package models

import (
	"testing"
)

var cleaningTests = []struct {
	input string
	expected string
}{
	{
		"234 + 11",
		"234 + 11",
	},
	{
		"   42 * 11",
		"42 * 11",
	},
	{
		"11 / 10      ",
		"11 / 10",
	},
	{
		"     732 - 41    ",
		"732 - 41",
	},
	{
		"1+1",
		"1 + 1",
	},
	{
		"1       *4",
		"1 * 4",
	},
	{
		"1/      3",
		"1 / 3",
	},
	{
		"fji3spod+ff13dfs",
		"3 + 13",
	},
	{
		" sdf 1.2sdf41.123dsf +sdf 21sdf.44sdf.2 ",
		"1.241123 + 21.442",
	},
	{
		"ksjdflkj",
		"error",
	},
	{
		"skdfj2132 + sdlkfj",
		"error",
	},
}

func TestCleanCalculation(t *testing.T) {
	for i, test := range cleaningTests {
		val := CleanCalculation(test.input)
		if val != test.expected {
			t.Errorf(" %d test failed.\nGot:\n%v\nExpected:\n%v", i, val, test.expected)
		}
	}
}
