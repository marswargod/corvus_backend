package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)
type basicFlight struct {
	FlightId  int    `json:"id" db:"flightId"`
	Time      string `json:"time" db:"time"`
}


type flight struct {
	FlightId  int    `json:"id" db:"flightId"`
	Time      string `json:"time" db:"time"`
	Sku       string `json:"sku" db:"sku"`
	Occupancy string `json:"occupancy" db:"occupancy"`
	Aisle     string `json:"aisle" db:"aisle"`
	Shelf     string `json:"shelf" db:"shelf"`
	Slot      string `json:"slot" db:"slot"`
}

func (f basicFlight) toFieldList() (fl []string) {
	rt := reflect.TypeOf(f)
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		fl = append(fl, field.Tag.Get("db"))
	}
	return
}

func (f flight) toFieldList() (fl []string) {
	rt := reflect.TypeOf(f)
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		fl = append(fl, field.Tag.Get("db"))
	}
	return
}

func (f basicFlight) toPlaceholderList() (phl []string) {
	rt := reflect.TypeOf(f)
	for i := 0; i < rt.NumField(); i++ {
		phl = append(phl, "?")
	}
	return
}

func (f flight) toPlaceholderList() (phl []string) {
	rt := reflect.TypeOf(f)
	for i := 0; i < rt.NumField(); i++ {
		phl = append(phl, "?")
	}
	return
}

func (f basicFlight) toInterfaceList() (il []interface{}) {
	v := reflect.ValueOf(f)
	for i := 0; i < v.NumField(); i++ {
		il = append(il, v.Field(i).Interface())
	}
	return
}

func (f flight) toInterfaceList() (il []interface{}) {
	v := reflect.ValueOf(f)
	for i := 0; i < v.NumField(); i++ {
		il = append(il, v.Field(i).Interface())
	}
	return
}

func convertInterfaceListToBasicFlight(il []interface{}) (f basicFlight) {
	v := reflect.ValueOf(&f).Elem()
	for i := 0; i < v.NumField(); i++ {
		switch il[i].(type) {
		case int64:
			v.Field(i).SetInt(il[i].(int64))
		case string:
			v.Field(i).SetString(il[i].(string))
		default:
			log.Println("default")
		}
	}
	return
}

func convertInterfaceListToFlight(il []interface{}) (f flight) {
	v := reflect.ValueOf(&f).Elem()
	for i := 0; i < v.NumField(); i++ {
		switch il[i].(type) {
		case int64:
			v.Field(i).SetInt(il[i].(int64))
		case string:
			v.Field(i).SetString(il[i].(string))
		default:
			log.Println("default")
		}
	}
	return
}

type basicFlightList []basicFlight

type flightList []flight

type flightFilter struct {
	FlightId int
	Sku      string
	Before   string
	After    string
	Aisle    string
	Limit    int
	Offset   int
	Sort     string
	Order_by string
}

func (ff flightFilter) toSqlSelect() (sqlstmt string) {

	var sel, ord, limit string

	// Format select statement using field list
	var f flight
	sel = fmt.Sprintf("select %s from v_flightList", strings.Join(f.toFieldList(), ", "))

	// Accumulate where clauses
	var where []string
	if ff.FlightId != 0 {
		where = append(where, fmt.Sprintf(`flightId='%v'`, ff.FlightId))
	} else {
		if ff.Sku != "" {
			where = append(where, fmt.Sprintf(`sku LIKE '%%%s%%'`, ff.Sku))
		}
		if ff.Before != "" {
			where = append(where, fmt.Sprintf(`before LIKE '%%%s%%'`, ff.Sku))
		}
		if ff.After != "" {
			where = append(where, fmt.Sprintf(`after LIKE '%%%s%%'`, ff.Sku))
		}
		if ff.Aisle != "" {
			where = append(where, fmt.Sprintf(`aisle LIKE '%%%s%%'`, ff.Sku))
		}
	}

	// Format order by
	if ff.Sort != "" {
		ord = fmt.Sprintf(" order by %s %s", ff.Sort, ff.Order_by)
	}

	// Format limit and offset
	var lim []string
	if ff.Limit != 0 {
		lim = append(lim, fmt.Sprintf("%v", ff.Limit))
		if ff.Offset != 0 {
			lim = append(lim, fmt.Sprintf("%v", ff.Offset))
		}
		limit = fmt.Sprintf(" LIMIT %s", strings.Join(lim, ", "))
	}

	// Format sql statement
	// Start with select clause
	sqlstmt = sel

	// Append where clause
	if len(where) > 0 {
		sqlstmt += " where "
		sqlstmt += strings.Join(where, " AND ")
	}

	// Append order by clause
	if ord != "" {
		sqlstmt += ord
	}

	// Append limit and offset
	if limit != "" {
		sqlstmt += limit
	}

	return
}

func FetchBasicFlights() (bfl basicFlightList, err error) {
		// Execute database query
		var rows *sql.Rows
		if rows, err = db.Query("select distinct flightId, time from v_flightList"); err != nil {
			return
		}
		defer rows.Close()
	
		// get number of fields for a flight from reflect
		var f basicFlight
		numCols := reflect.TypeOf(f).NumField()
	
		// Process query results
		for rows.Next() {
	
			// Create interface list and set pointers to interface list
			cols := make([]interface{}, numCols)
			ptrs := make([]interface{}, numCols)
			for i := 0; i < numCols; i++ {
				ptrs[i] = &cols[i]
			}
	
			// Load query results into interface list via the pointers
			if err = rows.Scan(ptrs...); err != nil {
				return
			}
	
			// append query results to flight list
			bfl = append(bfl, convertInterfaceListToBasicFlight(cols))
		}
		return
}

// FetchInventory performs a query on v_inventory and returns the results in a WmsList.
func FetchFlights(ff flightFilter) (fl flightList, err error) {
	// Execute database query
	var rows *sql.Rows
	if rows, err = db.Query(ff.toSqlSelect()); err != nil {
		return
	}
	defer rows.Close()

	// get number of fields for a flight from reflect
	var f flight
	numCols := reflect.TypeOf(f).NumField()

	// Process query results
	for rows.Next() {

		// Create interface list and set pointers to interface list
		cols := make([]interface{}, numCols)
		ptrs := make([]interface{}, numCols)
		for i := 0; i < numCols; i++ {
			ptrs[i] = &cols[i]
		}

		// Load query results into interface list via the pointers
		if err = rows.Scan(ptrs...); err != nil {
			return
		}

		// append query results to flight list
		fl = append(fl, convertInterfaceListToFlight(cols))
	}
	return
}

// createFilter implements a simple api to read a wms list in a json format
func createFilter(r *http.Request) (ff flightFilter, err error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}
	err = json.Unmarshal(body, &ff)
	if err != nil {
		log.Println(err)
	}
	return
}

func handleApiFlights(w http.ResponseWriter, r *http.Request) {
	// Fetch inventory based on page controls
	var ff flightFilter

	// Get segment list from request, set flight filter if the last segment is a flight aisle
	// either process a flight number or a flight filter
	sl := strings.Split(r.URL.Path, "/")
	if len(sl) > 0 {
		ls := sl[len(sl)-1]
		if ls != "" {
			ff.FlightId, _ = strconv.Atoi(ls)
		} else {
			// myff, _ := createFilter(r)
			// log.Println(myff)
			// switch r.Method {
			// case http.MethodGet:
			// 	log.Println("Serve the resource.")
			// case http.MethodPost:
			// 	log.Println("Create a new record.")
			// }
		}
	}

	if ff.FlightId == 0 {
	// Fetch inventory filtered by aisle filter
	wl, err := FetchBasicFlights()
	if err != nil {
		log.Println(err)
	}

	// Send filter inventory in json response
	if err := jsonApi(w, r, wl); err != nil {
		log.Println(err)
	}
	} else {

	// Fetch inventory filtered by aisle filter
	wl, err := FetchFlights(ff)
	if err != nil {
		log.Println(err)
	}

	// Send filter inventory in json response
	if err := jsonApi(w, r, wl); err != nil {
		log.Println(err)
	}
}
}
