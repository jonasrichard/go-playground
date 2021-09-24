package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"sort"

	"github.com/lucasepe/color"
)

type SourceInfo struct {
	Extension string
	Lines     int
	Size      int64
}

type GroupInfo struct {
	Name  string
	Lines int
	Size  int64
}

/*
  TODO
    - order by lines, name, size (asc desc)
    - skip binary files
*/

var result map[string]SourceInfo = make(map[string]SourceInfo)

func main() {
	filetypes := extensions()

	err := filepath.Walk(".", walker)
	if err != nil {
		panic(err)
	}

	groups := make(map[string]GroupInfo)

	for _, info := range result {
		ft, ok := filetypes[info.Extension]
		if !ok {
			ft = "Unknown"
		}

		group, ok := groups[ft]
		if !ok {
			group = GroupInfo{Name: ft}
		}

		group.Lines += info.Lines
		group.Size += info.Size

		groups[ft] = group
	}

	values := make([]GroupInfo, 0)
	for _, info := range groups {
		values = append(values, info)
	}

	sort.Slice(values, func(i, j int) bool { return values[i].Lines > values[j].Lines })

	color.Cyan("Extension     Lines       Size\n")

	for i := range values {
		info := values[i]
		msg := fmt.Sprintf("%10s   %6d   %8d\n", info.Name, info.Lines, info.Size)

		if info.Name == "Unknown" {
			color.Red(msg)
		} else {
			color.Green(msg)
		}
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

func extensions() map[string]string {
	var filetypes = []struct {
		name string
		exts []string
	}{
		{"C", []string{".c", ".cc", ".h"}},
		{"Config", []string{".conf", ".yml", ".yaml"}},
		{"Elixir", []string{".ex", ".exs"}},
		{"Erlang", []string{".erl", ".hrl"}},
		{"Go", []string{".go"}},
		{"HTML", []string{".html", ".css"}},
		{"Javascript", []string{".js"}},
		{"JSON", []string{".json"}},
		{"Ruby", []string{".rb", ".gemspec", ".rake", ".ru"}},
		{"Rust", []string{".rs"}},
		{"SQL", []string{".sql"}},
	}

	var result = make(map[string]string)

	for _, filetype := range filetypes {
		for _, ext := range filetype.exts {
			result[ext] = filetype.name
		}
	}

	return result
}
