// 	Corvus Warehouse Management System (cwms)
// 	Provides a suite of displays that allows a user to:
// 	• View the entire inventory
// 	• Navigate by aisle
// 	• Filter by discrepancies
//	• Compare Drone Inventory to Warehouse Inventory
package main

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3" // Reference for installed sqlite3 driver
)

// Min implements an integer min function
func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// Max implements an integer max function
func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

// db is a global declaration for the database
var db *sql.DB

// type LocalFileSystem implements an http fileserver handler to perform local file system functions.
type LocalFileSystem struct {
	fs http.FileSystem
}

// This Open opens a file in the local file system
func (l LocalFileSystem) Open(name string) (f http.File, err error) {
	f, err = l.fs.Open(name)
	if err != nil {
		log.Println("open", err)
		return nil, err
	}
	var fi os.FileInfo
	fi, err = f.Stat()
	if fi.IsDir() {
		log.Println("file info", err)
		return nil, os.ErrNotExist
	}
	return
}

// executeTemplate parses and executes the specified html template file
func executeTemplate(f string, tm map[string]interface{}, w http.ResponseWriter) (err error) {
	// parse html template file
	t, err := template.ParseFiles(f)
	if err != nil {
		return
	}

	// execute html template
	err = t.Execute(w, tm)
	if err != nil {
		return
	}
	return
}

// main
// 	• opens the database
// 	• sets up the mutex
// 	• sets up the http handlers
// 	• listens and serves
func main() {
	// Setup logger
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var err error

	// Open global database
	db, err = sql.Open("sqlite3", `.\wms3.db`)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	// Setup servemux to serve http handler routines
	mux := http.NewServeMux()

	// Setup file server handler
	files := http.FileServer(LocalFileSystem{http.Dir(`.\static\`)})
	mux.Handle("/static/", http.StripPrefix("/static/", files))

	// Setup http handlers
	mux.HandleFunc("/", imw(handleDashboard))
	mux.HandleFunc("/dashboard/", imw(handleDashboard))
	mux.HandleFunc("/inventory/", imw(handleInventory))
	mux.HandleFunc("/hybrid/", imw(handleHybrid))
	mux.HandleFunc("/schedule/", mmw(handleSchedule))
	mux.HandleFunc("/export/csv/", imw(handleExportInventoryCsv))
	mux.HandleFunc("/export/json/", imw(handleExportInventoryJson))
	mux.HandleFunc("/export/xml/", imw(handleExportInventoryXml))
	// restful api handlers
	mux.Handle("/api/", http.NotFoundHandler())
	mux.HandleFunc("/api/json/", imw(handleApiInventoryJson))
	mux.HandleFunc("/api/aisles/", handleApiAisles)
	mux.HandleFunc("/api/discrepancy/", handleApiDiscrepancies)
	mux.HandleFunc("/api/restrictions/", handleApiRestrictions)
	mux.HandleFunc("/api/flights/", handleApiFlights)
	mux.HandleFunc("/api/statistics/", handleApiStatistics)

	// Listen and serve mux
	http.ListenAndServe(":8080", mux)
}

// jsonApi implements a simple restful api to export data in a json format
func jsonApi(w http.ResponseWriter, data interface{}) (err error) {
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(data); err != nil {
		log.Println(err)
	}
	return
}
