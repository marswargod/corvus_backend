package main

import (
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

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
	err := jsonApi(w, r, wl)
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
	if err = json.NewEncoder(w).Encode(data); err != nil {
		log.Println(err)
	}
	return
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
