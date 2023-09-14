package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	. "stddev-service/numbers"
	. "stddev-service/stddev"
	"strconv"
)

var numGenerator Generator = NewRandomOrgNumbersGenerator()
var calculator Calculator = &StdDevsWithSumCalculator{}

func GetRandomMean(w http.ResponseWriter, r *http.Request) {
	requests, err := getRequestsParam(r)
	if err != nil {
		writeErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	length, err := getLengthParam(r)
	if err != nil {
		writeErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	intSlices, err := numGenerator.GetIntSlices(requests, length)
	if err != nil {
		writeDefaultErrorResponse(w, err)
		return
	}

	stddevs, err := calculator.CalculateStdDevs(intSlices)
	if err != nil {
		writeDefaultErrorResponse(w, err)
		return
	}

	writeResponse(w, stddevs)
}

func getRequestsParam(r *http.Request) (int, error) {
	isPresent := r.URL.Query().Has("requests")
	if !isPresent {
		return 0, fmt.Errorf("'requests' query param required")
	}
	requestsStr := r.URL.Query().Get("requests")
	requestsInt, err := strconv.Atoi(requestsStr)
	if err != nil {
		return 0, fmt.Errorf("invalid 'requests' query param")
	}
	if requestsInt < 1 || requestsInt > 10 {
		return 0, fmt.Errorf("'requests' param value should be between 1 and 10")
	}
	return requestsInt, nil
}

func getLengthParam(r *http.Request) (int, error) {
	isPresent := r.URL.Query().Has("length")
	if !isPresent {
		return 0, fmt.Errorf("'length' query param required")
	}
	lengthStr := r.URL.Query().Get("length")
	lengthInt, err := strconv.Atoi(lengthStr)
	if err != nil {
		return 0, fmt.Errorf("invalid 'length' query param")
	}
	if lengthInt < 1 || lengthInt > 10 {
		return 0, fmt.Errorf("'length' param value should be between 1 and 10")
	}
	return lengthInt, nil
}

func writeResponse(w http.ResponseWriter, stddevs []StdDev) {
	jsonData, err := json.Marshal(stddevs)
	if err != nil {
		writeDefaultErrorResponse(w, nil)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func writeErrorResponse(w http.ResponseWriter, err error, statusCode int) {
	log.Println(err.Error())
	resp := make(map[string]string)
	resp["message"] = err.Error()
	jsonData, err := json.Marshal(resp)
	if err != nil {
		writeDefaultErrorResponse(w, nil)
		return
	}
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func writeDefaultErrorResponse(w http.ResponseWriter, err error) {
	log.Println(err.Error())
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, `{"message": "internal error occurred"}`)
}
