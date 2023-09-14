package stddev

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var testCases = []struct {
	name     string
	input    [][]int
	expected []StdDev
}{
	{
		"empty",
		[][]int{},
		[]StdDev{},
	},
	{
		"single slice",
		[][]int{
			{1, 2, 3, 4, 5},
		},
		[]StdDev{
			{Stddev: 1, Data: []int{1, 2, 3, 4, 5}},
			{Stddev: 1, Data: []int{1, 2, 3, 4, 5}},
		},
	},
	{
		"common",
		[][]int{
			{1, 2, 3, 4, 5},
			{10, 20, 30, 40, 50},
		},
		[]StdDev{
			{Stddev: 1, Data: []int{1, 2, 3, 4, 5}},
			{Stddev: 14, Data: []int{10, 20, 30, 40, 50}},
			{Stddev: 16, Data: []int{1, 2, 3, 4, 5, 10, 20, 30, 40, 50}},
		},
	},
	{
		"duplicates and minuses",
		[][]int{
			{-1, 2, -3, 4, -5, 0, 0, 1, 1},
			{10, -20, 30, -40, 50},
		},
		[]StdDev{
			{Stddev: 2, Data: []int{-5, -3, -1, 0, 0, 1, 1, 2, 4}},
			{Stddev: 32, Data: []int{-40, -20, 10, 30, 50}},
			{Stddev: 19, Data: []int{-40, -20, -5, -3, -1, 0, 0, 1, 1, 2, 4, 10, 30, 50}},
		},
	},
}

func TestStdDevsWithSumCalculator(t *testing.T) {
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			//given
			calculator := StdDevsWithSumCalculator{}
			//when
			result, err := calculator.CalculateStdDevs(testCase.input)
			//then
			assert.NoError(t, err)
			assert.ElementsMatch(t, testCase.expected, result)
		})
	}
}
