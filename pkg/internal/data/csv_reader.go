package data

import (
	"bufio"
	"compress/gzip"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// readCSV reads CSV data from a URL and processes each row
func ReadCSV(url string, processor func([]string, map[string]int) (bool, error)) error {
	// Use a client with timeout to prevent hanging on slow responses
	client := &http.Client{
		Timeout: 1 * time.Minute,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	// Add compression support
	req.Header.Set("Accept-Encoding", "gzip")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Set up the reader based on response content encoding
	var bodyReader io.Reader = resp.Body

	// Check if the response is gzip-encoded
	if strings.Contains(resp.Header.Get("Content-Encoding"), "gzip") {
		gzipReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return fmt.Errorf("error creating gzip reader: %w", err)
		}
		defer gzipReader.Close()
		bodyReader = gzipReader
	}

	// Use buffered reader for better performance
	bufferedReader := bufio.NewReaderSize(bodyReader, 32*1024) // 32KB buffer

	reader := csv.NewReader(bufferedReader)
	reader.Comma = ';'

	// Performance optimizations
	reader.ReuseRecord = true // Reuse the same slice for each Read call
	reader.LazyQuotes = true  // Allow lazy quotes for faster parsing

	// Read headers
	headers, err := reader.Read()
	if err != nil {
		return err
	}

	// Create header map with pre-allocated capacity
	headerMap := make(map[string]int, len(headers))
	for i, header := range headers {
		headerMap[header] = i
	}

	// Process records one by one
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		res, err := processor(record, headerMap)
		if err != nil {
			return err
		}
		if res == false {
			break
		}
	}

	return nil
}
