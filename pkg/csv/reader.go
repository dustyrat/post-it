package csv

import (
	"encoding/csv"
	"io"
	"log"
)

type Csv struct {
	Headers []string
	Records []Record
}

type Record struct {
	Headers []string
	Fields  map[string]string
}

func Parse(reader io.Reader) Csv {
	r := csv.NewReader(reader)
	r.LazyQuotes = true
	r.TrimLeadingSpace = true
	records := make([]Record, 0)
	var headers []string
	for {
		line, err := r.Read()
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
			records = append(records, record)
		}
	}
	return Csv{Headers: headers, Records: records}
}
