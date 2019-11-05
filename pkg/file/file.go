package file

import (
	"log"
	"os"

	"github.com/DustyRat/post-it/pkg/client"
	"github.com/DustyRat/post-it/pkg/file/csv"
)

func ParseFile(file, method, rawUrl, body string) []*client.Request {
	input, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer input.Close()

	csvFile := csv.Parse(input, body)
	requests := make([]*client.Request, 0)
	for i := range csvFile.Records {
		request, _ := csvFile.Records[i].ToRequest(method, rawUrl)
		requests = append(requests, request)
	}
	return requests
}
