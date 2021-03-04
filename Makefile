PACKAGE=myTldr

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
