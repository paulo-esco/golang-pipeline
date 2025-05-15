package main

import (
	"sync"
)

func ExecutePipeline(jobs ...job) {
	var channels []chan any

	in := make(chan any)
	// out := channels[len(channels)-1]

	wg := new(sync.WaitGroup)
	for i, job := range jobs {
		channels = append(channels, make(chan any))
		wg.Add(1)
		go func() {
			defer wg.Done()

			jobIn := in
			if i > 0 {
				jobIn = channels[i-1]
			}
			jobOut := channels[i]

			job(jobIn, jobOut)
			close(jobOut)
		}()
	}
	wg.Wait()
}
