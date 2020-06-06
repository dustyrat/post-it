package csv

import (
	"encoding/csv"
	"os"
	"sync"
)

// Writer ...
type Writer struct {
	mutex  *sync.Mutex
	writer *csv.Writer
}

// NewWriter ...
func NewWriter(fileName string) (*Writer, error) {
	csvFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0777)
	if err != nil {
		return nil, err
	}
	w := csv.NewWriter(csvFile)
	return &Writer{writer: w, mutex: &sync.Mutex{}}, nil
}

// Write ...
func (w *Writer) Write(row []string) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.writer.Write(row)
}

// WriteAll ...
func (w *Writer) WriteAll(rows [][]string) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.writer.WriteAll(rows)
}

// Flush ...
func (w *Writer) Flush() {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.writer.Flush()
}
