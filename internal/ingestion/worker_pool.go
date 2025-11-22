package ingestion

import (
	"log"
	"time"
)

type Row = []interface{}

func StartWorkers(workerCount int, batchSize int) (chan Row, func()) {
	rowsCh := make(chan Row, 10000)
	stop := make(chan struct{})

	for i := 0; i < workerCount; i++ {
		go func(id int) {
			batch := make([]Row, 0, batchSize)
			ticker := time.NewTicker(2 * time.Second)
			defer ticker.Stop()
			for {
				select {
				case r, ok := <-rowsCh:
					if !ok {
						if len(batch) > 0 {
							if err := BatchInsert(batch); err != nil {
								log.Println("worker flush failed:", err)
							}
						}
						return
					}
					batch = append(batch, r)
					if len(batch) >= batchSize {
						if err := BatchInsert(batch); err != nil {
							log.Println("worker insert failed:", err)
						}
						batch = batch[:0]
					}
				case <-ticker.C:
					if len(batch) > 0 {
						if err := BatchInsert(batch); err != nil {
							log.Println("worker periodic insert failed:", err)
						}
						batch = batch[:0]
					}
				case <-stop:
					return
				}
			}
		}(i)
	}

	cancelFn := func() {
		close(rowsCh)
		close(stop)
	}

	return rowsCh, cancelFn
}
