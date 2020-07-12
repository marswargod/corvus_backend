// 	Corvus Warehouse Management System (cwms)
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
	"io/ioutil"
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

// fetchAisles performs a query on v_inventory and returns the results in a aisleList
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
		wl, err := FetchInventory(pc.toAisleFilter())
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

// PageControls holds page navigation and page control fields
type PageControls struct {
	Curr, Next, Prev string   // Page Nav
	SingleAisle      bool     // Has a single aisle been selected?
	Selection        string   // selection choice "all" or an aisle name
	Scope            string   // filter scope: blank or "issues"
	Aisles           []string // list of aisles
}

// MissionControls holds page navigation and mission control fields
type MissionControls struct {
	Curr, Next, Prev       string   // Mission Nav
	SingleDay              bool     // Has a single day been selected?
	Selection              string   // selection choice "all" or a specific date
	Scope                  string   // filter scope: blank or "issues"
	Days                   []string // list of days
	CurrentStatus          string   // Charging, Waiting, In Flight
	BatteryLevel           string   // Battery charge level
	TimeUntilNextFlight    string   // Time until next flight (when battery is charged)
	LastCompleteInventory  string   //
	DaysLeftInCurrentCycle string   //
	AveDaysToCompleteCycle string   //
}

func (pc PageControls) toAisleFilter() (af AisleFilter) {
	if pc.SingleAisle {
		af.Aisle = pc.Curr
	}
	if pc.Scope != "" {
		af.Discrepancy = "all"
	}
	return
}

// toSqlStmt generates a sql statement based on the current set of page controls
func (mc MissionControls) toSqlStmt() (sqlstmt string) {
	var sel, order string
	sel = `select entry, region, frequency, IFNULL(aisle, ""), IFNULL(block,""), IFNULL(slot,"") from v_schedule `
	order = `order by entry`
	sqlstmt = fmt.Sprintf("%s %s", sel, order)
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

// missionControls generates a set of mission nav and mission controls based on day and scope
func missionControls(day, scope string) (mc MissionControls, err error) {
	var days []string
	days, err = fetchDays(scope)
	if err != nil {
		return
	}
	if day == "" {
		day = days[0]
	}
	mc = missionNav(day, days)
	mc.Scope = scope
	mc.Selection = day
	mc.SingleDay = day != "all"
	mc.BatteryLevel = "50%"
	mc.CurrentStatus = "Charging"
	mc.TimeUntilNextFlight = "15 minutes"
	mc.LastCompleteInventory = "April 15, 2020"
	mc.DaysLeftInCurrentCycle = "20 days"
	mc.AveDaysToCompleteCycle = "45 day"
	return
}

// pageNav creates a new page controls and initializes the page navigation portion of page controls based on curr aisle and aisle list
func missionNav(curr string, dl []string) (mc MissionControls) {
	mc = MissionControls{Days: dl, Curr: dl[0], Next: dl[0], Prev: dl[0]}
	for d, v := range dl {
		if v == curr {
			mc.Curr = dl[d]
			mc.Next = dl[Min(len(dl)-1, d+1)]
			mc.Prev = dl[Max(0, d-1)]
			break
		}
	}
	return
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
	db, err = sql.Open("sqlite3", `C:\Users\James\go\src\cwms\wms3.db`)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	// Setup servemux to serve http handler routines
	mux := http.NewServeMux()

	// Setup file server handler
	files := http.FileServer(LocalFileSystem{http.Dir(`C:\Users\James\go\src\cwms\static\`)})
	mux.Handle("/static/", http.StripPrefix("/static/", files))

	// Setup http handlers
	mux.HandleFunc("/", imw(handleDashboard))
	mux.HandleFunc("/dashboard/", imw(handleDashboard))
	mux.HandleFunc("/inventory/", imw(handleInventory))
	mux.HandleFunc("/hybrid/", imw(handleHybrid))
	mux.HandleFunc("/schedule/", mmw(handleSchedule))
	mux.HandleFunc("/export/csv/", imw(handleExportInventoryCsv))
	mux.HandleFunc("/export/json/", imw(handleExportInventoryJson))
	mux.HandleFunc("/api/json/", imw(handleApiInventoryJson))
	mux.HandleFunc("/export/xml/", imw(handleExportInventoryXml))
	mux.HandleFunc("/api/", handleJsonApiRequest) // handleApiAisles
	mux.HandleFunc("/api/aisles/", handleApiAisles)
	mux.HandleFunc("/api/discrepancies/", handleApiDiscrepancies)
	mux.HandleFunc("/api/restrictions/", handleApiRestrictions)

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

// handleSchedule provides basic navigation features and downloads files in csv, json, or xml formats
func handleSchedule(tm map[string]interface{}, ml MmsList, w http.ResponseWriter, r *http.Request) {
	err := executeTemplate("schedule.html", tm, w)
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

// handleApiInventoryJson transfers the inventory via a restful api in a json format
func handleApiInventoryJson(tm map[string]interface{}, wl WmsList, w http.ResponseWriter, r *http.Request) {
	err := jsonApi(w, wl)
	if err != nil {
		log.Println(err)
	}
}

// // handleApiRequestJson interprets a json request
// func handleApiRequestJson(tm map[string]interface{}, wl WmsList, w http.ResponseWriter, r *http.Request) {
// 	err := jsonExternalApi(w, r)
// 	if err != nil {
// 		log.Println(err)
// 	}
// }

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
	if err = json.NewEncoder(w).Encode(data); err != nil {
		log.Println(err)
	}
	return
}

// jsonApi implements a simple restful api to export inventory in a json format
func jsonApi(w http.ResponseWriter, data interface{}) (err error) {
	if err = json.NewEncoder(w).Encode(data); err != nil {
		log.Println(err)
	}
	return
}

type Pick struct {
	SKU   string `json:"sku"`
	Aisle string `json:"aisle"`
	Shelf string `json:"shelf"`
	Slot  string `json:"slot"`
}

// external api struct
type WMSActions struct {
	Picks []Pick
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

// handleJsonApiRequest implements a simple api to read a wms list in a json format
func handleJsonApiRequest(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}
	var wmsl WmsList
	err = json.Unmarshal(body, &wmsl)
	if err != nil {
		log.Println(err)
	}
	log.Println(wmsl)
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
