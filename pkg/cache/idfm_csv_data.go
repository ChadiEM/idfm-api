package cache

import (
	"fmt"
	"idfm/pkg/internal/data"
	"log"
	"sync"
	"time"
)

type LinesCache struct {
	Data       []LineCacheData
	LastUpdate time.Time    // Last update timestamp
	mutex      sync.RWMutex // Lock for concurrent access
}

type LineCacheData struct {
	IDLine        string
	TransportMode string
	ShortNameLine string
	OperatorName  string
}

type StopsCache struct {
	Data       []StopCacheData
	LastUpdate time.Time    // Last update timestamp
	mutex      sync.RWMutex // Lock for concurrent access
}

type StopCacheData struct {
	StopID   string
	StopName string
	LineID   string
}

// Global map of cached CSV data
var (
	lineCache = &LinesCache{}
	stopCache = &StopsCache{}
)

// InitializeCaches preloads all CSV files used by the application
func InitializeCaches() {
	log.Println("Preloading CSV data...")
	if err := updateCSVCache(); err != nil {
		log.Printf("Error preloading CSV from %v", err)
	}
	log.Println("CSV data preloading complete")

	// Start background refresh
	go startPeriodicCSVRefresh()
}

// updateCSVCache fetches and updates the cache for a specific URL
func updateCSVCache() error {
	err := updateLineData()
	if err != nil {
		return err
	}
	err = updateStopData()
	if err != nil {
		return err
	}
	return nil
}

func updateLineData() error {
	log.Printf("Updating CSV cache for lines")

	lineData := make([]LineCacheData, 0)

	err := data.ReadCSV("https://data.iledefrance-mobilites.fr/explore/dataset/referentiel-des-lignes/download/?format=csv&timezone=Europe/Paris&lang=fr&use_labels_for_header=true&csv_separator=%3B",
		func(row []string, headers map[string]int) (bool, error) {
			lineData = append(lineData, LineCacheData{
				IDLine:        row[headers["ID_Line"]],
				TransportMode: row[headers["TransportMode"]],
				ShortNameLine: row[headers["ShortName_Line"]],
				OperatorName:  row[headers["OperatorName"]]})

			return true, nil
		})

	if err != nil {
		return fmt.Errorf("error fetching CSV: %w", err)
	}

	// Update the cache
	lineCache.mutex.Lock()
	defer lineCache.mutex.Unlock()

	lineCache.Data = lineData
	lineCache.LastUpdate = time.Now()

	log.Printf("CSV cache updated (%d rows)", len(lineData))
	return nil
}

func updateStopData() error {
	log.Printf("Updating CSV cache for stops")

	stopData := make([]StopCacheData, 0)

	err := data.ReadCSV("https://data.iledefrance-mobilites.fr/explore/dataset/arrets-lignes/download/?format=csv&timezone=Europe/Berlin&lang=fr&use_labels_for_header=true&csv_separator=%3B",
		func(row []string, headers map[string]int) (bool, error) {
			stopData = append(stopData, StopCacheData{
				StopID:   row[headers["stop_id"]],
				StopName: row[headers["stop_name"]],
				LineID:   row[headers["route_id"]],
			})
			return true, nil
		})

	if err != nil {
		return fmt.Errorf("error fetching CSV: %w", err)
	}

	// Update the cache
	stopCache.mutex.Lock()
	defer stopCache.mutex.Unlock()

	stopCache.Data = stopData
	stopCache.LastUpdate = time.Now()

	log.Printf("CSV cache updated (%d rows)", len(stopData))
	return nil
}

// startPeriodicCSVRefresh starts a goroutine to refresh CSV caches periodically
func startPeriodicCSVRefresh() {
	interval := 3 * time.Hour // Update every 3 hours
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		<-ticker.C
		log.Println("Starting periodic CSV cache refresh...")

		if err := updateCSVCache(); err != nil {
			log.Printf("Error refreshing CSV cache %s", err)
		}

		log.Println("Periodic CSV cache refresh complete")
	}
}

// processCacheData processes data from a CSV cache using a generic type
func processCacheData[T LineCacheData | StopCacheData](cache []T, mutex *sync.RWMutex, processor func(data T) (bool, error)) error {
	mutex.RLock()
	defer mutex.RUnlock()

	// Process cached data
	for _, record := range cache {
		continueProcessing, err := processor(record)
		if err != nil {
			return err
		}
		if !continueProcessing {
			break
		}
	}

	return nil
}

// ProcessLineCache processes data from the CSV cache
func ProcessLineCache(processor func(data LineCacheData) (bool, error)) error {
	return processCacheData(lineCache.Data, &lineCache.mutex, processor)
}

// ProcessStopCache processes data from the CSV cache
func ProcessStopCache(processor func(data StopCacheData) (bool, error)) error {
	return processCacheData(stopCache.Data, &stopCache.mutex, processor)
}
