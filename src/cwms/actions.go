package main

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
