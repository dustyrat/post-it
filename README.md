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
  -c, --connections int          connections (default 10)
      --header stringArray       HTTP headers to use ("K: V")
  -h, --help                     help for post-it
      --idle-timeout duration    Idle Connection timeout (default 500ms)
  -i, --input string             Input File (default "input.csv")
      --insecure-skip-verify     Insecure Skip Verify (default true)
  -o, --output string            Output File (default "output.csv")
      --record-body              Output body
      --record-headers           Output headers
      --response-status string   Response status to output. eg: any, 2xx, -2xx (non 2xx statuses), 4xx, 5xx, 200, 301, 404, 503... (default "any")
      --response-type string     Response type to output. eg: all, error, status, none (default "none")
  -t, --timeout duration         Connection timeout (default 3s)
  -u, --url string               Url. Should be in the format 'http://localhost:3000/path/{column_name}' if input file is specified

Use "post-it [command] --help" for more information about a command.
```

#### DELETE
```
The DELETE method deletes the specified resource.

Usage:
  post-it DELETE [flags]

Aliases:
  DELETE, delete

Examples:
post-it DELETE -u http://localhost:3000/path/{column_name}

Flags:
  -h, --help   help for DELETE
```

#### GET
```
The HTTP GET method requests a representation of the specified resource.

Usage:
  post-it GET [flags]

Aliases:
  GET, get

Examples:
post-it GET -u http://localhost:3000/path/{column_name}

Flags:
  -h, --help   help for GET
```

#### HEAD
```
The HEAD method asks for a response identical to that of a GET request, but without the response body.

Usage:
  post-it HEAD [flags]

Aliases:
  HEAD, head

Examples:
post-it HEAD -u http://localhost:3000/path/{column_name}

Flags:
  -h, --help   help for HEAD
```

#### PATCH
```
The PATCH method is used to apply partial modifications to a resource.

Usage:
  post-it PATCH [flags]

Aliases:
  PATCH, patch

Examples:
post-it PATCH -u http://localhost:3000/path/{column_name}

Flags:
  -h, --help   help for PATCH
```

#### POST
```
The POST method is used to submit an entity to the specified resource, often causing a change in state or side effects on the server.

Usage:
  post-it POST [flags]

Aliases:
  POST, post

Examples:
post-it POST -u http://localhost:3000/path/{column_name}

Flags:
  -h, --help   help for POST
```

#### PUT
```
The PUT method replaces all current representations of the target resource with the request payload.

Usage:
  post-it PUT [flags]

Aliases:
  PUT, put

Examples:
post-it PUT -u http://localhost:3000/path/{column_name}

Flags:
  -h, --help   help for PUT
```