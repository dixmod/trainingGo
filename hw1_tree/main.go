package main

import (
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"sort"
)

var exceptionNames = [1]string{".DS_Store"}

const (
	headForEnd      = "└───"
	headForNoEnd    = "├───"
	tailForEnd      = "\t"
	tailForNoEnd    = "│\t"
	defaultFileSize = "empty"
)

func main() {
	out := os.Stdout

	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}

	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"

	err := dirTree(out, path, printFiles)

	if err != nil {
		panic(err.Error())
	}
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	err := showTree(out, path, printFiles, "")

	if err != nil {
		panic(err.Error())
	}

	return nil
}

func showTree(out io.Writer, path string, printFiles bool, tail string) error {
	files := getListDirAndFiles(path)

	countFiles := len(files)

	for indexFile := 0; indexFile < countFiles; indexFile++ {
		file := files[indexFile]

		if isExceptionName(file.Name()) {
			continue
		}

		isEndLoop := checkEndLoop(indexFile, countFiles)

		if file.IsDir() {
			printDirName(file, tail, isEndLoop)
			showTree(out, path+"/"+file.Name(), printFiles, tail+getTailByIsLast(isEndLoop))
		} else if printFiles {
			printFileName(file, tail, isEndLoop)
		}
	}

	return nil
}

func getListDirAndFiles(path string) []fs.FileInfo {
	files, err := ioutil.ReadDir(path)

	if err != nil {
		panic(err)
	}

	sort.Slice(files, func(currentIndex, nextIndex int) bool {
		return files[currentIndex].Name() < files[nextIndex].Name()
	})

	return files
}

func isExceptionName(fileName string) bool {
	for _, val := range exceptionNames {
		if val == fileName {
			return true
		}
	}

	return false
}

func checkEndLoop(indexFile int, countFiles int) bool {
	return indexFile >= countFiles-1
}

func printDirName(file fs.FileInfo, tail string, isLast bool) {
	fmt.Printf("%s%s%s\n", tail, getHeadByIsLast(isLast), file.Name())

	return
}

func printFileName(file fs.FileInfo, tail string, isLast bool) {
	fileSize := defaultFileSize

	if 0 < file.Size() {
		fileSize = fmt.Sprintf("%vb", file.Size())
	}

	fmt.Printf("%s%s%s (%s)\n", tail, getHeadByIsLast(isLast), file.Name(), fileSize)

	return
}

func getTailByIsLast(isLast bool) string {
	if isLast {
		return tailForEnd
	}

	return tailForNoEnd
}

func getHeadByIsLast(isLast bool) string {
	if isLast {
		return headForEnd
	}

	return headForNoEnd
}
