package main

import (
	"log"
	"net/http"
)

// restriction definition matches database table
// xml and json reflection tags determine how the restrictions appear the response
type Statistics struct {
	Except30 []int `json:"exceptionsLast30Days"`
}

// FetchRestrictions performs a query on restrictions and returns the results in a RestrictionList.
func FetchStatistics() (s Statistics, err error) {
	s.Except30 = []int{0, 0, 0, 1, 2, 3, 2, 0, 1, 0, 2, 3, 2, 3, 2, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}
	return
}

// handleApiStatistics is the endpoint for statistics restful api
// accepts:
//  /api/statistics
func handleApiStatistics(w http.ResponseWriter, r *http.Request) {
	// Fetch Statistics
	if s, err := FetchStatistics(); err != nil {
		log.Println(err)
	} else if err = jsonApi(w, r, s, false); err != nil {
		log.Println(err)
	}
}
