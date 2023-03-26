//go:build windows
// +build windows

package main

import (
	"fmt"
	"io/fs"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	modadvapi32           = syscall.NewLazyDLL("advapi32.dll")
	procGetSecurityInfo   = modadvapi32.NewProc("GetSecurityInfo")
	procLookupAccountSidW = modadvapi32.NewProc("LookupAccountSidW")

	modkernel32     = syscall.NewLazyDLL("kernel32.dll")
	procCreateFileW = modkernel32.NewProc("CreateFileW")
	procCloseHandle = modkernel32.NewProc("CloseHandle")
)

const (
	// https://learn.microsoft.com/en-us/windows/win32/api/accctrl/ne-accctrl-se_object_type
	SE_FILE_OBJECT = 0x00001

	// https://learn.microsoft.com/ja-jp/openspecs/windows_protocols/ms-dtyp/23e75ca3-98fd-4396-84e5-86cd9d40d343
	OWNER_SECURITY_INFORMATION = 0x00001

	// https://learn.microsoft.com/en-us/windows/win32/api/fileapi/nf-fileapi-createfilea
	// https://learn.microsoft.com/en-us/openspecs/windows_protocols/ms-samr/262970b7-cd4a-41f4-8c4d-5a27f0092aaa
	GENERIC_READ          = 0x80000000
	FILE_SHARE_READ       = 0x00000001
	OPEN_EXISTING         = 0x00000003
	FILE_ATTRIBUTE_NORMAL = 0x00000080
	// https://learn.microsoft.com/en-us/windows/win32/fileio/obtaining-a-handle-to-a-directory
	FILE_FLAG_BACKUP_SEMANTICS = 0x02000000

	// https://pkg.go.dev/math#pkg-constants uint64
	INVALID_HANDLE_VALUE = 18446744073709551615
)

type SE_OBJECT_TYPE uint32
type SECURITY_INFORMATION uint32

// https://learn.microsoft.com/en-us/windows/win32/api/fileapi/nf-fileapi-createfilew
func CreateFileW(
	filePath string,
	access uint32,
	mode uint32,
	sa *uint,
	createmode uint32,
	attrs uint32,
	templatefile *uint,
) (r1 uintptr, err error) {
	// https://go.dev/src/syscall/dll_windows.go
	r1, _, err = procCreateFileW.Call(
		uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(filePath))),
		uintptr(access),
		uintptr(mode),
		uintptr(unsafe.Pointer(sa)),
		uintptr(createmode),
		uintptr(attrs),
		uintptr(unsafe.Pointer(templatefile)),
	)
	if r1 == 0 || r1 == INVALID_HANDLE_VALUE {
		return 0, err
	}

	return
}

// https://learn.microsoft.com/en-us/windows/win32/api/handleapi/nf-handleapi-closehandle
func CloseHandle(
	handle uintptr,
) (r1 uintptr, err error) {
	r1, _, err = procCloseHandle.Call(
		handle,
	)

	return
}

// https://learn.microsoft.com/en-us/windows/win32/api/aclapi/nf-aclapi-getsecurityinfo
func GetSecurityInfo(
	handle uintptr,
	objectType SE_OBJECT_TYPE,
	securityInformation SECURITY_INFORMATION,
	owner **struct{},
	group **struct{},
	dacl **struct{},
	sacl **struct{},
	sd **struct{}, //**SECURITY_DESCRIPTOR,
) (r1 uintptr, err error) {
	r1, _, err = procGetSecurityInfo.Call(
		uintptr(handle),
		uintptr(objectType),
		uintptr(securityInformation),
		uintptr(unsafe.Pointer(owner)),
		uintptr(unsafe.Pointer(group)),
		uintptr(unsafe.Pointer(dacl)),
		uintptr(unsafe.Pointer(sacl)),
		uintptr(unsafe.Pointer(sd)),
	)

	return
}

// https://learn.microsoft.com/en-us/windows/win32/api/winbase/nf-winbase-lookupaccountsidw
func LookupAccountSid(
	systemName *uint16,
	sid *struct{}, // *SID,
	name *uint16,
	nameLen *uint32,
	refdDomainName *uint16,
	refdDomainNameLen *uint32,
	use *uint32,
) (r1 uintptr, err error) {
	r1, _, err = procLookupAccountSidW.Call(
		uintptr(unsafe.Pointer(systemName)),
		uintptr(unsafe.Pointer(sid)),
		uintptr(unsafe.Pointer(name)),
		uintptr(unsafe.Pointer(nameLen)),
		uintptr(unsafe.Pointer(refdDomainName)),
		uintptr(unsafe.Pointer(refdDomainNameLen)),
		uintptr(unsafe.Pointer(use)),
		0, 0)

	return
}

// https://learn.microsoft.com/en-us/windows/win32/secauthz/finding-the-owner-of-a-file-object-in-c--
func GetOwnerGroup(filePath string, fileinfo fs.FileInfo) (n string, dn string) {
	hFile, _ := CreateFileW(
		filePath,
		GENERIC_READ,
		FILE_SHARE_READ,
		nil,
		OPEN_EXISTING,
		FILE_FLAG_BACKUP_SEMANTICS,
		nil,
	)
	if hFile == 0 {
		n, dn = "<UNKNOWN>", "<UNKNOWN>"
		return
	}
	defer CloseHandle(hFile)

	var pSD *struct{} //*SECURITY_DESCRIPTOR
	var pSidOwner *struct{}
	ret, _ := GetSecurityInfo(
		hFile,
		SE_FILE_OBJECT,
		OWNER_SECURITY_INFORMATION,
		&pSidOwner,
		nil,
		nil,
		nil,
		&pSD,
	)

	// If the function fails, the return value is a nonzero
	if ret != 0 {
		n, dn = "<UNKNOWN>", "<UNKNOWN>"
		return
	}

	// First call will fail, but we will get necessary buffer sizes
	var nameLen, domainLen, sidUse uint32
	LookupAccountSid(
		nil,
		pSidOwner,
		nil,
		&nameLen,
		nil,
		&domainLen,
		&sidUse,
	)

	// Allocate memory for size returned by previous call
	name := make([]uint16, nameLen)
	domainName := make([]uint16, domainLen)
	ret, _ = LookupAccountSid(
		nil,
		pSidOwner,
		&name[0],
		&nameLen,
		&domainName[0],
		&domainLen,
		&sidUse,
	)
	// If the function fails, it returns zero.
	if ret == 0 {
		n, dn = "<UNKNOWN>", "<UNKNOWN>"
		return
	}

	n = syscall.UTF16ToString(name)
	dn = syscall.UTF16ToString(domainName)

	return
}

func RemoveHiddenFile(files []string) []string {
	tmp := []string{}
	for _, file := range files {
		pointer, err := syscall.UTF16PtrFromString(file)
		if err != nil {
			fmt.Println(err)
			break
		}
		attributes, err := syscall.GetFileAttributes(pointer)
		if err != nil {
			fmt.Println(err)
			break
		}
		if attributes&syscall.FILE_ATTRIBUTE_HIDDEN != 0 {
			// Hidden files
			continue
		}
		tmp = append(tmp, file)
	}
	return tmp
}
