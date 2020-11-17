package main

import (
	"flag"
	happa "github.com/oamam/happa/lib"
	"log"
	"net/http"
)

func main() {
	defer func(){
		if err := recover(); err != nil {
			log.Fatal("Caught error : ", err)
		}
	}()

	url := flag.String("u", "", "target url")
	method := flag.String("m", "GET", "http method")
	var headers happa.Headers
	flag.Var(&headers, "h", "request header")
	duration := flag.Int("d", 60, "duration")
	rate := flag.Int("r", 10, "request rate (rps)")
	workerNumber :=  flag.Int("w", 5, "worker thread number")
	flag.Parse()

	a := happa.NewAttacker(&http.Client{}, url, method, headers, duration, rate, workerNumber)
	if err := a.Run(); err != nil {
		log.Fatal("Caught error : ", err)
	}
}