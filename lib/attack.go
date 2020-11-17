package happa

import (
	"math"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Attacker struct {
	client *http.Client
	req *http.Request
	url *string
	method *string
	headers []string
	body []byte
	duration *int
	rate *int
	workerNumber *int
}

func NewAttacker(client *http.Client, url, method *string, headers []string, duration, rate, workerNumber *int) *Attacker {
	return &Attacker {
		client: client,
		url: url,
		method: method,
		headers: headers,
		duration: duration,
		rate: rate,
		workerNumber: workerNumber,
	}
}

func (a *Attacker) shouldSendSignal(rc *int, elapsed *time.Duration) bool {
	if *a.rate > 0 {
		return *a.duration * *a.rate > *rc
	}
	return int64(time.Duration(*a.duration) * time.Second) > elapsed.Nanoseconds()
}

func (a *Attacker) calcWaitTime(rc *int, elapsed *time.Duration) time.Duration {
	if *a.rate <= 0 || elapsed.Nanoseconds() * int64(*a.rate) > int64(*rc) {
		return 0
	}
	return time.Duration(math.Round(1e9 / float64(*a.rate)))
}

func (a *Attacker) sendSignal(start time.Time, rc *int, wg *sync.WaitGroup, signalCh chan <- struct{}, resultCh chan <- *Result) {
	defer close(resultCh)
	defer wg.Wait()
	defer close(signalCh)

	for {
		current := time.Now()
		elapsed := current.Sub(start)
		if !a.shouldSendSignal(rc, &elapsed) {
			return
		}
		t := a.calcWaitTime(rc, &elapsed)
		time.Sleep(t)
		select {
		case signalCh <- struct{}{}:
			*rc++
		default:
		}
	}
}

func (a *Attacker) receiveSignal(wg *sync.WaitGroup, signalCh <- chan struct{}, resultCh chan <- *Result) {
	defer wg.Done()
	for range signalCh {
		resultCh <- a.attack()
	}
}

func (a *Attacker) attack() *Result {
	start := time.Now()
	res, err := a.client.Do(a.req)
	end := time.Now()
	if err != nil {
		return &Result{}
	}
	latency := end.Sub(start)
	return &Result{
		statusCode: res.StatusCode,
		latency: latency,
	}
}

func (a *Attacker) parseHeader() http.Header {
	var h http.Header = map[string][]string{}
	for _, s := range a.headers {
		ss := strings.SplitN(s, ":", 2)
		if len(ss) != 2 {
			continue
		}
		h[ss[0]] = []string{strings.TrimSpace(ss[1])}
	}
	return h
}


func (a *Attacker) makeRequest() error {
	req, err := http.NewRequest(*a.method, *a.url, nil)
	if err != nil {
		return err
	}
	req.Header = a.parseHeader()
	a.req = req
	return nil
}

func (a *Attacker) Run() error {
	var wg sync.WaitGroup
	signalCh := make(chan struct{})
	resultCh := make(chan *Result)

	if err := a.makeRequest(); err != nil {
		return err
	}

	start := time.Now()
	rc := 0
	go a.sendSignal(start, &rc, &wg, signalCh, resultCh)
	for i := 0; i < *a.workerNumber; i++{
		wg.Add(1)
		go a.receiveSignal(&wg, signalCh, resultCh)
	}

	results := NewReport(start, &rc)
	results.receive(resultCh)
	results.output()

	return nil
}