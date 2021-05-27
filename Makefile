PACKAGE=myTldr
VERSION=0.7

linters:
	golangci-lint run --disable-all \
	    --enable gosimple \
	    --enable errcheck \
	    --enable golint \
	    --enable govet \
	    --enable staticcheck \
	    --enable gocritic \
	    --enable gosec \
	    --enable scopelint \
	    --enable prealloc \
	    --enable maligned \
	    --enable ineffassign \
	    --enable unparam \
	    --enable deadcode \
	    --enable unused \
	    --enable varcheck \
	    --enable unconvert \
	    --enable misspell \
	    --enable goconst \
	    --enable gochecknoinits \
	    --enable gochecknoglobals \
	    --enable nakedret \
	    --enable gocyclo \
	    --enable goimports \
	    --enable gofmt


all: linux-386 linux-amd64 linux-arm linux-arm64 darwin-amd64

linux-386:
	GOOS=linux GOARCH=386 go build -o ${PACKAGE}-${VERSION}-linux-386
	tar -czf ${PACKAGE}-${VERSION}-linux-386.tar.gz ${PACKAGE}-${VERSION}-linux-386

linux-amd64:
	GOOS=linux GOARCH=amd64 go build -o ${PACKAGE}-${VERSION}-linux-amd64
	tar -czf ${PACKAGE}-${VERSION}-linux-amd64.tar.gz ${PACKAGE}-${VERSION}-linux-amd64

linux-arm:
	GOOS=linux GOARCH=arm go build -o ${PACKAGE}-${VERSION}-linux-arm
	tar -czf ${PACKAGE}-${VERSION}-linux-arm.tar.gz ${PACKAGE}-${VERSION}-linux-arm

linux-arm64:
	GOOS=linux GOARCH=arm64 go build -o ${PACKAGE}-${VERSION}-linux-arm64
	tar -czf ${PACKAGE}-${VERSION}-linux-arm64.tar.gz ${PACKAGE}-${VERSION}-linux-arm64

darwin-amd64:
	GOOS=darwin GOARCH=amd64 go build -o ${PACKAGE}-${VERSION}-darwin-amd64
	tar -czf ${PACKAGE}-${VERSION}-darwin-amd64.tar.gz ${PACKAGE}-${VERSION}-darwin-amd64
