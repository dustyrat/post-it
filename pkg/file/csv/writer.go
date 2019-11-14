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
	w.writer.Write(row)
	w.mutex.Unlock()
}

// WriteAll ...
func (w *Writer) WriteAll(rows [][]string) {
	w.mutex.Lock()
	w.writer.WriteAll(rows)
	w.mutex.Unlock()
}

// Flush ...
func (w *Writer) Flush() {
	w.mutex.Lock()
	w.writer.Flush()
	w.mutex.Unlock()
}
