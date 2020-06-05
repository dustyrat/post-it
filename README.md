
# post-it
[![Go Report Card](https://goreportcard.com/badge/github.com/DustyRat/post-it)](https://goreportcard.com/report/github.com/DustyRat/post-it)
[![GoDoc](https://godoc.org/github.com/DustyRat/post-it?status.svg)](http://godoc.org/github.com/DustyRat/post-it)
[![Coverage](https://gocover.io/_badge/github.com/DustyRat/post-it)](https://gocover.io/github.com/DustyRat/post-it)

post-it is a HTTP(S) testing tool. It is written in Go programming language.
> Insperation taken from https://www.postman.com/ and https://github.com/codesenberg/bombardier


## Installation
You can grab binaries in the [releases](https://github.com/DustyRat/post-it/releases) section.

Alternatively, to get latest and greatest run: `go get -u github.com/DustyRat/post-it`

## Usage

```
post-it is a HTTP(S) CLI library for calling a variaty of urls from an input file.

All methods use the request_body column for requests.

Usage:
  post-it [command] <url> [flags]

Available Commands:
  DELETE      The DELETE method deletes the specified resource.
  GET         The HTTP GET method requests a representation of the specified resource.
  HEAD        The HEAD method asks for a response identical to that of a GET request, but without the response body.
  PATCH       The PATCH method is used to apply partial modifications to a resource.
  POST        The POST method is used to submit an entity to the specified resource, often causing a change in state or side effects on the server.
  PUT         The PUT method replaces all current representations of the target resource with the request payload.
  help        Help about any command

Flags:
  -c, --connections int          Concurrent connections (default 10)
  -e, --errors                   Record erorrs to output file
  -H, --header stringArray       HTTP headers to use ("K: V")
  -h, --help                     help for post-it
  -g, --histogram                Print histogram statistics
  -i, --input string             Input File (default "input.csv")
      --insecure                 Insecure Skip Verify (default true)
  -l, --latencies                Print latency statistics
  -o, --output string            Output File (default "output.csv")
  -b, --record-body              Record body to output file under the response_body column.
      --record-headers           Record headers to output file under the headers column.
  -s, --response-status string   Record response status to output file under the headers status. eg: any, none, 2xx, -2xx (non 2xx statuses), 4xx, 5xx, 200, 301, 404, 503... (default "-2xx")
  -t, --timeout duration         Connection timeout (default 3s)

Use "post-it [command] --help" for more information about a command.
```

## Examples
Basic:
> Simple STD output. Any non 2xx responses will be saved in output.csv.
> Uses default input file 'input.csv'
```
post-it GET "http://localhost:3000/get/{id}"
```
Output:
```
10 / 10  [=====================================] complete  12.5/s Elapsed: 0s   

Responses
   OK: 200 | 
        10 | 
Statistics
           |     Average |      STDDEV | Max
   Req/sec |       12.49 |          NA | NA
   Latency |    497.95ms |    246.37ms | 799.77ms
```

Specific Input:
> Uses './test/input.csv' as input file
```
post-it GET "http://localhost:3000/get/{id}" -i ./test/input.csv
```
Output:
```
10 / 10  [=====================================] complete  10.4/s Elapsed: 0s   

Responses
   OK: 200 | 
        10 | 
Statistics
           |     Average |      STDDEV | Max
   Req/sec |       10.41 |          NA | NA
   Latency |    599.45ms |    294.34ms | 959.44ms
```
---
Specific Output:
> Outputs request results to './results.csv'
```
post-it GET "http://localhost:3000/get/{id}" -o ./results.csv
```
Output:
```
10 / 10  [=====================================] complete  9.5/s Elapsed: 1s   

Responses
   OK: 200 | 
        10 | 
Statistics
           |     Average |      STDDEV | Max
   Req/sec |        9.51 |          NA | NA
   Latency |    597.58ms |    341.38ms | 1.05s
```
---
Specific Status:
> Outputs 2xx request results to './output.csv'
```
post-it GET "http://localhost:3000/get/{id}" -s 2xx
```
Output:
```
10 / 10  [=====================================] complete  11.5/s Elapsed: 0s   

Responses
   OK: 200 | 
        10 | 
Statistics
           |     Average |      STDDEV | Max
   Req/sec |       11.51 |          NA | NA
   Latency |    531.55ms |    251.08ms | 868.61ms
```
---
Errors Only:
> Output only client errors to './output.csv'
```
post-it GET "http://localhost:3000/get/{id}" -e -s none
```
Output:
```
10 / 10  [=====================================] complete  1810.6/s Elapsed: 0s   

Responses
   Errors | 
       10 | 
Statistics
           |    Average |    STDDEV | Max
   Req/sec |       0.00 |        NA | NA
   Latency |         0s |        0s | 0s
```
---
Latency & Historgram:
> Outputs request results to './results.csv'
```
post-it GET "http://localhost:3000/get/{id}" -lg
```
Output:
```
10 / 10  [=====================================] complete  9.8/s Elapsed: 1s   

Responses
   OK: 200 | 
        10 | 
Statistics
           |     Average |     STDDEV | Max
   Req/sec |        9.84 |         NA | NA
   Latency |    580.18ms |    292.9ms | 1.02s
Latency Distibution
   50.00% | 515.75ms
   75.00% | 870.99ms
   90.00% | 894.22ms
   95.00% | 1.02s
   99.00% | 1.02s
Histogram
   Bucket | Count
      1ms | 0
    2.5ms | 0
      5ms | 0
    7.5ms | 0
     10ms | 0
     25ms | 0
     50ms | 0
     75ms | 1
    100ms | 0
    250ms | 3
    500ms | 2
    750ms | 3
       1s | 1
     2.5s | 0
       5s | 0
     7.5s | 0
      10s | 0
     Inf+ | 0
```

## Licensing
"The code in this project is licensed under MIT license."