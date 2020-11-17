package happa

import (
	"testing"
	"time"
)

func TestCalcPercentileLatency(t *testing.T) {
	rc := 0
	r := NewReport(time.Now(), &rc)
	latencies := []time.Duration {
		time.Duration(10) * time.Millisecond,
		time.Duration(10) * time.Millisecond,
		time.Duration(10) * time.Millisecond,
		time.Duration(50) * time.Millisecond,
		time.Duration(100) * time.Millisecond,
	}
	for _, test := range []struct {
		name string
		expected time.Duration
		tile float64
	} {
		{"10%tile", 0 * time.Millisecond, 0.1},
		{"20%tile", 10 * time.Millisecond, 0.2},
		{"30%tile", 10 * time.Millisecond, 0.3},
		{"40%tile", 10 * time.Millisecond, 0.4},
		{"50%tile", 10 * time.Millisecond, 0.5},
		{"60%tile", 10 * time.Millisecond, 0.6},
		{"70%tile", 10 * time.Millisecond, 0.7},
		{"80%tile", 20 * time.Millisecond, 0.8},
		{"90%tile", 20 * time.Millisecond, 0.9},
		{"99%tile", 20 * time.Millisecond, 0.99},
	} {
		t.Run(test.name, func(t *testing.T) {
			actual := r.calcPercentileLatency(r.calcTileIndex(5, test.tile), latencies)
			if actual != test.expected {
				t.Errorf("failed : actual = %s, expected = %s", actual, test.expected)
			}
		})
	}
}

func TestOutputEmptyResults(t *testing.T) {
	c := 100
	rs := make([]Result, c)
	for i := 0; i < c; i++ {
		rs[i] = Result{}
	}
	rc := 100
	r := NewReport(time.Now(), &rc)
	r.results = rs
	r.end = time.Now().Add(1 * time.Second)
	r.output()
}