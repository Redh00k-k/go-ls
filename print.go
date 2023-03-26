package main

import (
	"fmt"
	"io/fs"
	"os"
	"strconv"
	"time"
)

func printMode(fileinfo fs.FileInfo) {
	// c:char device  d:directory  t:sticky  D:device file  L:symlink
	// exp. Dcrw-------
	fmt.Printf("%-11s  ", fileinfo.Mode())
}

func printOwnerGroup(filePath string, fileinfo fs.FileInfo) {
	owner, group := GetOwnerGroup(filePath, fileinfo)
	fmt.Printf("%s %s  ", owner, group)
}

func printSize(sizeLen int, fileinfo fs.FileInfo) {
	fmt.Printf("%*d  ", sizeLen, fileinfo.Size())
}

func printUpdateDate(up time.Time) {
	fmt.Printf("%s  ", up.Format("Jan _2 15:04:05 2006"))
}

func printFilename(file fs.FileInfo) {
	fmt.Printf("%s  ", file.Name())
}

func DisplayLongFormat(files []string) {
	// Get filesize lengeth for printSize
	var sizeLen = 0
	for _, file := range files {
		fileinfo, _ := os.Stat(file)
		l := len(strconv.FormatInt(fileinfo.Size(), 10))
		if sizeLen < l {
			sizeLen = l
		}
	}

	for _, file := range files {
		fileinfo, _ := os.Stat(file)
		printMode(fileinfo)
		printOwnerGroup(file, fileinfo)
		printSize(sizeLen, fileinfo)
		printUpdateDate(fileinfo.ModTime())
		printFilename(fileinfo)
		fmt.Println()
	}
}

func DisplayShortFormat(files []string) {
	for _, file := range files {
		fileinfo, _ := os.Stat(file)
		printFilename(fileinfo)
	}
	fmt.Println()
}
