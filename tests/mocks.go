package tests

import (
	"errors"
	. "stddev-service/stddev"
)

type MockedNumGenerator struct{}

func (m *MockedNumGenerator) GetIntSlices(numberOfSlices, length int) ([][]int, error) {
	return [][]int{
		{0, 1, 2, 3, 6},
		{0, 0, -1, 5, 2},
	}, nil
}

type MockedErrorNumGenerator struct{}

func (m *MockedErrorNumGenerator) GetIntSlices(numberOfSlices, length int) ([][]int, error) {
	return nil, errors.New("error")
}

type MockedCalculator struct{}

func (m *MockedCalculator) CalculateStdDevs(intSlices [][]int) ([]StdDev, error) {
	return []StdDev{
		{Stddev: 2, Data: []int{0, 1, 2, 3, 6}},
		{Stddev: 2, Data: []int{-1, 0, 0, 2, 5}},
		{Stddev: 2, Data: []int{-1, 0, 0, 0, 1, 2, 2, 3, 5, 6}},
	}, nil
}

type MockedErrorCalculator struct{}

func (m *MockedErrorCalculator) CalculateStdDevs(intSlices [][]int) ([]StdDev, error) {
	return nil, errors.New("error")
}
