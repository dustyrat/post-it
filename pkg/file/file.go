package file

import (
	"log"
	"os"

	"github.com/DustyRat/post-it/pkg/client"
	"github.com/DustyRat/post-it/pkg/file/csv"
)

func ParseFile(file, method, rawUrl string, batchSize int) [][]*client.Request {
	input, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer input.Close()

	csvFile := csv.Parse(input, "")
	requests := make([]*client.Request, 0)
	for i := range csvFile.Records {
		request, _ := csvFile.Records[i].ToRequest(method, rawUrl)
		requests = append(requests, request)
	}

	var chunks [][]*client.Request
	for batchSize < len(requests) {
		requests, chunks = requests[batchSize:], append(chunks, requests[0:batchSize:batchSize])
	}
	chunks = append(chunks, requests)
	return chunks
}
