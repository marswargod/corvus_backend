package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"
)

type flight struct {
	FlightId  int    `json:"id" db:"flightId"`
	Time      string `json:"time" db:"time"`
	Sku       string `json:"sku" db:"sku"`
	Occupancy string `json:"occupancy" db:"occupancy"`
	Aisle     string `json:"aisle" db:"aisle"`
	Shelf     string `json:"shelf" db:"shelf"`
	Slot      string `json:"slot" db:"slot"`
}

func (f flight) toFieldList() (fl []string) {
	rt := reflect.TypeOf(f)
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		fl = append(fl, field.Tag.Get("db"))
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

func (f flight) toInterfaceList() (il []interface{}) {
	v := reflect.ValueOf(f)
	for i := 0; i < v.NumField(); i++ {
		il = append(il, v.Field(i).Interface())
	}
	return
}

func convertInterfaceListToFlight(il []interface{}) (f flight) {
	v := reflect.ValueOf(&f).Elem()
	for i := 1; i < v.NumField(); i++ {
		switch il[i].(type) {
		case int64:
			v.Field(i).SetInt(il[i].(int64))
		case string:
			v.Field(i).SetString(il[i].(string))
		}
	}
	return
}

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
			where = append(where, fmt.Sprintf(`sku LIKE '%%%s%%'`, ff.Sku))
		}
		if ff.Aisle != "" {
			where = append(where, fmt.Sprintf(`sku LIKE '%%%s%%'`, ff.Sku))
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

func handleApiFlights(w http.ResponseWriter, r *http.Request) {
	// Fetch inventory based on page controls
	var ff flightFilter

	// Get segment list from request, set aisle filter if the last segment is a specific aisle
	sl := strings.Split(r.URL.Path, "/")
	if len(sl) > 0 {
		ls := sl[len(sl)-1]
		if ls != "" {
			ff.Sku = ls
		}
	}

	// Fetch inventory filtered by aisle filter
	wl, err := FetchFlights(ff)
	if err != nil {
		log.Println(err)
	}

	// Send filter inventory in json response
	if err = jsonApi(w, wl); err != nil {
		log.Println(err)
	}
}
