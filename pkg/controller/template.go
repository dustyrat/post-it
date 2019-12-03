package controller

const (
	text = `
{{ printf "%-10v %10v %10v %10v" "Statistics" "Avg" "Stdev" "Max" }}
{{ printf "  %-10v %10.2f %10.2f %10.2f" "Reqs/sec" .Request.Mean .Request.Stddev .Request.Max }}
{{ printf "  %-10v %10v %10v %10v" "Latency" .Response.Mean .Response.Stddev .Response.Max -}}
{{ print "\n  Latency Distribution" -}}
{{ range $pc, $lat := .Response.Percentiles -}}
	{{ printf "\n     %2.0f%% %10s" (Multiply $pc 100) $lat -}}
{{ end }}
`
	text1 = `
{{- printf "%10v %10v %10v %10v" "Statistics" "Avg" "Stdev" "Max" }}
{{ with .Results (FloatsToArray 0.5 0.75 0.9 0.95 0.99) }}
	{{- printf "  %-10v %10.2f %10.2f %10.2f" "Reqs/sec" .Mean .Stddev .Max -}}
{{ else }}
	{{- print "  There wasn't enough data to compute statistics for requests." }}
{{ end }}
{{ with .Results.Percentiles (FloatsToArray 0.5 0.75 0.9 0.95 0.99) }}
	{{- printf "  %-10v %10v %10v %10v" "Latency" (FormatTimeUs .Mean) (FormatTimeUs .Stddev) (FormatTimeUs .Max) }}
	{{- if WithLatencies }}
  		{{- "\n  Latency Distribution" }}
		{{- range $pc, $lat := .Percentiles }}
			{{- printf "\n     %2.0f%% %10s" (Multiply $pc 100) (FormatTimeUsUint64 $lat) -}}
		{{ end -}}
	{{ end }}
{{ else }}
	{{- print "  There wasn't enough data to compute statistics for latencies." }}
{{ end -}}
{{ with .Result -}}
{{ "  HTTP codes:" }}
{{ printf "    1xx - %v, 2xx - %v, 3xx - %v, 4xx - %v, 5xx - %v" .Req1XX .Req2XX .Req3XX .Req4XX .Req5XX }}
	{{- printf "\n    others - %v" .Others }}
	{{- with .Errors }}
		{{- "\n  Errors:"}}
		{{- range . }}
			{{- printf "\n    %10v - %v" .Error .Count }}
		{{- end -}}
	{{ end -}}
{{ end }}
{{ printf "  %-10v %10v/s\n" "Throughput:" (FormatBinary .Result.Throughput)}}`
)
