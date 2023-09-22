package numbers

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestShouldMakeIntSlicesFromTwoRequests(t *testing.T) {
	// given
	requestsCount := 0
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		switch requestsCount {
		case 0:
			w.Write([]byte("3\n4\n99"))
		case 1:
			w.Write([]byte("6\n20\n11"))
		}
		requestsCount++
	}))
	defer testServer.Close()
	generator := BufferedChannelRandomOrgNumbersGenerator{RandomOrgNumbersGenerator{url: testServer.URL, min: 1, max: 100, timeout: time.Second}}

	// when
	res, err := generator.GetIntSlices(2, 3)

	// then
	assert.ElementsMatch(t, [][]int{{3, 4, 99}, {6, 20, 11}}, res)
	assert.NoError(t, err)
	assert.Equal(t, 2, requestsCount)
}

func TestShouldReturnErrorWhenTimeoutOnOneOfRequests(t *testing.T) {
	//given
	requestsCount := 0
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		switch requestsCount {
		case 0:
			requestsCount++
			w.Write([]byte("3\n4\n99"))
		case 1:
			requestsCount++
			time.Sleep(110 * time.Millisecond)
			w.Write([]byte("6\n20\n11"))
		}
	}))
	defer testServer.Close()
	generator := &BufferedChannelRandomOrgNumbersGenerator{RandomOrgNumbersGenerator{url: testServer.URL, min: 1, max: 100, timeout: 100 * time.Millisecond}}

	// when
	res, err := generator.GetIntSlices(2, 3)

	// then
	assert.Nil(t, res)
	assert.ErrorContains(t, err, "context deadline exceeded")
	assert.Equal(t, 2, requestsCount)
}

func TestShouldReturnErrorWhenNot200OnOneOfRequests(t *testing.T) {
	// given
	requestsCount := 0
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch requestsCount {
		case 0:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("3\n4\n99"))
		case 1:
			w.WriteHeader(http.StatusBadRequest)
		}
		requestsCount++
	}))
	defer testServer.Close()
	generator := &BufferedChannelRandomOrgNumbersGenerator{RandomOrgNumbersGenerator{url: testServer.URL, min: 1, max: 100, timeout: time.Second}}

	// when
	res, err := generator.GetIntSlices(2, 3)

	// then
	assert.Nil(t, res)
	assert.ErrorContains(t, err, "got status code 400")
	assert.Equal(t, 2, requestsCount)
}
