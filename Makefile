PACKAGE=myTldr

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
	GOOS=linux GOARCH=386 go build -o ${PACKAGE}-linux-386
	tar -czf ${PACKAGE}-linux-386.tar.gz ${PACKAGE}-linux-386

linux-amd64:
	GOOS=linux GOARCH=amd64 go build -o ${PACKAGE}-linux-amd64
	tar -czf ${PACKAGE}-linux-amd64.tar.gz ${PACKAGE}-linux-amd64

linux-arm:
	GOOS=linux GOARCH=arm go build -o ${PACKAGE}-linux-arm
	tar -czf ${PACKAGE}-linux-arm.tar.gz ${PACKAGE}-linux-arm

linux-arm64:
	GOOS=linux GOARCH=arm64 go build -o ${PACKAGE}-linux-arm64
	tar -czf ${PACKAGE}-linux-arm64.tar.gz ${PACKAGE}-linux-arm64

darwin-amd64:
	GOOS=darwin GOARCH=amd64 go build -o ${PACKAGE}-darwin-amd64
	tar -czf ${PACKAGE}-darwin-amd64.tar.gz ${PACKAGE}-darwin-amd64
