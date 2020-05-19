package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"sort"
)

type SourceInfo struct {
	Extension string
	Lines     int
	Size      int64
}

/*
  TODO
    - order by lines, name, size (asc desc)
    - coloring
    - intelligent ext filter (all mode)
    - skip binary files
*/

var result map[string]SourceInfo = make(map[string]SourceInfo)

func main() {
	err := filepath.Walk(".", walker)
	if err != nil {
		panic(err)
	}

	values := make([]SourceInfo, 0, len(result))
	for _, info := range result {
		values = append(values, info)
	}

	sort.Slice(values, func(i, j int) bool { return values[i].Lines > values[j].Lines })

	fmt.Println("Extension     Lines       Size")

	for i := range values {
		info := values[i]
		fmt.Printf("%10s   %6d   %8d\n", info.Extension, info.Lines, info.Size)
	}
}

func walker(filename string, info os.FileInfo, err error) error {
	if err != nil || info.IsDir() {
		return nil
	}

	ext := path.Ext(filename)
	if ext == "" {
		return nil
	}

	var lines int
	lines, err = lineCount(filename)

	if err != nil {
		return err
	}

	size := info.Size()

	source, ok := result[ext]
	if !ok {
		source = SourceInfo{Extension: ext}
	}

	source.Lines += lines
	source.Size += size

	result[ext] = source

	return nil
}

func lineCount(file string) (int, error) {
	f, err := os.Open(file)
	if err != nil {
		return 0, err
	}

	defer f.Close()

	r := bufio.NewReader(f)

	var lines int

	for {
		_, err := r.ReadString('\n')
		if err == io.EOF {
			return lines, nil
		} else if err != nil {
			return 0, err
		}

		lines++
	}
}
