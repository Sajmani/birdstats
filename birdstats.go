// The birdstats command prints stats for a collection of eBird submissions in a CSV file.
// It also prints per-species stats.
package main

import (
	"fmt"
	"log"
	"maps"
	"os"
	"slices"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/Sajmani/birdstats/ebird"
)

// submission contains per-submission stats
type submission struct {
	species int
	km      float64
	dur     time.Duration
}

// species contains per-species stats
type species struct {
	commonName  string
	submissions int
	count       int
	mlAssets    int
}

func main() {
	recs, err := ebird.Records(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	subs := make(map[string]*submission)
	specs := make(map[string]*species)
	for rec := range recs {
		sub := subs[rec.SubmissionID]
		if sub == nil {
			km, err := strconv.ParseFloat(rec.DistanceTraveledKm, 64)
			if err != nil && rec.DistanceTraveledKm != "" {
				log.Fatal(rec.Line, err)
			}
			durMin, err := strconv.Atoi(rec.DurationMin)
			if err != nil && rec.DurationMin != "" {
				log.Fatal(rec.Line, err)
			}
			sub = &submission{
				km:  km,
				dur: time.Duration(durMin) * time.Minute,
			}
			subs[rec.SubmissionID] = sub
		}
		sub.species++

		spec := specs[rec.ScientificName]
		if spec == nil {
			spec = &species{
				commonName: rec.CommonName,
			}
			specs[rec.ScientificName] = spec
		}
		spec.submissions++
		count, err := strconv.Atoi(rec.Count)
		if err != nil && rec.Count != "" && rec.Count != "X" {
			log.Fatal(rec.Line, err)
		}
		spec.count += count
		if mlAssets := strings.TrimSpace(rec.MLCatalogNumbers); mlAssets != "" {
			spec.mlAssets += len(strings.Split(mlAssets, " "))
		}
	}

	var (
		kmTotal  float64
		durTotal time.Duration
	)
	for _, sub := range subs {
		kmTotal += sub.km
		durTotal += sub.dur
	}
	miTotal := kmTotal * 0.621371

	fmt.Printf("%d species; %d submissions\n", len(specs), len(subs))
	fmt.Printf("%.0f total km; %.1f avg km\n", kmTotal, kmTotal/float64(len(subs)))
	fmt.Printf("%.0f total mi; %.1f avg mi\n", miTotal, miTotal/float64(len(subs)))
	fmt.Printf("%s total time; %s avg time\n", durTotal, (durTotal / time.Duration(len(subs))).Truncate(time.Minute))
	fmt.Printf("————————————————————————————————————————————————————————————————————————————————\n")

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	for _, name := range slices.Sorted(maps.Keys(specs)) {
		spec := specs[name]
		fmt.Fprintln(w, name, "\t"+spec.commonName+"\t", spec.count, "total\t", spec.mlAssets, "pics/sounds")
	}
	w.Flush()
}
