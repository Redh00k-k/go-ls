package main

import (
	"os"
	"path/filepath"
)

func GatherFileInfo(files []string) map[string]fInfo {
	inf := map[string]fInfo{}
	for _, fPath := range files {
		fileinfo, _ := os.Stat(fPath)

		_, filename := filepath.Split(fPath)

		var fi fInfo
		fi.fileName = filename
		fi.fileMode = fileinfo.Mode()
		fi.fileMode = fileinfo.Mode()
		fi.fileSize = fileinfo.Size()
		owner, group := GetOwnerGroup(fPath, fileinfo)
		fi.ownerName = owner
		fi.groupName = group
		fi.updateTime = fileinfo.ModTime()

		inf[fPath] = fi
	}
	return inf
}
