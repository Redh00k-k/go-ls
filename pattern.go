package main

import (
	"os"
	"path/filepath"
)

func checkDir(path string) (string, bool) {
	fileinfo, err := os.Stat(path)
	if err == nil {
		if fileinfo.IsDir() {
			// If "path" is a directory, add "/*"
			return path + "/*", false
		}
	}
	// cannnot find the file or directory
	return path, true
}

func GeneratePattern(path string) (string, bool) {
	dirname, filename := filepath.Split(path)
	var pattern string
	isSpecifyFile := false
	switch {
	case dirname != "" && filename != "":
		// ex. gols ../text.txt ./dir
		// pattern = ../text.txt, ./dir/
		pattern, isSpecifyFile = checkDir(path)
	case dirname != "" && filename == "":
		// ex. gols ../
		// pattern = ../*
		pattern = dirname + "*"
	case dirname == "" && filename != "":
		// ex. gols test.txt
		// pattern = test.txt
		pattern, isSpecifyFile = checkDir(path)
	default:
		// ex. gols
		// Never come here.
		pattern = "./*"
	}

	return pattern, isSpecifyFile
}
