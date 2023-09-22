package numbers

import (
	"context"
	"fmt"
	"net/http"
)

// 1. spawn X goroutines with nonbuffered channel
// 2. iterate for X results

type IterateForXResultsRandomOrgNumbersGenerator struct {
	RandomOrgNumbersGenerator
}

func NewIterateForXResultsRandomOrgNumbersGenerator() *IterateForXResultsRandomOrgNumbersGenerator {
	return &IterateForXResultsRandomOrgNumbersGenerator{NewRandomOrgNumbersGenerator()}
}

func (g *IterateForXResultsRandomOrgNumbersGenerator) GetIntSlices(numberOfSlices, length int) ([][]int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), g.timeout)
	defer cancel()

	url := fmt.Sprintf("%s?num=%d&min=%d&max=%d&col=1&base=10&format=plain&rnd=new", g.url, length, g.min, g.max)

	intSlices := make([][]int, 0, numberOfSlices)
	results := make(chan FetchRandomNumbersResult, numberOfSlices)

	for i := 0; i < numberOfSlices; i++ {
		go g.fetchRandomNumbers(ctx, url, length, results)
	}

	for i := 0; i < numberOfSlices; i++ {
		result := <-results
		if result.err != nil {
			return nil, fmt.Errorf("could not get numbers: %v", result.err)
		}
		intSlices = append(intSlices, result.numbers)
	}

	return intSlices, nil
}

func (g *IterateForXResultsRandomOrgNumbersGenerator) fetchRandomNumbers(ctx context.Context, url string, length int, results chan<- FetchRandomNumbersResult) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		results <- FetchRandomNumbersResult{nil, fmt.Errorf("could not make a request context: %v", err)}
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		results <- FetchRandomNumbersResult{nil, fmt.Errorf("request failed: %v", err)}
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		results <- FetchRandomNumbersResult{nil, fmt.Errorf("got status code %d", resp.StatusCode)}
		return
	}

	ints, err := intsFromBody(resp.Body, length)
	if err != nil {
		results <- FetchRandomNumbersResult{nil, fmt.Errorf("could not get numbers from body: %v", err)}
		return
	}

	results <- FetchRandomNumbersResult{ints, nil}
}
