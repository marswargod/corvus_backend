package main

import (
	"database/sql"
	"log"
	"net/http"
)

type Mms struct {
	Entry     int    `json:"entry"`
	Region    string `json:"region"`
	Frequency string `json:"frequency"`
	Aisle     string `json:"aisle"`
	Shelf     string `json:"shelf"`
	Slot      string `json:"slot"`
}

func (m Mms) toSummary() (ms MmsSummary) {
	ms.Region = m.Region
	ms.Frequency = m.Frequency
	ms.Entry = m.Entry
	ms.StartAisle = m.Aisle
	ms.EndAisle = m.Aisle
	return
}

type MmsList []Mms

func (ml MmsList) toSummaryList() (msl MmsSummaryList) {
	for _, m := range ml {
		if m.Region == "" { // skip no name entries
			continue
		}

		if (len(msl) == 0) || (msl[len(msl)-1].Region != m.Region) { // append new entries
			msl = append(msl, m.toSummary())
		} else {
			if m.Aisle != "" {
				msl[len(msl)-1].EndAisle = m.Aisle // extend existing entries
			}
		}
	}
	return
}

type MmsSummary struct {
	Region     string
	Frequency  string
	Entry      int
	StartAisle string
	EndAisle   string
}

type MmsSummaryList []MmsSummary

// mmwHandler is a middleware handler function signature used by the mmw middleware
type mmwHandler func(tm map[string]interface{}, ml MmsList, w http.ResponseWriter, r *http.Request)

// mmw mission middleware fetches missions, statistics, and mission controls and loads them into the template map
//	mmw does "everything":
//	• loads missions into the template map (tm)
//	• loads mission controls into the template map
func mmw(next mmwHandler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Fetch url parameters
		urlParams := r.URL.Query()

		// Create mission controls
		mc, err := missionControls(urlParams.Get("day"), urlParams.Get("scope"))
		if err != nil {
			log.Println(err)
		}

		// Fetch missions based on mission controls
		ml, err := FetchSchedule(mc)
		if err != nil {
			log.Println(err)
		}

		// Create template map
		tm := make(map[string]interface{})

		// Load mission controls into template map
		tm["MissionControls"] = mc

		// Load missions into template map
		tm["Missions"] = ml.toSummaryList()

		next(tm, ml, w, r)
	})
}

// FetchSchedule performs a query on v_inventory and returns the results in a WmsList.
func FetchSchedule(mc MissionControls) (ml MmsList, err error) {
	// Execute database query
	var rows *sql.Rows
	rows, err = db.Query(mc.toSqlStmt())

	if err != nil {
		return
	}
	defer rows.Close()

	// Process database query results
	var record Mms
	for rows.Next() {
		err = rows.Scan(&record.Entry,
			&record.Region,
			&record.Frequency,
			&record.Aisle,
			&record.Shelf,
			&record.Slot)
		if err != nil {
			return
		}
		ml = append(ml, record)
	}
	return
}

// fetchDays performs a query on v_schedule and returns the results in a dayList
func fetchDays(filter string) (dayList []string, err error) {
	// Execute database query
	var rows *sql.Rows
	rows, err = db.Query(`select distinct region from v_schedule order by entry`)
	if err != nil {
		return
	}
	defer rows.Close()

	// Process database query results
	var day string
	for rows.Next() {
		err = rows.Scan(&day)
		if err != nil {
			return
		}
		dayList = append(dayList, day)
	}
	return
}
