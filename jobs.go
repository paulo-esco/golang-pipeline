package main

import (
	"sort"
	"strconv"
	"strings"
	"sync"
)

func SingleHash(in, out chan any) {
	mu := new(sync.Mutex)
	mainWg := new(sync.WaitGroup)
	for inp := range in {
		mainWg.Add(1)
		go func(input any) {
			defer mainWg.Done()
			var crc1, md, crc2 string

			data := strconv.Itoa(input.(int))

			wg := new(sync.WaitGroup)

			wg.Add(2)
			go func() {
				defer wg.Done()
				crc1 = DataSignerCrc32(data)
			}()
			go func() {
				defer wg.Done()
				mu.Lock()
				md = DataSignerMd5(data)
				mu.Unlock()
				crc2 = DataSignerCrc32(md)
			}()
			wg.Wait()

			res := crc1 + "~" + crc2

			out <- res
			// fmt.Println("SingleHash:", res)
		}(inp)
	}
	mainWg.Wait()
}

func MultiHash(in, out chan any) {
	mainWg := new(sync.WaitGroup)
	for inp := range in {
		mainWg.Add(1)
		go func(input any) {
			defer mainWg.Done()
			res := make([]string, 6)
			data := input.(string)

			wg := new(sync.WaitGroup)

			f := func(i int) {
				defer wg.Done()
				res[i] = DataSignerCrc32(
					strconv.Itoa(i) + data,
				)
			}

			for i := range 6 {
				wg.Add(1)
				go f(i)
			}
			wg.Wait()

			out <- strings.Join(res, "")
			// fmt.Println("MultiHash:", strings.Join(res, ""))
		}(inp)
	}
	mainWg.Wait()
}

type ByAlpha []string

func (b ByAlpha) Len() int           { return len(b) }
func (b ByAlpha) Less(i, j int) bool { return b[i] < b[j] }
func (b ByAlpha) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }

func CombineResults(in, out chan any) {
	datas := make([]string, 0)
	for data := range in {
		datas = append(datas, data.(string))
	}

	sort.Sort(ByAlpha(datas))

	out <- strings.Join(datas, "_")
	// fmt.Println("CombineResults:", strings.Join(datas, "_"))
}
