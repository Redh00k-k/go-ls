package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"golang.org/x/exp/slices"
)

func (f *fInfo) printMode() {
	// c:char device  d:directory  t:sticky  D:device file  L:symlink
	// exp. Dcrw-------
	fmt.Printf("%-11s  ", f.fileMode)
}

func (f *fInfo) printOwnerGroup(filePath string, ownerLen int, groupLen int) {
	fmt.Printf("%-*s %-*s  ", ownerLen, f.ownerName, groupLen, f.groupName)
}

func (f *fInfo) printSize(sizeLen int) {
	fmt.Printf("%*d  ", sizeLen, f.fileSize)
}

func (f *fInfo) printUpdateDate() {
	fmt.Printf("%s  ", f.updateTime.Format("Jan _2 15:04:05 2006"))
}

func (f *fInfo) printFilename() {
	// https://pkg.go.dev/os

	// For Windows
	// https://pkg.go.dev/strings#Split
	execExt := strings.FieldsFunc(strings.ToLower(os.Getenv("PATHEXT")), func(r rune) bool { return r == ';' })
	ext := filepath.Ext(f.fileName)

	symFlag := false
	switch m := f.fileMode; {
	case (m&fs.ModeSymlink) != 0 && f.linkPath != "":
		// Symbolic Link
		color.Set(color.FgCyan)
		symFlag = true
	case (m & fs.ModeDir) != 0:
		color.Set(color.BgBlue)
	case (m&0111) != 0 || slices.Contains(execExt, ext) != false:
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

func DisplayLongFormat(files []fInfo) {
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

	for _, file := range files {
		file.printMode()
		file.printOwnerGroup(file.filePath, ownerLen, gourpLen)
		file.printSize(filesizeLen)
		file.printUpdateDate()
		file.printFilename()
		fmt.Println()
	}
}

func DisplayShortFormat(files []fInfo) {
	for _, file := range files {
		file.printFilename()
	}
	fmt.Println()
}
