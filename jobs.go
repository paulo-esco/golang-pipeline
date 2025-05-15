package main

import (
	"sort"
	"strconv"
	"strings"
)

func SingleHash(in, out chan interface{}) {
	input, _ := <-in
	data := input.(string)
	out <- DataSignerCrc32(data) + "~" +
		DataSignerMd5(DataSignerMd5(data))
}

func MultiHash(in, out chan interface{}) {
	res := ""
	input, _ := <-in
	data := input.(string)
	for i := range 6 {
		res += DataSignerCrc32(strconv.Itoa(i) + data)
	}
	out <- res
}

type ByAlpha []string

func (b ByAlpha) Len() int           { return len(b) }
func (b ByAlpha) Less(i, j int) bool { return b[i] < b[j] }
func (b ByAlpha) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }

func CombineResults(in, out chan interface{}) {
	datas := make([]string, 0)
	for data := range in {
		datas = append(datas, data.(string))
	}

	sort.Sort(ByAlpha(datas))

	out <- strings.Join(datas, "_")
}
