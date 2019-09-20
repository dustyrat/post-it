package controller

import (
	"errors"
	"fmt"
	"github.com/cheggaaa/pb/v3"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"post-it/pkg/client"
	"runtime/debug"
	"strings"
)

type worker struct {
	id     int
	client *client.Client
	method string
	url    string
	chunk  []Record
	writer *Writer

	to    int
	from  int
	batch int

	//recordBody bool
	stats    *Stats
	progress *pb.ProgressBar
}

func (w *worker) Work(id int) {
	w.id = id
	defer w.done()

	for i := range w.chunk {
		w.call(w.chunk[i])
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

func (w *worker) call(record Record) {
	defer w.progress.Increment()
	defer w.writer.Flush()

	entry := Entry{Record: record}
	defer func(){
		if !client.IsSuccessful(entry.Status) {
			w.writer.Write(entry.Strings())
		}
	}()

	rawurl := w.url
	for k, v := range record.Fields {
		rawurl = strings.Replace(rawurl, fmt.Sprintf("{{%s}}", k), v, 1)
	}

	_url, err := url.Parse(rawurl)
	if err != nil {
		entry.Error = err
		return
	}

	var response *client.Response
	if method, ok := record.Fields["method"]; ok {
		response, err = w.client.Do(method, _url, http.Header{}, nil)
	} else {
		response, err = w.client.Do(w.method, _url, http.Header{}, nil)
	}
	if err != nil {
		entry.Error = err
		return
	}

	entry.Status = response.StatusCode
	w.stats.Increment(response.StatusCode)

	if !client.IsSuccessful(response.StatusCode) {
		entry.Error = errors.New(string(response.Body))
	}
}
