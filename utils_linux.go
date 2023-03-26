//go:build (linux && 386) || (darwin && !cgo)
// +build linux,386 darwin,!cgo

package main

import (
	"io/fs"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

func GetOwnerGroup(filePath string, fileinfo fs.FileInfo) (owner string, group string) {
	var uid, gid int
	if stat, ok := fileinfo.Sys().(*syscall.Stat_t); ok {
		uid = int(stat.Uid)
		gid = int(stat.Gid)
	} else {
		uid = os.Getuid()
		gid = os.Getgid()
	}

	u, err := user.LookupId(strconv.Itoa(uid))
	if err != nil {
		owner = strconv.Itoa(uid)
	} else {
		owner = u.Username
	}

	g, err := user.LookupGroupId(strconv.Itoa(gid))
	if err != nil {
		group = strconv.Itoa(gid)
	} else {
		group = g.Name
	}

	return
}

func RemoveHiddenFile(filePath []string) []string {
	tmp := []string{}
	for _, path := range filePath {
		_, filename := filepath.Split(path)
		if strings.HasPrefix(filename, ".") {
			// ".<file name>" is hidden files
			continue
		}
		tmp = append(tmp, path)
	}
	return tmp
}
