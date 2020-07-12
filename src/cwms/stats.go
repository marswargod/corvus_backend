package main

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
