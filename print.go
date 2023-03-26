package main

import (
	"fmt"
	"strconv"
)

func printMode(f fInfo) {
	// c:char device  d:directory  t:sticky  D:device file  L:symlink
	// exp. Dcrw-------
	fmt.Printf("%-11s  ", f.fileMode)
}

func printOwnerGroup(filePath string, f fInfo, ownerLen int, groupLen int) {
	fmt.Printf("%-*s %-*s  ", ownerLen, f.ownerName, groupLen, f.groupName)
}

func printSize(sizeLen int, f fInfo) {
	fmt.Printf("%*d  ", sizeLen, f.fileSize)
}

func printUpdateDate(f fInfo) {
	fmt.Printf("%s  ", f.updateTime.Format("Jan _2 15:04:05 2006"))
}

func printFilename(f fInfo) {
	fmt.Printf("%s  ", f.fileName)
}

func DisplayLongFormat(files map[string]fInfo) {
	// Get filesize lengeth for printSize
	var filesizeLen, ownerLen, gourpLen = 0, 0, 0
	for _, f := range files {
		l := len(strconv.FormatInt(f.fileSize, 10))
		if filesizeLen < l {
			filesizeLen = l
		}

		l = len(f.ownerName)
		if ownerLen < l {
			ownerLen = l
		}

		l = len(f.groupName)
		if gourpLen < l {
			gourpLen = l
		}
	}

	for path, file := range files {
		printMode(file)
		printOwnerGroup(path, file, ownerLen, gourpLen)
		printSize(filesizeLen, file)
		printUpdateDate(file)
		printFilename(file)
		fmt.Println()
	}
}

func DisplayShortFormat(files map[string]fInfo) {
	for _, file := range files {
		printFilename(file)
	}
	fmt.Println()
}
