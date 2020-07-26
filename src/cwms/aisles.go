package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Wms is Warehouse Management System inventory database record structure that matches the fields in the v_inventory view
// 	xml reflection tags are included for xml marshalling
type Wms struct {
	Id          int    `xml:"id,attr" json:"id"`
	StartTime   string `xml:"time>start" json:"startTime"`
	StopTime    string `xml:"time>stop" json:"stopTime"`
	SKU         string `xml:"item>SKU" json:"sku"`
	Discrepancy string `xml:"item>Discrepancy,omitempty" json:"discrepancy"`
	Aisle       string `xml:"position>Aisle" json:"aisle"`
	Shelf       string `xml:"position>Shelf" json:"shelf"`
	Slot        string `xml:"position>Slot" json:"slot"`
}

// WmsList is a slice of Wms
type WmsList []Wms

// toSlice converts a WmsList to a [][]string for use when generating csv output
//	Hardcoded alert!! toSlice will need to be updated if v_inventory is refactored.
//	There is a more generic approach: the rows.Columns method can be used to get the column names from the query.
//	Of course, the column names would need to be carried in page controls or a global variable perhaps when the query
//	is performed.
func (wl WmsList) toSlice() (s [][]string) {
	// Prepend the column headers
	s = append(s, []string{"Id", "Start Time", "Stop Time", "SKU", "Aisle", "Shelf", "Slot", "Discrepancy"})
	for _, v := range wl {
		s = append(s, []string{strconv.Itoa(v.Id), v.StartTime, v.StopTime, v.SKU, v.Aisle, v.Shelf, v.Slot, v.Discrepancy})
	}
	return
}

// FetchInventory performs a query on v_inventory and returns the results in a WmsList.
func FetchInventory(af AisleFilter) (wl WmsList, err error) {
	// Execute database query
	var rows *sql.Rows
	rows, err = db.Query(af.toSqlStmt())

	if err != nil {
		return
	}
	defer rows.Close()

	// Process database query results
	var record Wms
	for rows.Next() {
		err = rows.Scan(&record.Id, &record.StartTime, &record.StopTime, &record.SKU, &record.Aisle, &record.Shelf, &record.Slot, &record.Discrepancy)
		if err != nil {
			return
		}
		wl = append(wl, record)
	}
	return
}

// fetchAisles performs a query on v_inventory and returns the results in a aisleList
func fetchAisles(filter string) (aisleList []string, err error) {
	// Execute database query
	var rows *sql.Rows
	rows, err = db.Query(`select distinct aisle from v_inventory order by aisle`)
	if err != nil {
		return
	}
	defer rows.Close()

	// Process database query results
	var aisle string
	for rows.Next() {
		err = rows.Scan(&aisle)
		if err != nil {
			return
		}
		aisleList = append(aisleList, aisle)
	}
	return
}

// AisleFilter holds Aisle filter information
// Aisle and Discrepancy filters a cumulative
type AisleFilter struct {
	Aisle       string // Filter on Aisle
	Discrepancy string // Filter on Discrepancies
}

// toSqlStmt generates a sql statement
func (af AisleFilter) toSqlStmt() (sqlstmt string) {
	var sel, order string
	var where []string
	sel = `select inventoryId, startTime, stopTime, sku, aisle, block, slot, discrepancy from v_inventory `
	if af.Aisle != "" {
		where = append(where, fmt.Sprintf(`aisle ='%s'`, af.Aisle))
	}
	if af.Discrepancy == "all" {
		where = append(where, `discrepancy !="" `)
	} else if af.Discrepancy != "" {
		where = append(where, fmt.Sprintf(`discrepancy ='%s'`, af.Discrepancy))
	}
	order = `order by aisle, block, slot`
	if len(where) > 0 {
		sqlstmt = fmt.Sprintf("%s where %s %s", sel, strings.Join(where, " and "), order)
	} else {
		sqlstmt = fmt.Sprintf("%s %s", sel, order)
	}
	return
}

func handleApiAisles(w http.ResponseWriter, r *http.Request) {
	// Fetch inventory based on page controls
	var af AisleFilter

	// Get segment list from request, set aisle filter if the last segment is a specific aisle
	sl := strings.Split(r.URL.Path, "/")
	if len(sl) > 0 {
		ls := sl[len(sl)-1]
		if ls != "" {
			af.Aisle = ls
		}
	}

	if af.Aisle == "" {
		asl, err := fetchAisleStats()
		if err != nil {
			log.Println(err)
		}
		// Send filtered inventory in json response
		if err = jsonApi(w, asl); err != nil {
			log.Println(err)
		}
	} else {
		// Fetch inventory filtered by aisle filter
		wl, err := FetchInventory(af)
		if err != nil {
			log.Println(err)
		}
		// Send filtered inventory in json response
		if err = jsonApi(w, wl); err != nil {
			log.Println(err)
		}
	}
}

func handleApiDiscrepancies(w http.ResponseWriter, r *http.Request) {
	// Fetch inventory based on page controls
	var af AisleFilter

	af.Discrepancy = "all"

	// Get segment list from request, set discrepancy filter if the last segment is a specific aisle
	sl := strings.Split(r.URL.Path, "/")
	if len(sl) > 0 {
		ls := sl[len(sl)-1]
		log.Println(ls)
		if ls != "" {
			af.Discrepancy = ls
		}
	}

	// Fetch inventory filtered by aisle filter
	wl, err := FetchInventory(af)
	if err != nil {
		log.Println(err)
	}

	// Send filter inventory in json response
	if err = jsonApi(w, wl); err != nil {
		log.Println(err)
	}
}

type aisleStats struct {
	Id              string `db:"aisle" json:"id"`
	NumberOccupied  int    `db:"numberOccupied" json:"numberOccupied"`
	NumberEmpty     int    `db:"numberEmpty" json:"numberEmpty"`
	NumberException int    `db:"numberException" json:"numberException"`
	NumberUnscanned int    `db:"numberUnscanned" json:"numberUnscanned"`
	LastScanned     string `db:"lastScanned" json:"lastScanned"`
}

type aisleStatsList []aisleStats

func fetchAisleStats() (asl aisleStatsList, err error) {
	as := aisleStats{Id: "1a", NumberOccupied: 10, NumberEmpty: 5, NumberException: 6, NumberUnscanned: 1, LastScanned: "2020-04-04T19:22:45.004Z"}
	al := []string{"1a", "1b", "1c", "2a", "2b", "2c", "3a", "3b", "3c", "4a"}
	for i := 0; i < 10; i++ {
		as.Id = al[i]
		as.NumberOccupied += i
		as.NumberEmpty += i
		as.NumberException += i
		asl = append(asl, as)
	}
	return
}
