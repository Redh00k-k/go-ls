package main

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"time"

	"github.com/pborman/getopt/v2"
)

type fInfo struct {
	fileName   string
	fileMode   fs.FileMode
	linkPath   string
	fileSize   int64
	ownerName  string
	groupName  string
	updateTime time.Time
}

func EnumFilePath(filePaths []string, isAll bool) map[string]map[string]fInfo {
	allFilePath := map[string]map[string]fInfo{}
	if len(filePaths) == 0 {
		// ex. gols
		files, err := filepath.Glob("./*")
		if err != nil {
			panic(err)
		}

		if !isAll {
			files = RemoveHiddenFile(files)
		}

		inf := GatherFileInfo(files)
		allFilePath["./"] = inf
	} else {
		for _, path := range filePaths {
			pattern, isSpecifyFile := GeneratePattern(path)
			files, err := filepath.Glob(pattern)
			if err != nil {
				panic(err)
			}

			if !isAll && !isSpecifyFile {
				// Remove hidden files if no "-a" and no file is specified.
				files = RemoveHiddenFile(files)
			}

			inf := GatherFileInfo(files)

			allFilePath[path] = inf
		}
	}
	return allFilePath
}

func main() {
	var (
		isLongfmt bool
		isAll     bool
	)

	getopt.Flag(&isLongfmt, 'l', "list with long format - show permissions")
	getopt.Flag(&isAll, 'a', "list all files including hidden file starting with '.'")
	getopt.Parse()
	filePaths := getopt.Args()

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
