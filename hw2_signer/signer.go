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
	mu := &sync.Mutex{}
	wg := &sync.WaitGroup{}

	for i := range in {
		wg.Add(1)

		go WorkerSingleHash(i, out, wg, mu)
	}

	wg.Wait()
}

func WorkerSingleHash(in interface{}, out chan interface{}, wg *sync.WaitGroup, mu *sync.Mutex) {
	defer wg.Done()
	buf := make(chan string)

	data := strconv.Itoa(in.(int))

	fmt.Printf("%s SingleHash data %s\n", data, data)

	mu.Lock()
	go WorkerMd5(data, buf)
	md5res := <-buf
	fmt.Printf("%s SingleHash md5(data) %s\n", data, md5res)
	mu.Unlock()

	go WorkerCrc32(md5res, buf)
	crc32res := <-buf
	fmt.Printf("%s SingleHash crc32(md5(data)) %s\n", data, crc32res)

	go WorkerCrc32(data, buf)
	crc32datares := <-buf
	fmt.Printf("%s SingleHash crc32(data) %s\n", data, crc32datares)

	result := crc32datares + "~" + crc32res
	fmt.Printf("%s SingleHash result %s\n", data, result)

	out <- result
}

func WorkerMd5(data string, out chan string) {
	out <- DataSignerMd5(data)
}

func WorkerCrc32(data string, out chan string) {
	out <- DataSignerCrc32(data)
}

func MultiHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}

	for itemData := range in {
		wg.Add(1)

		go WorkerMultiHash(itemData.(string), out, wg)
	}

	wg.Wait()
}

func WorkerMultiHash(data string, out chan interface{}, wg *sync.WaitGroup) {
	defer wg.Done()

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

	wgDataSignerCrc32.Wait()

	result := strings.Join(results[:], "")

	fmt.Printf("%s MultiHash result: %s\n", data, result)

	out <- result
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
