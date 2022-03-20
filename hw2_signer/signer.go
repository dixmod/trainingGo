package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

func main() {
	ExecutePipeline("0", "1")
}

func ExecutePipeline(dataInput ...string) {
	var results []string

	for index, data := range dataInput {
		singleHash := SingleHash(index, data)
		result := MultiHash(singleHash)

		results = append(results, result)

		fmt.Println("")
	}

	CombineResults(results)
}

func SingleHash(index int, data string) string {
	fmt.Printf("%d SingleHash data %s\n", index, data)

	md5res := DataSignerMd5(data)
	fmt.Printf("%d SingleHash md5(data) %s\n", index, md5res)

	crc32res := DataSignerCrc32(md5res)
	fmt.Printf("%d SingleHash crc32(md5(data)) %s\n", index, crc32res)

	crc32datares := DataSignerCrc32(data)
	fmt.Printf("%d SingleHash crc32(data) %s\n", index, crc32datares)

	result := crc32datares + "~" + crc32res
	fmt.Printf("%d SingleHash result %s\n", index, result)

	return result
}

func MultiHash(data string) string {
	ths := [6]int{0, 1, 2, 3, 4, 5}
	var results [len(ths)]string

	for index, th := range ths {
		th := strconv.Itoa(th)
		results[index] = DataSignerCrc32(th + data)

		fmt.Printf("%s MultiHash: crc32(th+step1)) %s %s\n", data, th, results[index])
	}

	result := strings.Join(results[:], "")

	fmt.Printf("%s MultiHash result: %s\n", data, result)

	return result
}

func CombineResults(data []string) string {
	sort.Strings(data)
	combineResults := strings.Join(data[:], "_")

	fmt.Printf("CombineResults %s\n", combineResults)

	return combineResults
}
