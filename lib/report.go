package happa

import (
	"fmt"
	"math"
	"sort"
	"time"
)

type Result struct {
	statusCode int
	latency time.Duration
}

type Report struct {
	results []Result
	start time.Time
	end time.Time
	rc *int
}

func NewReport(start time.Time, rc *int) *Report{
	return &Report{
		results: []Result{},
		start: start,
		rc: rc,
	}
}

func (r *Report) receive(resultCh <- chan *Result) {
	for {
		select {
		case res, ok := <- resultCh:
			if !ok {
				r.end = time.Now()
				return
			}
			r.results = append(r.results, *res)
		}
	}
}

func (r *Report) calcTileIndex(len int, tile float64) int64 {
	return int64(math.Floor(float64(len) * tile))
}

func (r *Report) calcPercentileLatency(i int64, latencies []time.Duration) time.Duration {
	if i == 0 {
		return time.Duration(0)
	}
	var sum time.Duration
	for _, l := range latencies[:i] {
		sum += l
	}
	return time.Duration(sum.Nanoseconds() / i)
}

func (r *Report) output() {
	var sumLatency time.Duration
	var maxLatency time.Duration
	resultsNumber := len(r.results)
	latencies := make([]time.Duration, resultsNumber)
	statusCounts := make(map[int]int)

	for i, res := range r.results {
		sumLatency += res.latency
		latencies[i] = res.latency
		if res.latency > maxLatency {
			maxLatency = res.latency
		}
		statusCounts[res.statusCode] += 1
	}

	sort.Slice(latencies, func(i, j int) bool { return latencies[i] < latencies[j] })
	i50 := r.calcTileIndex(resultsNumber, 0.5)
	i90 := r.calcTileIndex(resultsNumber, 0.9)
	i99 := r.calcTileIndex(resultsNumber, 0.99)

	fmt.Println("[Requests]")
	fmt.Printf("%-8s %12d\n", "total", *r.rc)
	fmt.Printf("%-8s %12.3f\n", "rate", float64(*r.rc) / r.end.Sub(r.start).Seconds())
	fmt.Println("[Duration]")
	fmt.Printf("%-8s %12s\n", "total", r.end.Sub(r.start))
	fmt.Println("[Status]")
	for code, count := range statusCounts {
		fmt.Printf("%d:%d\n", code, count)
	}
	fmt.Println("[Latency]")
	fmt.Printf("%-8s %12s\n", "avg", time.Duration(sumLatency.Nanoseconds() / int64(resultsNumber)))
	fmt.Printf("%-8s %12s\n", "max", maxLatency)
	fmt.Printf("%-8s %12s\n", "50%tile", r.calcPercentileLatency(i50, latencies))
	fmt.Printf("%-8s %12s\n", "90%tile", r.calcPercentileLatency(i90, latencies))
	fmt.Printf("%-8s %12s\n", "99%tile", r.calcPercentileLatency(i99, latencies))
}