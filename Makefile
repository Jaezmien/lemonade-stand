FILENAME = stand
BUILD_DIR = ./build

GOARCH = amd64
VERSION = 1.0.0
COMMIT = $(shell git rev-parse --short HEAD)

LDFLAGS = -ldflags "-X main.BuildVersion=${VERSION} -X main.BuildCommit=${COMMIT}"

all: clean linux windows

test:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

clean:
	rm -rf "${BUILD_DIR}"
	mkdir "${BUILD_DIR}"

linux:
	GOOS=linux GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BUILD_DIR}/${FILENAME}-linux-${GOARCH} .

windows:
	GOOS=windows GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BUILD_DIR}/${FILENAME}-windows-${GOARCH}.exe .
