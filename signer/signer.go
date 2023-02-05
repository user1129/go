package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/cheggaaa/pb/v3"
	"github.com/valyala/fastrand"
)

func ExecutePipeline(jobs ...job) {
	in := make(chan interface{})
	for _, j := range jobs {
		out := make(chan interface{})
		go j(in, out)
		in = out
	}
}

type job func(in, out chan interface{})

func SingleHash(in, out chan interface{}) {
	worker := func(data int) string {
		md5 := DataSignerMd5(strconv.Itoa(data))
		crc32 := DataSignerCrc32(strconv.Itoa(data))
		crc32md5 := DataSignerCrc32(md5)
		return crc32 + "~" + crc32md5
	}

	for input := range in {
		data, ok := input.(int)
		if !ok {
			continue
		}
		out <- worker(data)
	}
}

func MultiHash(in, out chan interface{}) {
	worker := func(data string) string {
		var results []string
		for i := 0; i < 6; i++ {
			results = append(results, DataSignerCrc32(strconv.Itoa(i)+data))
		}
		return strings.Join(results, "")
	}

	for input := range in {
		data, ok := input.(string)
		if !ok {
			continue
		}
		out <- worker(data)
	}
}

func CombineResults(in, out chan interface{}) {
	var results []string
	for input := range in {
		data, ok := input.(string)
		if !ok {
			continue
		}
		results = append(results, data)
	}
	sort.Strings(results)
	out <- strings.Join(results, "_")
}

func DataSignerMd5(data string) string {
	return uuid.NewV4().String()
}

func DataSignerCrc32(data string) string {
	return strconv.FormatUint(uint64(fastrand.Uint32()), 10)
}

func main() {
	inputData := make(chan interface{})
	numJobs := 100

	bar := pb.StartNew(numJobs)
	for i := 0; i < numJobs; i++ {
		inputData <- i
		bar.Increment()
	}
	close(inputData)
	ExecutePipeline(
		job(SingleHash),
		job(MultiHash),
		job(CombineResults),
	)
	bar.Finish()

	fmt.Println("Done.")
}
