package main

import (
	"io/fs"
	"os"
	"path/filepath"
)

func GatherFileInfo(files []string) []fInfo {
	inf := []fInfo{}
	for _, fPath := range files {
		fileinfo, _ := os.Lstat(fPath)

		_, filename := filepath.Split(fPath)

		var fi fInfo
		fi.fileName = filename
		fi.filePath = fPath
		fi.fileMode = fileinfo.Mode()
		if fi.fileMode&fs.ModeSymlink != 0 {
			fi.linkPath, _ = os.Readlink(fPath)
		}

		fi.fileSize = fileinfo.Size()
		owner, group := GetOwnerGroup(fPath, fileinfo)
		fi.ownerName = owner
		fi.groupName = group
		fi.updateTime = fileinfo.ModTime()

		inf = append(inf, fi)
	}

	return inf
}
