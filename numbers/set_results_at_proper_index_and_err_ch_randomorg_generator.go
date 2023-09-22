package numbers

import (
	"context"
	"fmt"
	"net/http"
	"sync"
)

// 1. spawn 5 goroutines with a nonbuffered channel, error channel and a slice with individual index
// 2. every goroutine set their result in the slice on the given index
// 3. if error, then goroutine sends it to the error channel
// 4. wait all goroutines to finish
// 5. check if there is an error in the error channel - and return proper result

type SetResultsAtProperIndexAndErrChannelRandomOrgGenerator struct {
	RandomOrgNumbersGenerator
}

func NewSetResultsAtProperIndexAndErrChannelRandomOrgGenerator() *SetResultsAtProperIndexAndErrChannelRandomOrgGenerator {
	return &SetResultsAtProperIndexAndErrChannelRandomOrgGenerator{NewRandomOrgNumbersGenerator()}
}

func (g *SetResultsAtProperIndexAndErrChannelRandomOrgGenerator) GetIntSlices(numberOfSlices, length int) ([][]int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), g.timeout)
	defer cancel()

	url := fmt.Sprintf("%s?num=%d&min=%d&max=%d&col=1&base=10&format=plain&rnd=new", g.url, length, g.min, g.max)

	intSlices := make([][]int, numberOfSlices)
	var wg sync.WaitGroup
	errCh := make(chan error)

	for i := 0; i < numberOfSlices; i++ {
		go g.fetchRandomNumbers(ctx, &wg, errCh, i, intSlices, url, length)
		wg.Add(1)
	}
	wg.Wait()

	select {
	case err := <-errCh:
		return nil, err
	default:
		return intSlices, nil
	}
}

func (g *SetResultsAtProperIndexAndErrChannelRandomOrgGenerator) fetchRandomNumbers(ctx context.Context, wg *sync.WaitGroup, errCh chan<- error, i int, intSlices [][]int, url string, length int) {
	defer wg.Done()
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		errCh <- fmt.Errorf("could not make a request context: %v", err)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		errCh <- fmt.Errorf("request failed: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errCh <- fmt.Errorf("got status code %d", resp.StatusCode)
		return
	}

	ints, err := intsFromBody(resp.Body, length)
	if err != nil {
		errCh <- fmt.Errorf("could not get numbers from body: %v", err)
		return
	}
	intSlices[i] = ints
}
