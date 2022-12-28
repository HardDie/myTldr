PACKAGE=myTldr
VERSION=v0.9.0

.PHONY: default
default: help

.PHONY: help
help: ## help information about make commands
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "[36m%-30s[0m %s\n", $$1, $$2}'

.PHONY: linters
linters: ## run linter
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

.PHONY: build_all
build_all: ## build binaries for all platforms
	make linux-386
	make linux-amd64
	make linux-arm
	make linux-arm64
	make darwin-amd64

.PHONY: linux-386
linux-386: ## Build binary for GOOS=linux GOARCH=386
	GOOS=linux GOARCH=386 go build -o ${PACKAGE}-${VERSION}-linux-386
	tar -czf ${PACKAGE}-${VERSION}-linux-386.tar.gz ${PACKAGE}-${VERSION}-linux-386
	rm ${PACKAGE}-${VERSION}-linux-386

.PHONY: linux-amd64
linux-amd64: ## Build binary for GOOS=linux GOARCH=amd64
	GOOS=linux GOARCH=amd64 go build -o ${PACKAGE}-${VERSION}-linux-amd64
	tar -czf ${PACKAGE}-${VERSION}-linux-amd64.tar.gz ${PACKAGE}-${VERSION}-linux-amd64
	rm ${PACKAGE}-${VERSION}-linux-amd64

.PHONY: linux-arm
linux-arm: ## Build binary for GOOS=linux GOARCH=arm
	GOOS=linux GOARCH=arm go build -o ${PACKAGE}-${VERSION}-linux-arm
	tar -czf ${PACKAGE}-${VERSION}-linux-arm.tar.gz ${PACKAGE}-${VERSION}-linux-arm
	rm ${PACKAGE}-${VERSION}-linux-arm

.PHONY: linux-arm64
linux-arm64: ## Build binary for GOOS=linux GOARCH=arm64
	GOOS=linux GOARCH=arm64 go build -o ${PACKAGE}-${VERSION}-linux-arm64
	tar -czf ${PACKAGE}-${VERSION}-linux-arm64.tar.gz ${PACKAGE}-${VERSION}-linux-arm64
	rm ${PACKAGE}-${VERSION}-linux-arm64

.PHONY: darwin-amd64
darwin-amd64: ## Build binary for GOOS=darwin GOARCH=amd64
	GOOS=darwin GOARCH=amd64 go build -o ${PACKAGE}-${VERSION}-darwin-amd64
	tar -czf ${PACKAGE}-${VERSION}-darwin-amd64.tar.gz ${PACKAGE}-${VERSION}-darwin-amd64
	rm ${PACKAGE}-${VERSION}-darwin-amd64
