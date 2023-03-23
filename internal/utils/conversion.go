package utils

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

// daysIn calculates and returns how many days in a month for a given year
func DaysIn(month, year int) int {
	return time.Date(year, time.Month(month+1), 0, 0, 0, 0, 0, time.UTC).Day()
}

// renderJSON renders `v` as a JSON and writes it in a response writer
func RenderJSON(w http.ResponseWriter, v any) {
	json, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

// `StrToSlice` slices `str` into all substrings separated by `,` and returns
// slice of the substrings between those separators
func StrToSlice(str string) []string {
	return strings.Split(str, ",")
}

func StrToTime(str string) time.Time {
	// TODO: implement
	return time.Now()
}
