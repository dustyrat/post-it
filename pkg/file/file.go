package file

import (
	"log"
	"os"

	"github.com/DustyRat/post-it/pkg/client"
	"github.com/DustyRat/post-it/pkg/file/csv"
)

// Data ...
type Data struct {
	Record  csv.Record
	Request *client.Request
}

// ParseFile ...
func ParseFile(file, method, rawUrl, body string) ([]string, []*Data) {
	input, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer input.Close()

	csvFile := csv.Parse(input, body)
	data := make([]*Data, 0)
	for i := range csvFile.Records {
		request, _ := csvFile.Records[i].Request(method, rawUrl)
		data = append(data, &Data{
			Record:  csvFile.Records[i],
			Request: request,
		})
	}
	return csvFile.Headers, data
}
