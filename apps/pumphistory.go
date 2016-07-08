package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/ecc1/medtronic"
)

var (
	all      = flag.Bool("a", false, "get entire pump history")
	numHours = flag.Int("n", 6, "number of `hours` of history to get")
)

func main() {
	flag.Parse()
	pump := medtronic.Open()
	pump.Wakeup()
	newer := pump.Family() >= 23
	numPages := pump.HistoryPageCount()
	cutoff := time.Now().Add(-time.Duration(*numHours) * time.Hour)
	if *all {
		log.Printf("retrieving %d pages of records", numPages)
	} else {
		log.Printf("retrieving records since %s", cutoff.Format(medtronic.TimeLayout))
	}
	results := []medtronic.HistoryRecord{}
loop:
	for page := 0; page < numPages && pump.Error() == nil; page++ {
		log.Printf("scanning page %d", page)
		data := pump.HistoryPage(page)
		records, err := medtronic.DecodeHistoryRecords(data, newer)
		if err != nil {
			pump.SetError(err)
		}
		for _, r := range records {
			t := r.Time
			if !*all && !t.IsZero() && t.Before(cutoff) {
				log.Printf("stopping at timestamp %s", t.Format(medtronic.TimeLayout))
				break loop
			}
			results = append(results, r)
		}
	}
	if pump.Error() != nil {
		log.Fatal(pump.Error())
	}
	b, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
}
