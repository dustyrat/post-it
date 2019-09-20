package controller

import (
	"bufio"
	"encoding/csv"
	"errors"
	"github.com/cheggaaa/pb/v3"
	"github.com/goinggo/work"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"post-it/pkg/client"
	"time"
)

type Controller struct {
	Client *client.Client
	Method string
	Url    string
	headers []string

	BatchSize  int
	ThreadPool int

	//RecordBody bool

	Stats *Stats
}

func (c *Controller) resetStats() {
	c.Stats = &Stats{
		Responses: make(map[int]int, 0),
		Entries: make([]Entry, 0),
	}
}

func (c *Controller) Run(input, output string) error {
	c.resetStats()
	defer func(){
		c.Stats.Print()
	}()

	from := 1
	to := c.BatchSize
	batch := 1

	inputfile, err := os.Open(input)
	if err != nil {
		log.Fatal(err)
	}
	defer inputfile.Close()

	reader := bufio.NewReader(inputfile)
	headers, records := parse(reader)
	headers = append(headers, "status")
	headers = append(headers, "error")
	//if c.RecordBody {
	//	headers = append(headers, "body")
	//}

	writer, err := NewWriter(output)
	if err != nil {
		log.Fatal(err)
	}
	writer.Write(headers)
	writer.Flush()

	total := len(records)
	progress := pb.Full.Start(total)
	defer progress.Finish()

	wp, err := work.New(c.ThreadPool, time.Hour*24, func(message string){})
	if err != nil {
		return errors.New("error creating worker pools")
	}

	var chunks [][]Record
	for c.BatchSize < len(records) {
		records, chunks = records[c.BatchSize:], append(chunks, records[0:c.BatchSize:c.BatchSize])
	}
	chunks = append(chunks, records)

	for i := range chunks {
		w := worker{
			writer: writer,
			client: c.Client,
			method: c.Method,
			url: c.Url,
			chunk: chunks[i],
			batch: batch,
			from: from,
			to: to,
			//recordBody: c.RecordBody,
			stats: c.Stats,
			progress: progress,
		}
		wp.Run(&w)

		from = from + c.BatchSize
		to = to + c.BatchSize
		batch++
	}

	// wait for all worker threads to finish doing their work
	wp.Shutdown()
	return nil
}

type Record struct {
	Headers []string
	Fields map[string]string
}

func parse(reader io.Reader) ([]string, []Record) {
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
			record := Record{ Headers: headers, Fields: make(map[string]string) }
			for i := range headers {
				record.Fields[headers[i]] = line[i]
			}
			records = append(records, record)
		}
	}
	return headers, records
}
