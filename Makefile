BINARY_NAME=go-ls

build:
ifeq ($(OS), Windows_NT)
	go build -o ${BINARY_NAME}.exe -ldflags="-s -w" -trimpath main.go print.go pattern.go gather.go utils_windows.go
	
# 	set GOOS=linux
# Note: "make" cannot set environment variables for shells that invoke "make". 
# 	go build -o ${BINARY_NAME} -ldflags="-s -w" -trimpath main.go print.go pattern.go utils_linux.go
else
	GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME} -ldflags="-s -w" -trimpath main.go print.go pattern.go gather.go utils_linux.go

# 	Install "x86_64-w64-mingw32-gcc" and "x86_64-w64-mingw32-g++"
# 	GOARCH=amd64 GOOS=windows go build -o ${BINARY_NAME} -ldflags="-s -w" -trimpath main.go print.go pattern.go utils_windows.go
endif

run: build
	./${BINARY_NAME}

clean:
	go clean
ifeq ($(OS), Windows_NT)
	del ${BINARY_NAME}
	del ${BINARY_NAME}.exe
else
	rm ${BINARY_NAME}
	rm ${BINARY_NAME}.exe
endif