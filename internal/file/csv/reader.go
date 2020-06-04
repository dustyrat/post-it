package csv

import (
	"bytes"
	"encoding/csv"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"sync"

	internal "github.com/DustyRat/post-it/internal/http"
)

// Record ...
type Record struct {
	Headers []string
	Body    []byte
	Fields  map[string]string

	Request *internal.Request
}

// Reader ...
type Reader struct {
	file   *os.File
	mutex  *sync.Mutex
	reader *csv.Reader

	count   int
	headers []string
	body    string

	method string
	rawurl string
}

// NewReader ...
func NewReader(file *os.File, method, rawurl, body string) *Reader {
	if file == nil {
		log.Fatal(errors.New("no file provided"))
	}

	reader := csv.NewReader(file)
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true

	var headers []string
	var count int
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		if headers == nil {
			headers = line
		} else {
			count++
		}
	}
	_, err := file.Seek(0, io.SeekStart)
	if err != nil {
		log.Fatal(err)
	}

	_, err = reader.Read()
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}

	return &Reader{
		file:   file,
		reader: reader,
		mutex:  &sync.Mutex{},

		count:   count,
		headers: headers,

		body:   body,
		method: method,
		rawurl: rawurl,
	}
}

// Read ...
func (r *Reader) Read() *Record {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	line, err := r.reader.Read()
	if err == io.EOF {
		return nil
	} else if err != nil {
		log.Fatal(err)
	}

	record := Record{Headers: r.headers, Fields: make(map[string]string)}
	for i := range r.headers {
		record.Fields[r.headers[i]] = line[i]
	}
	if b, ok := record.Fields[r.body]; ok {
		record.Body = []byte(b)
	}
	request, _ := internal.NewRequest(r.method, r.rawurl, http.Header{}, bytes.NewBuffer(record.Body), record.Fields)
	record.Request = request
	return &record
}

// Headers ...
func (r Reader) Headers() []string {
	return r.headers
}

// Count ...
func (r Reader) Count() int {
	return r.count
}
