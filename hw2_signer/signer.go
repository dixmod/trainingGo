package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
)

func main() {
	inputData := []int{1, 2}

	hashSignJobs := []job{
		job(func(in, out chan interface{}) {
			for _, fibNum := range inputData {
				out <- fibNum
			}
		}),
		job(SingleHash),
		job(MultiHash),
		job(CombineResults),
		job(func(in, out chan interface{}) {
			_ = <-in
		}),
	}

	ExecutePipeline(hashSignJobs...)
}

func ExecutePipeline(jobs ...job) {
	wg := &sync.WaitGroup{}
	in := make(chan interface{})

	for _, job := range jobs {
		wg.Add(1)

		out := make(chan interface{})
		go WorkerJob(job, in, out, wg)
		in = out
	}

	wg.Wait()
}

func WorkerJob(job job, in, out chan interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(out)

	job(in, out)
}

func SingleHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	md5QueueChanel := make(chan interface{}, 1)

	for itemData := range in {
		WorkerSingleHash(itemData, out, wg, md5QueueChanel)
	}

	wg.Wait()
}

func WorkerSingleHash(in interface{}, out chan interface{}, wg *sync.WaitGroup, md5QueueChanel chan interface{}) {
	data := strconv.Itoa(in.(int))
	crc32Md5Chanel := make(chan string, 1)

	fmt.Printf("%s SingleHash data %s\n", data, data)

	wg.Add(1)
	go func(data string, out chan<- string, queueChanel chan interface{}) {
		defer wg.Done()
		queueChanel <- data

		md5Res := DataSignerMd5(data)
		fmt.Printf("%s SingleHash md5(data) %s\n", <-queueChanel, md5Res)

		out <- DataSignerCrc32(md5Res)
	}(data, crc32Md5Chanel, md5QueueChanel)

	wg.Add(1)
	go func(data string, crc32Md5Chanel <-chan string, out chan<- interface{}) {
		defer wg.Done()

		crc32DataRes := DataSignerCrc32(data)
		fmt.Printf("%s SingleHash crc32(data) %s\n", data, crc32DataRes)

		result := crc32DataRes + "~" + <-crc32Md5Chanel
		fmt.Printf("%s SingleHash result %s\n", data, result)

		out <- result
	}(data, crc32Md5Chanel, out)
}

func MultiHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}

	for itemData := range in {
		WorkerMultiHash(itemData, out, wg)
	}

	wg.Wait()
}

func WorkerMultiHash(in interface{}, out chan interface{}, wg *sync.WaitGroup) {
	data := in.(string)
	wgDataSignerCrc32 := &sync.WaitGroup{}

	ths := [6]int{0, 1, 2, 3, 4, 5}
	const countThs = len(ths)
	results := make([]string, countThs)

	for index, th := range ths {
		th := strconv.Itoa(th)

		wgDataSignerCrc32.Add(1)

		go func(results []string, data string, index int, wgDataSignerCrc32 *sync.WaitGroup) {
			defer wgDataSignerCrc32.Done()

			results[index] = DataSignerCrc32(th + data)

			fmt.Printf("%s MultiHash: crc32(th+step1)) %s %s\n", data, th, results[index])
		}(results, data, index, wgDataSignerCrc32)
	}

	wg.Add(1)

	go func(out chan<- interface{}) {
		defer wg.Done()

		wgDataSignerCrc32.Wait()

		result := strings.Join(results[:], "")

		fmt.Printf("%s MultiHash result: %s\n", data, result)

		out <- result
	}(out)
}

func CombineResults(in, out chan interface{}) {
	var data []string

	for i := range in {
		data = append(data, i.(string))
	}

	sort.Strings(data)
	combineResults := strings.Join(data[:], "_")

	fmt.Printf("CombineResults %s\n", combineResults)

	out <- combineResults
}
