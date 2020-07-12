package main

import "fmt"

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

// toSqlStmt generates a sql statement based on the current set of page controls
func (mc MissionControls) toSqlStmt() (sqlstmt string) {
	var sel, order string
	sel = `select entry, region, frequency, IFNULL(aisle, ""), IFNULL(block,""), IFNULL(slot,"") from v_schedule `
	order = `order by entry`
	sqlstmt = fmt.Sprintf("%s %s", sel, order)
	return
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

// missionNav creates a new page controls and initializes the mission navigation portion of mission controls based on curr aisle and aisle list
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
