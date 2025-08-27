// Package ebird provides helper functions for working with eBird data.
package ebird

import (
	"encoding/csv"
	"iter"
	"log"
	"os"
	"time"
)

// Record contains the fields in MyEBirdData.csv records.
type Record struct {
	Line               int // line in the CSV file
	SubmissionID       string
	CommonName         string
	ScientificName     string
	TaxonomicOrder     string
	Count              string // "X" or integer
	StateProvince      string
	County             string
	LocationID         string
	Location           string
	Latitude           string
	Longitude          string
	Date               string // YYYY-MM-DD
	Time               string // 07:00 AM
	Protocol           string
	DurationMin        string
	AllObsReported     string // "1" means yes
	DistanceTraveledKm string
	AreaCoveredHa      string
	NumberOfObservers  string
	BreedingCode       string
	ObservationDetails string
	ChecklistComments  string
	MLCatalogNumbers   string
}

func (r Record) Observed() (time.Time, error) {
	if r.Time == "" {
		return time.Parse("2006-01-02", r.Date)
	}
	return time.Parse("2006-01-02 03:04 PM", r.Date+" "+r.Time)
}

func Records(filename string) (iter.Seq[Record], error) {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("ebird.Records(%s): %v", filename, err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	// eBird's CSV export returns a variable number of fields per record,
	// so disable this check. This means we need to explicitly check len(rec)
	// before accessing fields that might not be there.
	r.FieldsPerRecord = -1
	recs, err := r.ReadAll()
	if err != nil {
		log.Fatalf("Error reading CSV records from %s: %v", filename, err)
	}
	if len(recs) < 1 {
		log.Fatalf("No records found in %s", filename)
	}
	field := make(map[string]int)
	for i, f := range recs[0] {
		field[f] = i
	}
	recs = recs[1:]
	log.Printf("Read %d eBird observations", len(recs))
	return func(yield func(Record) bool) {
		for i, rec := range recs {
			stringField := func(key string) string {
				if f := field[key]; f < len(rec) {
					return rec[f]
				}
				return ""
			}
			if !yield(Record{
				Line:               i + 2, // header was line 1
				SubmissionID:       stringField("Submission ID"),
				CommonName:         stringField("Common Name"),
				ScientificName:     stringField("Scientific Name"),
				TaxonomicOrder:     stringField("Taxonomic Order"),
				Count:              stringField("Count"),
				StateProvince:      stringField("State/Province"),
				County:             stringField("County"),
				LocationID:         stringField("Location ID"),
				Location:           stringField("Location"),
				Latitude:           stringField("Latitude"),
				Longitude:          stringField("Longitude"),
				Date:               stringField("Date"),
				Time:               stringField("Time"),
				Protocol:           stringField("Protocol"),
				DurationMin:        stringField("Duration (Min)"),
				AllObsReported:     stringField("All Obs Reported"),
				DistanceTraveledKm: stringField("Distance Traveled (km)"),
				AreaCoveredHa:      stringField("Area Covered (ha)"),
				NumberOfObservers:  stringField("Number of Observers"),
				BreedingCode:       stringField("Breeding Code"),
				ObservationDetails: stringField("Observation Details"),
				ChecklistComments:  stringField("Checklist Comments"),
				MLCatalogNumbers:   stringField("ML Catalog Numbers"),
			}) {
				return
			}
		}
	}, nil
}
