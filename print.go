package main

import (
	"fmt"
	"io/fs"
	"strconv"

	"github.com/fatih/color"
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
	// https://pkg.go.dev/os

	symFlag := false
	switch m := f.fileMode; {
	case (m&fs.ModeSymlink) != 0 && f.linkPath != "":
		// Symbolic Link
		color.Set(color.FgCyan)
		symFlag = true
	case (m & fs.ModeDir) != 0:
		color.Set(color.BgBlue)
	case (m & 0111) != 0:
		// 0111 = --x--x--x
		color.Set(color.FgGreen)
	default:
		color.Set(color.FgWhite)
	}

	if symFlag {
		// Symbolic Link
		fmt.Printf("%s  →  %s", f.fileName, f.linkPath)
	} else {
		fmt.Printf("%s", f.fileName)
	}
	color.Unset()

	// Put a space between file names
	fmt.Printf("  ")
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
