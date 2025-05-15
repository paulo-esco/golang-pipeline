package main

import (
	"sync"
)

func ExecutePipeline(jobs... job) {
	channels := make([]chan interface{}, len(jobs))

	in := make(chan interface{})
	out := channels[len(channels)-1]

	wg := new(sync.WaitGroup)
	for i, job := range jobs {
		wg.Add(1)
		go func() {
			defer wg.Done()

			jobIn := in
			if i > 0 {
				jobIn = channels[i-1]
			}
			jobOut := channels[i]

			job(jobIn, jobOut)
		}()
	}
	wg.Wait()
}