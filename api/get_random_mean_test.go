package api

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"stddev-service/tests"
	"testing"
)

var testCases = []struct {
	name       string
	urlParams  string
	statusCode int
	body       string
}{
	{
		"should return 200",
		"requests=2&length=5",
		200,
		`[{"stddev":2,"data":[0,1,2,3,6]},{"stddev":2,"data":[-1,0,0,2,5]},{"stddev":2,"data":[-1,0,0,0,1,2,2,3,5,6]}]`,
	},
	{
		"no requests param",
		"length=5",
		400,
		`{"message":"'requests' query param required"}`,
	},
	{
		"no length param",
		"requests=2",
		400,
		`{"message":"'length' query param required"}`,
	},
	{
		"invalid requests param",
		"requests=abc&length=5",
		400,
		`{"message":"invalid 'requests' query param"}`,
	},
	{
		"invalid length param",
		"requests=2&length=abc",
		400,
		`{"message":"invalid 'length' query param"}`,
	},
}

func TestGetRandomMean(t *testing.T) {
	numGenerator = &tests.MockedNumGenerator{}
	calculator = &tests.MockedCalculator{}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// given
			req, _ := http.NewRequest("GET", "/random/mean?"+testCase.urlParams, nil)
			rr := httptest.NewRecorder()
			// when
			GetRandomMean(rr, req)
			// then
			assert.Equal(t, rr.Code, testCase.statusCode)
			assert.Equal(t, rr.Header().Get("Content-Type"), "application/json")
			assert.Equal(t, rr.Body.String(), testCase.body)
		})
	}
}

func TestOnGeneratorError(t *testing.T) {
	// given
	numGenerator = &tests.MockedErrorNumGenerator{}
	calculator = &tests.MockedCalculator{}
	req, _ := http.NewRequest("GET", "/random/mean?requests=2&length=5", nil)
	rr := httptest.NewRecorder()
	// when
	GetRandomMean(rr, req)
	// then
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
	assert.Equal(t, `{"message": "internal error occurred"}`, rr.Body.String())
}

func TestOnCalculatorError(t *testing.T) {
	// given
	numGenerator = &tests.MockedNumGenerator{}
	calculator = &tests.MockedErrorCalculator{}
	req, _ := http.NewRequest("GET", "/random/mean?requests=2&length=5", nil)
	rr := httptest.NewRecorder()
	// when
	GetRandomMean(rr, req)
	// then
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
	assert.Equal(t, `{"message": "internal error occurred"}`, rr.Body.String())
}
