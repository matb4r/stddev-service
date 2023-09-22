package numbers

import (
	"context"
	"fmt"
	"net/http"
	"sync"
)

// 1. spawn 5 goroutines with buffered channel
// 2. wait for them to finish
// 3. close the channel
// 4. iterate for all results from the channel

type BufferedChannelRandomOrgNumbersGenerator struct {
	RandomOrgNumbersGenerator
}

func NewBufferedChannelRandomOrgNumbersGenerator() *BufferedChannelRandomOrgNumbersGenerator {
	return &BufferedChannelRandomOrgNumbersGenerator{NewRandomOrgNumbersGenerator()}
}

func (g *BufferedChannelRandomOrgNumbersGenerator) GetIntSlices(numberOfSlices, length int) ([][]int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), g.timeout)
	defer cancel()

	url := fmt.Sprintf("%s?num=%d&min=%d&max=%d&col=1&base=10&format=plain&rnd=new", g.url, length, g.min, g.max)

	intSlices := make([][]int, 0, numberOfSlices)
	results := make(chan FetchRandomNumbersResult, numberOfSlices)
	var wg sync.WaitGroup

	for i := 0; i < numberOfSlices; i++ {
		wg.Add(1)
		go g.fetchRandomNumbers(ctx, &wg, url, length, results)
	}
	wg.Wait()
	close(results)

	for result := range results {
		if result.err != nil {
			return nil, fmt.Errorf("could not get numbers: %v", result.err)
		}
		intSlices = append(intSlices, result.numbers)
	}

	return intSlices, nil
}

func (g *BufferedChannelRandomOrgNumbersGenerator) fetchRandomNumbers(ctx context.Context, wg *sync.WaitGroup, url string, length int, results chan<- FetchRandomNumbersResult) {
	defer wg.Done()
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
