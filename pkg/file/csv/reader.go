/*
Copyright © 2019 Dustin Ratcliffe <dustin.k.ratcliffe@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package csv

import (
	"bytes"
	"encoding/csv"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/DustyRat/post-it/pkg/client"

	"github.com/spkg/bom"
)

type Csv struct {
	Headers []string
	Records []Record
}

type Record struct {
	Headers []string
	Body    []byte
	Fields  map[string]string
}

func (r *Record) ToRequest(method, rawurl string) (*client.Request, error) {
	return client.NewRequest(method, rawurl, http.Header{}, bytes.NewBuffer(r.Body), r.Fields)
}

func Parse(file *os.File, body string) Csv {
	if file == nil {
		return Csv{}
	}

	reader := csv.NewReader(bom.NewReader(file))
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true
	records := make([]Record, 0)
	var headers []string
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
			record := Record{Headers: headers, Fields: make(map[string]string)}
			for i := range headers {
				record.Fields[headers[i]] = line[i]
			}
			if b, ok := record.Fields[body]; ok {
				record.Body = []byte(b)
			}
			records = append(records, record)
		}
	}
	return Csv{Headers: headers, Records: records}
}
