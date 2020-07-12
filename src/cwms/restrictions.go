package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// restriction definition matches database table
// xml and json reflection tags determine how the restrictions appear the response
type Restriction struct {
	Id          int    `xml:"id,attr" json:"id"`
	Name        string `xml:"name,attr" json:"name"`
	StartTime   string `xml:"time>start" json:"start time"`
	StopTime    string `xml:"time>stop" json:"stop time"`
	Periodicity string `xml:"periodicity" json:"periodicity"`
	Region      int    `xml:"region" json:"region"`
}
type RestrictionList []Restriction

// RestrictionFilter holds Restriction filter information
// Restriction and tbd filters a cumulative
type RestrictionFilter struct {
	Name string // Filter on name
}

// toSqlStmt generates a sql statement based on the current set of page controls
func (rf RestrictionFilter) toSqlStmt() (sqlstmt string) {
	var sel, order string
	var where []string
	sel = `select restrictionId, name, startTime, stopTime, periodicity, regionId from restrictions `
	if rf.Name != "" {
		where = append(where, fmt.Sprintf(`name ='%s'`, rf.Name))
	}
	order = `order by regionId`
	if len(where) > 0 {
		sqlstmt = fmt.Sprintf("%s where %s %s", sel, strings.Join(where, " and "), order)
	} else {
		sqlstmt = fmt.Sprintf("%s %s", sel, order)
	}
	return
}

// FetchRestrictions performs a query on restrictions and returns the results in a RestrictionList.
func FetchRestrictions(rf RestrictionFilter) (rl RestrictionList, err error) {
	// Execute database query
	var rows *sql.Rows
	rows, err = db.Query(rf.toSqlStmt())

	if err != nil {
		return
	}
	defer rows.Close()

	// Process database query results
	var record Restriction
	for rows.Next() {
		err = rows.Scan(&record.Id, &record.Name, &record.StartTime, &record.StopTime, &record.Periodicity, &record.Region)
		if err != nil {
			return
		}
		rl = append(rl, record)
	}
	return
}

// handleApiRestrictions is the endpoint for restrictions restful api
// accepts:
//  /restrictions
//	/restrictions/:name
// Sets restrictions filter based on name and writes a json response with
// a list of restrictions.
func handleApiRestrictions(w http.ResponseWriter, r *http.Request) {
	// Fetch restrictions based on filter
	var rf RestrictionFilter

	// Get segment list from request, set discrepancy filter if the last segment is a specific aisle
	sl := strings.Split(r.URL.Path, "/")
	if len(sl) > 0 {
		ls := sl[len(sl)-1]
		if ls != "" {
			rf.Name = ls
		}
	}

	// Fetch restrictions filtered by restriction filter
	rl, err := FetchRestrictions(rf)
	if err != nil {
		log.Println(err)
	}

	// Send filtered restriction list in json response
	if err = jsonApi(w, rl); err != nil {
		log.Println(err)
	}
}
