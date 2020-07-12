package main

// PageControls holds page navigation and page control fields
type PageControls struct {
	Curr, Next, Prev string   // Page Nav
	SingleAisle      bool     // Has a single aisle been selected?
	Selection        string   // selection choice "all" or an aisle name
	Scope            string   // filter scope: blank or "issues"
	Aisles           []string // list of aisles
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
