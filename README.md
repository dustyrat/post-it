#### post-it
```
post-it is a HTTP(S) CLI library for calling a variaty of urls from an input file.

Usage:
  post-it [command]

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
  -s, --response-status string   Record response status to output file under the headers status. eg: any, none, 2xx, -2xx (non 2xx statuses), 4xx, 5xx, 200, 301, 404, 503... (default "-2xx")
  -t, --timeout duration         Connection timeout (default 3s)
  -u, --url string               Url. Should be in the format 'http://localhost:3000/path/{column_name}' if input file is specified

Use "post-it [command] --help" for more information about a command.
```
