package worker

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"runtime/debug"
	"sort"
	"strconv"
	"strings"

	"github.com/DustyRat/post-it/pkg/options"

	"github.com/DustyRat/post-it/pkg/file/csv"

	"github.com/DustyRat/post-it/pkg/client"

	"github.com/vbauerster/mpb/v5"

	log "github.com/sirupsen/logrus"
)

var (
	statusDxx *regexp.Regexp
	statusDdd *regexp.Regexp
)

func init() {
	statusDxx = regexp.MustCompile("^-?\\dxx")
	statusDdd = regexp.MustCompile("\\d{3}")
}

type worker struct {
	id      int
	record  csv.Record
	request *client.Request

	pool     *Pool
	progress *mpb.Bar
}

type entry struct {
	record  csv.Record
	request client.Request
	err     error
}

// Strings ...
func (e *entry) Strings(flags options.Flags) []string {
	out := make([]string, 0)
	for _, header := range e.record.Headers {
		if value, ok := e.record.Fields[header]; ok {
			out = append(out, value)
		} else {
			out = append(out, "")
		}
	}
	response := e.request.Response
	out = append(out, strconv.Itoa(response.StatusCode))

	if flags.Headers {
		out = append(out, toString(response.Header))
	}

	if flags.Body {
		out = append(out, string(response.Body))
	}

	if e.err != nil {
		out = append(out, e.err.Error())
	} else {
		out = append(out, "")
	}
	return out
}

// Work ...
func (w *worker) Work(id int) {
	w.id = id
	entry := &entry{record: w.record, request: *w.request}
	defer func() {
		defer w.done()
		// w.write(entry)
		go write(w.pool.writer, *w.pool.options, *entry)
		w.progress.Increment()
		w.pool.increment()
	}()

	response, err := w.pool.client.Do(w.request.Method, w.request.URL, w.request.Header, w.request.Body)
	if err != nil {
		entry.err = err
		w.pool.stats.Errors.Increment(err)
	}

	entry.request.Response = response
	w.pool.stats.Codes.Increment(response.StatusCode)
	w.pool.stats.Latencies.Increment(response.Duration)
}

func write(w *csv.Writer, options options.Options, entry entry) {
	if w == nil {
		return
	}
	defer w.Flush()

	request := entry.request
	response := request.Response

	switch options.Flags.Type {
	case "all":
		output := entry.Strings(options.Flags)
		w.Write(output)
	case "error":
		if entry.err != nil {
			output := entry.Strings(options.Flags)
			w.Write(output)
		}
	case "status":
		if options.Flags.Status == "any" {
			output := entry.Strings(options.Flags)
			w.Write(output)
		} else if statusDxx.MatchString(options.Flags.Status) {
			status, _ := strconv.Atoi(strings.Replace(options.Flags.Status, "xx", "00", 1))
			if status > 0 {
				a := status
				b := a + 100
				if client.InRange(response.StatusCode, a, b) {
					output := entry.Strings(options.Flags)
					w.Write(output)
				}
			} else {
				a := -status
				b := a + 100
				if !client.InRange(response.StatusCode, a, b) {
					output := entry.Strings(options.Flags)
					w.Write(output)
				}
			}
		} else if statusDdd.MatchString(options.Flags.Status) {
			status, _ := strconv.Atoi(options.Flags.Status)
			if response.StatusCode == status {
				output := entry.Strings(options.Flags)
				w.Write(output)
			}
		}
	}
}

func (w *worker) done() {
	if r := recover(); r != nil {
		log.Debug("recovered from ", r)
		stack := debug.Stack()
		var err error
		switch t := r.(type) {
		case string:
			err = errors.New(t)
		case error:
			err = t
		default:
			err = errors.New("unknown error")
		}
		log.WithFields(log.Fields{"stacktrace": string(stack)}).Errorf("[%d] %s", w.id, err)
	}
}

func toString(headers http.Header) string {
	keys := make([]string, 0)
	for key := range headers {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	str := ""
	for i, key := range keys {
		str += fmt.Sprintf("%s: ", key)
		values := headers[key]
		for j, value := range values {
			str += value
			if j < len(values)-1 {
				str += ", "
			}
		}

		if i < len(headers)-1 {
			str += "; "
		} else {
			str += ";"
		}
	}
	return str
}
