package common

type Page struct {
	Count int64 `json:"count"`
	Rows  any   `json:"rows"`
}
