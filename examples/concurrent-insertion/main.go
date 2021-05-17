package main

import (
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	"github.com/nakabonne/tstorage"
)

func main() {
	tmpDir, err := ioutil.TempDir("", "tstorage")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	storage, err := tstorage.NewStorage(
		tstorage.WithDataPath(tmpDir),
		tstorage.WithPartitionDuration(5*time.Hour),
	)
	if err != nil {
		log.Fatal(err)
	}
	var wg sync.WaitGroup
	for i := int64(1600000); i < 1610000; i++ {
		wg.Add(1)
		go func(timestamp int64) {
			err = storage.InsertRows([]tstorage.Row{
				{
					Metric:    "metric1",
					DataPoint: tstorage.DataPoint{Timestamp: timestamp},
				},
			})
			if err != nil {
				log.Fatal(err)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()

	iterator, _, err := storage.SelectRows("metric1", nil, 1600000, 1610001)
	if err != nil {
		log.Fatal(err)
	}
	for iterator.Next() {
		log.Printf("timestamp: %v, value: %v\n", iterator.DataPoint().Timestamp, iterator.DataPoint().Value)
	}
}