package main

import (
	"flag"
	"fmt"
	"path/filepath"
)

func EnumFilePath(filePaths []string, isAll bool) map[string][]string {
	allFilePath := map[string][]string{}
	if len(filePaths) == 0 {
		// ex. gols
		files, err := filepath.Glob("./*")
		if err != nil {
			panic(err)
		}

		if !isAll {
			files = RemoveHiddenFile(files)
		}
		allFilePath["./"] = files
	} else {
		for _, path := range filePaths {
			pattern, isSpecifyFile := GeneratePattern(path)
			if pattern == "" {
				allFilePath = map[string][]string{}
			}
			files, err := filepath.Glob(pattern)
			if err != nil {
				panic(err)
			}

			if !isAll && !isSpecifyFile {
				// Remove hidden files if no "-a" and no file is specified.
				files = RemoveHiddenFile(files)
			}
			allFilePath[path] = files
		}
	}
	return allFilePath
}

func main() {
	var isLongfmt, isAll bool
	flag.BoolVar(&isLongfmt, "l", false, "list with long format - show permissions")
	flag.BoolVar(&isAll, "a", false, "list all files including hidden file starting with '.'")
	flag.Parse()

	filePaths := flag.Args()
	allFilePath := EnumFilePath(filePaths, isAll)

	for path, files := range allFilePath {
		fmt.Println(path, ":")
		fmt.Println("Total files: ", len(files))

		if !isLongfmt {
			DisplayShortFormat(files)
		} else {
			DisplayLongFormat(files)
		}
		fmt.Println()
	}
}
