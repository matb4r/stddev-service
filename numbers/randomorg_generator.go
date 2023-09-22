package numbers

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

type RandomOrgNumbersGenerator struct {
	url      string
	min, max int
	timeout  time.Duration
}

func NewRandomOrgNumbersGenerator() RandomOrgNumbersGenerator {
	return RandomOrgNumbersGenerator{
		url:     "https://www.random.org/integers/",
		min:     1,
		max:     100,
		timeout: 3000 * time.Millisecond,
	}
}

type FetchRandomNumbersResult struct {
	numbers []int
	err     error
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
