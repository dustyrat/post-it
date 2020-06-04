package stats

import (
	"fmt"
	"math"
	"net/http"
	"os"
	"sort"
	"strconv"
	"text/tabwriter"
	"time"

	internal "github.com/DustyRat/post-it/internal/http"
	"github.com/DustyRat/post-it/internal/options"

	io_prometheus_client "github.com/prometheus/client_model/go"
)

// Print ...
func Print(opts options.Options, elapsed time.Duration) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.AlignRight|tabwriter.Debug)

	var summaries, histograms, counters []*io_prometheus_client.Metric
	metrics, _ := internal.Gatherer.Gather()
	for _, metric := range metrics {
		switch metric.GetType() {
		case io_prometheus_client.MetricType_COUNTER:
			counters = metric.GetMetric()
		case io_prometheus_client.MetricType_HISTOGRAM:
			histograms = metric.GetMetric()
		case io_prometheus_client.MetricType_SUMMARY:
			summaries = metric.GetMetric()
		}
	}

	histgram := make(map[time.Duration]uint64)
	hbuckets := make(sort.IntSlice, 0)
	var count uint64
	var average, max time.Duration
	var rate float64
	for _, h := range histograms {
		histogram := h.GetHistogram()
		average = time.Duration(histogram.GetSampleSum()*float64(time.Second)) / time.Duration(histogram.GetSampleCount())
		rate = float64(histogram.GetSampleCount()) / elapsed.Seconds()
		count = histogram.GetSampleCount()
		for _, bucket := range histogram.GetBucket() {
			upperBound := time.Duration(bucket.GetUpperBound() * float64(time.Second))
			count := bucket.GetCumulativeCount()
			histgram[upperBound] = count
			hbuckets = append(hbuckets, int(upperBound))
		}
	}
	hbuckets.Sort()

	for i := range hbuckets {
		current := time.Duration(hbuckets[i])
		if i < len(hbuckets)-1 {
			next := time.Duration(hbuckets[i+1])
			histgram[current] = histgram[next] - histgram[current]
			continue
		}
		histgram[current] = count - histgram[current]
	}

	var cumlative float64
	for bucket, count := range histgram {
		cumlative += float64(count) * math.Pow(float64(bucket-average), 2)
	}
	stddev := time.Duration(math.Sqrt(cumlative / float64(count)))

	quantiles := make(map[float64]time.Duration)
	qbuckets := make(sort.Float64Slice, 0)
	for _, s := range summaries {
		summary := s.GetSummary()
		for _, quantile := range summary.GetQuantile() {
			value := time.Duration(quantile.GetValue() * float64(time.Second))
			if quantile.GetQuantile() < 1 {
				bucket := quantile.GetQuantile()
				quantiles[bucket] = value
				qbuckets = append(qbuckets, bucket)
			} else {
				max = value
			}
		}
	}
	qbuckets.Sort()

	fmt.Fprintln(w, "\nResponses")
	statuses := make(map[int]int)
	codes := make(sort.IntSlice, 0)
	for _, counter := range counters {
		var code int
		for _, label := range counter.GetLabel() {
			if label.GetName() == "code" {
				code, _ = strconv.Atoi(label.GetValue())
				codes = append(codes, code)
				break
			}
		}
		value := counter.Counter.GetValue()
		statuses[code] = int(value)
	}
	codes.Sort()

	headers, line := "", ""
	values := make([]interface{}, 0)
	for _, code := range codes {
		if code != 0 {
			headers += fmt.Sprintf("%s: %d \t ", http.StatusText(code), code)
			line += "%d \t "
			values = append(values, statuses[code])
		} else {
			headers += "Errors \t "
			line += "%d \t "
			values = append(values, statuses[code])
		}
	}
	fmt.Fprintln(w, headers)
	fmt.Fprintf(w, fmt.Sprintf("%s\n", line), values...)

	fmt.Fprintln(w, "Statistics")
	fmt.Fprintln(w, " \t Average \t STDDEV \t Max")
	fmt.Fprintln(w, fmt.Sprintf("Req/sec \t %.2f \t %s \t %s", rate, "NA", "NA"))
	fmt.Fprintln(w, fmt.Sprintf("Latency \t %s \t %s \t %s", round(average, 2), round(stddev, 2), round(max, 2)))

	if opts.Latency {
		fmt.Fprintln(w, "Latency Distibution")
		for _, bucket := range qbuckets {
			quantile := quantiles[bucket]
			fmt.Fprintln(w, fmt.Sprintf("%.2f%% \t %s", bucket*100.0, round(quantile, 2)))
		}
	}

	if opts.Histogram {
		fmt.Fprintln(w, "Histogram")
		fmt.Fprintln(w, "Bucket \t Count")
		for _, current := range hbuckets {
			bucket := time.Duration(current)
			fmt.Fprintln(w, fmt.Sprintf("%s \t %d", round(bucket, 2), histgram[bucket]))
			count -= histgram[bucket]
		}
		fmt.Fprintln(w, fmt.Sprintf("Inf+ \t %d", count))
	}

	w.Flush()
}

func round(d time.Duration, digits int) time.Duration {
	var divs = []time.Duration{time.Duration(1), time.Duration(10), time.Duration(100), time.Duration(1000)}
	switch {
	case d > time.Second:
		d = d.Round(time.Second / divs[digits])
	case d > time.Millisecond:
		d = d.Round(time.Millisecond / divs[digits])
	case d > time.Microsecond:
		d = d.Round(time.Microsecond / divs[digits])
	}
	return d
}
