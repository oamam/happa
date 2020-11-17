# happa
This is a simple load testing tool.

## Usage

```
flag needs an argument: -h
Usage of ./happa:
  -d int
        duration (default 60)
  -h value
        request header
  -m string
        http method (default "GET")
  -r int
        request rate (rps) (default 10)
  -u string
        target url
  -w int
        worker thread number (default 5)
```

## Report Example

```
[Requests]
total             300
rate           89.217
[Duration]
total    3.362581716s
[Status]
200:300
[Latency]
avg       55.583952ms
max      287.348957ms
50%tile   46.501472ms
90%tile    49.55348ms
99%tile   53.266449ms
```
