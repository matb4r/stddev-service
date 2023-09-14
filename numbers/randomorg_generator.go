package numbers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type FetchRandomNumbersResult struct {
	numbers []int
	err     error
}

type RandomOrgNumbersGenerator struct {
	url      string
	min, max int
	timeout  time.Duration
}

func NewRandomOrgNumbersGenerator() Generator {
	return &RandomOrgNumbersGenerator{
		url:     "https://www.random.org/integers/",
		min:     1,
		max:     100,
		timeout: 3000 * time.Millisecond,
	}
}

var client = http.DefaultClient

func (g *RandomOrgNumbersGenerator) GetIntSlices(numberOfSlices, length int) ([][]int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), g.timeout)
	defer cancel()

	url := fmt.Sprintf("%s?num=%d&min=%d&max=%d&col=1&base=10&format=plain&rnd=new", g.url, length, g.min, g.max)

	intSlices := make([][]int, 0, numberOfSlices)
	results := make(chan FetchRandomNumbersResult, numberOfSlices)
	var wg sync.WaitGroup

	for i := 0; i < numberOfSlices; i++ {
		wg.Add(1)
		go fetchRandomNumbers(ctx, &wg, url, length, results)
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

func fetchRandomNumbers(ctx context.Context, wg *sync.WaitGroup, url string, length int, results chan<- FetchRandomNumbersResult) {
	defer wg.Done()
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		results <- FetchRandomNumbersResult{nil, fmt.Errorf("could not make a request context: %v", err)}
		return
	}

	resp, err := client.Do(req)
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

func intsFromBody(body io.ReadCloser, length int) ([]int, error) {
	bytes, err := io.ReadAll(body)
	if err != nil {
		return nil, fmt.Errorf("could not read response body: %v", err)

	}
	err = body.Close()
	if err != nil {
		return nil, fmt.Errorf("could not read response body: %v", err)
	}
	s := string(bytes)
	ss := strings.Split(s, "\n")[:length]
	ints := make([]int, 0, length)
	for i := 0; i < len(ss); i++ {
		intVal, err := strconv.Atoi(ss[i])
		if err != nil {
			return nil, fmt.Errorf("could not parse response body: %v", err)
		}
		ints = append(ints, intVal)
	}
	return ints, nil
}
