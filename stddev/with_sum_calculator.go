package stddev

import (
	"fmt"
	"github.com/montanaflynn/stats"
	"sort"
)

type StdDevsWithSumCalculator struct{}

func (c *StdDevsWithSumCalculator) CalculateStdDevs(intSlices [][]int) ([]StdDev, error) {
	if len(intSlices) == 0 || len(intSlices[0]) == 0 {
		return []StdDev{}, nil
	}

	stdDevs := make([]StdDev, 0, len(intSlices)+1)
	for _, ints := range intSlices {
		sdev, err := stdDev(ints)
		if err != nil {
			return nil, err
		}
		stdDevs = append(stdDevs, sdev)
	}
	sdevOfsum, err := stdDevOfSumOfIntSlices(intSlices)
	if err != nil {
		return nil, err
	}
	stdDevs = append(stdDevs, sdevOfsum)

	return stdDevs, nil
}

func stdDevOfSumOfIntSlices(intSlices [][]int) (StdDev, error) {
	allInts := make([]int, 0)
	for _, slice := range intSlices {
		allInts = append(allInts, slice...)
	}
	sdev, err := stdDev(allInts)
	if err != nil {
		return StdDev{}, err
	}
	return sdev, nil
}

func stdDev(ints []int) (StdDev, error) {
	sort.Ints(ints)
	floats := intsToFloats(ints)
	sdev, err := stats.StandardDeviation(floats)
	if err != nil {
		return StdDev{}, fmt.Errorf("could not calculate stddev: %v", err)
	}
	return StdDev{Stddev: int(sdev), Data: ints}, nil
}

func intsToFloats(ints []int) []float64 {
	floats := make([]float64, 0, len(ints))
	for _, val := range ints {
		floats = append(floats, float64(val))
	}
	return floats
}
