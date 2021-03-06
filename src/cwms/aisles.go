package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"
	"net/http"
	"strings"
	"encoding/json"
)

// Wms is Warehouse Management System inventory database record structure that matches the fields in the v_inventory view
// 	xml reflection tags are included for xml marshalling
type Wms struct {
	Id          int    `xml:"id,attr" json:"id" csv:"-"`
	StartTime   time.Time `xml:"time>start" json:"startTime"  csv:"-"`
	StopTime    time.Time `xml:"time>stop" json:"stopTime"  csv:"-"`
	SKU         NullString `xml:"item>SKU" json:"sku"  csv:"sku"`
	Discrepancy NullString `xml:"item>Discrepancy,omitempty" json:"discrepancy"  csv:"-"`
	Aisle       string `xml:"position>Aisle" json:"aisle"  csv:"aisle"`
	Block       string `xml:"position>Block" json:"block"  csv:"block"`
	Slot        string `xml:"position>Slot" json:"slot"  csv:"slot"`
	Shelf       string `xml:"position>Shelf" json:"shelf"  csv:"-"`
	DisplayName       string `xml:"position>DisplayName" json:"displayname"  csv:"display_name"`
	Image       NullString `xml:"position>Image" json:"image"  csv:"image_url"`
}

// NullString is an alias for sql.NullString data type
type NullString struct {
	sql.NullString
}


// MarshalJSON for NullString
func (ns *NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("\"\""), nil //TODO this is dumb, should be []byte("null")
	}
	return json.Marshal(ns.String)
}

// UnmarshalJSON for NullString
func (ns *NullString) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &ns.String)
	ns.Valid = (err == nil)
	return err
}


// MarshalCSV for NullString
func (ns *NullString) MarshalCSV() ([]byte, error) {
	if !ns.Valid {
		return []byte(""), nil //TODO this is dumb, should be []byte("null")
	}
	return []byte(ns.String), nil
}

// WmsList is a slice of Wms
type WmsList []Wms

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
		err = rows.Scan(&record.Id, &record.StartTime, &record.StopTime, &record.SKU, &record.Aisle, &record.Block, &record.Slot, &record.Shelf, &record.DisplayName, &record.Discrepancy, &record.Image)
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
	sel = `select inventoryId, startTime, stopTime, sku, aisle, block, slot, shelf, displayName, discrepancy, imageUrl from v_inventory `
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
		if err = jsonApi(w, r, asl, false); err != nil {
			log.Println(err)
		}
	} else {
		// Fetch inventory filtered by aisle filter
		wl, err := FetchInventory(af)
		if err != nil {
			log.Println(err)
		}
		// Send filtered inventory in json response
		if err = jsonApi(w, r, wl, true); err != nil {
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
	if err = jsonApi(w, r, wl, false); err != nil {
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
	// Execute database query
	var rows *sql.Rows
	if rows, err = db.Query("select distinct aisle, numberException, numberEmpty, numberOccupied, numberUnscanned, lastScanned from v_aisleStats"); err != nil {
		return
	}
	defer rows.Close()

	var as aisleStats
	// Process query results
	for rows.Next() {
		// Load query results into interface list via the pointers
		if err = rows.Scan(StructForScan(&as)...); err != nil {
			return
		}

		// append query results to flight list
		asl = append(asl, as)
	}
	return
}