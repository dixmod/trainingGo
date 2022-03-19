package main

import (
	"fmt"
	"strconv"
	"strings"
)

func ExecutePipeline(dataInput ...string) {
	var results [...]string

	for index, data := range dataInput {
		singleHash := SingleHash(index, data)
		results[index] = MultiHash(singleHash)

		fmt.Println("")
	}

	CombineResults(results)
}

func main() {
	ExecutePipeline("0", "1")
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
	var results [6]string

	for index, th := range []int{0, 1, 2, 3, 4, 5} {
		th := strconv.Itoa(th)
		results[index] = DataSignerCrc32(th + data)

		fmt.Printf("%s MultiHash: crc32(th+step1)) %s %s\n", data, th, results[index])
	}

	result := strings.Join(results[:], "")

	fmt.Printf("%s MultiHash result: %s\n", data, result)

	return result
}

func CombineResults(data []string) string {
	return strings.Join(data[:], "_")
}
