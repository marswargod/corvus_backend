// Corvus Warehouse Management System (cwms)
// 	Provides a suite of displays that allows a user to:
// 	• View the entire inventory
// 	• Navigate by aisle
// 	• Filter by discrepancies
//	• Compare Drone Inventory to Warehouse Inventory
package main

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3" // Reference for installed sqlite3 driver
)

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

// Wms is Warehouse Management System inventory database record structure that matches the fields in the v_inventory view
// 	xml reflection tags are included for xml marshalling
type Wms struct {
	Id          int    `xml:"id,attr"`
	StartTime   string `xml:"time>start"`
	StopTime    string `xml:"time>stop"`
	SKU         string `xml:"item>SKU"`
	Discrepancy string `xml:"item>Discrepancy,omitempty"`
	Aisle       string `xml:"position>Aisle"`
	Shelf       string `xml:"position>Shelf"`
	Slot        string `xml:"position>Slot"`
}

// WmsList is a slice of Wms
type WmsList []Wms

// toSlice converts a WmsList to a [][]string
func (wl WmsList) toSlice() (s [][]string) {
	s = append(s, []string{"Id", "Start Time", "Stop Time", "SKU", "Aisle", "Shelf", "Slot", "Discrepancy"})
	for _, v := range wl {
		s = append(s, []string{strconv.Itoa(v.Id), v.StartTime, v.StopTime, v.SKU, v.Aisle, v.Shelf, v.Slot, v.Discrepancy})
	}
	return
}

// FetchInventory performs a query on v_inventory and returns the results in a WmsList.
func FetchInventory(pc PageControls) (wl WmsList, err error) {
	// Execute database query
	var rows *sql.Rows
	rows, err = db.Query(pc.toSqlStmt())

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

// Stats contains a set of statistics derived from v_inventory view
type Stats struct {
	TotalSkus   int
	SkuIssues   int
	EmptySlots  int
	FilledSlots int
	TotalSlots  int
}

// fetchStats performs various queries on v_inventory and returns the results in Stats
func fetchStats() (stats Stats, err error) {
	err = db.QueryRow(`select count(1) from v_inventory`).Scan(&stats.TotalSkus)
	if err != nil {
		return
	}

	err = db.QueryRow(`select count(1) from v_inventory where discrepancy !=""`).Scan(&stats.SkuIssues)
	if err != nil {
		return
	}

	err = db.QueryRow(`select count (1) from positions`).Scan(&stats.TotalSlots)
	if err != nil {
		return
	}

	err = db.QueryRow(`select count (1) from positions`).Scan(&stats.FilledSlots)
	if err != nil {
		return
	}

	return
}

// imwHandler is a middleware handler function signature used by the imw middleware
type imwHandler func(tm map[string]interface{}, wl WmsList, w http.ResponseWriter, r *http.Request)

// imw inventory middleware fetches inventory, statistics, and page controls and loads them into the template map
//	imw does "everything":
//	• loads inventory into the template map (tm)
//	• loads page controls into the template map
//	• loads statistics into the template map
func imw(next imwHandler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Fetch url parameters
		urlParams := r.URL.Query()

		// Create page controls
		pc, err := pageControls(urlParams.Get("aisle"), urlParams.Get("scope"))
		if err != nil {
			log.Println(err)
		}

		// Fetch inventory based on page controls
		wl, err := FetchInventory(pc)
		if err != nil {
			log.Println(err)
		}

		// Fetch inventory statistics
		stats, err := fetchStats()
		if err != nil {
			log.Println(err)
		}

		// Create template map
		tm := make(map[string]interface{})

		// Load page controls into template map
		tm["PageControls"] = pc

		// Load inventory into template map
		tm["Inventory"] = wl

		// Fetch stats into template map
		tm["Stats"] = stats

		next(tm, wl, w, r)
	})
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

// PageControls holds page navigation and page control fields
type PageControls struct {
	Curr, Next, Prev string   // Page Nav
	SingleAisle      bool     // Has a single aisle been selected?
	Selection        string   // selection choice "all" or an aisle name
	Scope            string   // filter scope: blank or "issues"
	Aisles           []string // list of aisles
}

// toSqlStmt generates a sql statement based on the current set of page controls
func (pc PageControls) toSqlStmt() (sqlstmt string) {
	var sel, order string
	var where []string
	sel = `select inventoryId, startTime, stopTime, sku, aisle, block, slot, discrepancy from v_inventory `
	if pc.SingleAisle {
		where = append(where, fmt.Sprintf(`aisle ='%s'`, pc.Curr))
	}
	if pc.Scope != "" {
		where = append(where, `discrepancy !="" `)
	}
	order = `order by aisle, block, slot`
	if len(where) > 0 {
		sqlstmt = fmt.Sprintf("%s where %s %s", sel, strings.Join(where, " and "), order)
	} else {
		sqlstmt = fmt.Sprintf("%s %s", sel, order)
	}
	return
}

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

// pageControls generates a set of page nav and page controls based on aisle and scope
func pageControls(aisle, scope string) (pc PageControls, err error) {
	var aisles []string
	aisles, err = fetchAisles(scope)
	if err != nil {
		return
	}
	if aisle == "" {
		aisle = aisles[0]
	}
	pc = pageNav(aisle, aisles)
	pc.Scope = scope
	pc.Selection = aisle
	pc.SingleAisle = aisle != "all"
	return
}

// pageNav creates a new page controls and initializes the page navigation portion of page controls based on curr aisle and aisle list
func pageNav(curr string, al []string) (pc PageControls) {
	pc = PageControls{Aisles: al, Curr: al[0], Next: al[0], Prev: al[0]}
	for c, v := range al {
		if v == curr {
			pc.Curr = al[c]
			pc.Next = al[Min(len(al)-1, c+1)]
			pc.Prev = al[Max(0, c-1)]
			break
		}
	}
	return
}

// == Main routine

// main opens the database and sets up the http server
func main() {
	// Setup logger
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var err error

	// Open global database
	db, err = sql.Open("sqlite3", `C:\Users\James\projects\test\wms3.db`)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	// Setup servemux to serve http handler routines
	mux := http.NewServeMux()

	// Setup file server handler
	files := http.FileServer(LocalFileSystem{http.Dir(`C:\Users\James\projects\test\static\`)})
	mux.Handle("/static/", http.StripPrefix("/static/", files))

	// Setup http handlers
	mux.HandleFunc("/", imw(handleDashboard))
	mux.HandleFunc("/dashboard/", imw(handleDashboard))
	mux.HandleFunc("/inventory/", imw(handleInventory))
	mux.HandleFunc("/hybrid/", imw(handleHybrid))
	mux.HandleFunc("/export/csv/", imw(handleExportInventoryCsv))
	mux.HandleFunc("/export/json/", imw(handleExportInventoryJson))
	mux.HandleFunc("/export/xml/", imw(handleExportInventoryXml))

	// Listen and serve mux
	http.ListenAndServe(":8080", mux)
}

// handleInventory provides inventory comparison functions
func handleInventory(tm map[string]interface{}, wl WmsList, w http.ResponseWriter, r *http.Request) {
	err := executeTemplate("inventory.html", tm, w)
	if err != nil {
		log.Println(err)
	}
}

// handleDashboard provides main navigation to all webpages
func handleDashboard(tm map[string]interface{}, wl WmsList, w http.ResponseWriter, r *http.Request) {
	err := executeTemplate("dashboard.html", tm, w)
	if err != nil {
		log.Println(err)
	}
}

// handleHybrid provides basic navigation features and downloads files in csv, json, or xml formats
func handleHybrid(tm map[string]interface{}, wl WmsList, w http.ResponseWriter, r *http.Request) {
	err := executeTemplate("hybrid.html", tm, w)
	if err != nil {
		log.Println(err)
	}
}

// handleExportInventoryCsv downloads the inventory to a CSV file
func handleExportInventoryCsv(tm map[string]interface{}, wl WmsList, w http.ResponseWriter, r *http.Request) {
	err := csvDownload(w, "inventory.csv", wl.toSlice())
	if err != nil {
		log.Println(err)
	}
}

// handleExportInventory downloads the inventory to a JSON file
func handleExportInventoryJson(tm map[string]interface{}, wl WmsList, w http.ResponseWriter, r *http.Request) {
	err := jsonDownload(w, "inventory_json.txt", wl)
	if err != nil {
		log.Println(err)
	}
}

// handleExportInventoryXml downloads the inventory to an XML file
func handleExportInventoryXml(tm map[string]interface{}, wl WmsList, w http.ResponseWriter, r *http.Request) {
	err := xmlDownload(w, "inventory_xml.txt", wl)
	if err != nil {
		log.Println(err)
	}
}

// csvDownload downloads the inventory to a csv file via the web browser
func csvDownload(w http.ResponseWriter, filename string, data [][]string) (err error) {
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=%s", filename))
	writer := csv.NewWriter(w)
	defer writer.Flush()

	for _, value := range data {
		if err = writer.Write(value); err != nil {
			return
		}
	}
	return
}

// jsonDownload downloads the inventory to a json file via the web browser
func jsonDownload(w http.ResponseWriter, filename string, data WmsList) (err error) {
	w.Header().Set("Content-Type", "text")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=%s", filename))
	writer := json.NewEncoder(w)

	for _, value := range data {
		if err = writer.Encode(value); err != nil {
			return
		}
	}
	return
}

// csvExport exports the inventory to a local csv file
func csvExport(filename string, data [][]string) (err error) {
	var file *os.File
	file, err = os.Create(filename)
	if err != nil {
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, value := range data {
		if err = writer.Write(value); err != nil {
			return
		}
	}
	return
}

// jsonExport exports the inventory to a local JSON file
func jsonExport(filename string, data [][]string) (err error) {
	var file *os.File

	if file, err = os.Create(filename); err != nil {
		return
	}
	defer file.Close()

	writer := json.NewEncoder(file)
	for _, value := range data {
		if err = writer.Encode(value); err != nil {
			return
		}
	}
	return
}

// xmlDownload downloads the inventory to a xml file via the web browser
func xmlDownload(w http.ResponseWriter, filename string, data WmsList) (err error) {
	w.Header().Set("Content-Type", "text")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=%s", filename))
	output, err := xml.MarshalIndent(data, " ", "   ")
	if err != nil {
		log.Println(err)
	}
	w.Write(output)

	return
}
