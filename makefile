WIN_BIN=dosserver.exe
LINUX_BIN=dosserver
LINUX_ARM64_BIN=dosserver-arm64
LINUX_ARM_BIN=dosserver-arm
MAC_BIN=dosserver-mac

build: clean
	GOOS=windows GOARCH=amd64 go build -v -o dist/$(WIN_BIN)
	GOOS=linux GOARCH=amd64 go build -v -o dist/$(LINUX_BIN)
	GOOS=linux GOARCH=arm64 go build -v -o dist/$(LINUX_ARM64_BIN)
	GOOS=linux GOARCH=arm go build -v -o dist/$(LINUX_ARM_BIN)
	GOOS=darwin GOARCH=amd64 go build -v -o dist/$(MAC_BIN)

clean:
	rm -rf dist

docker:
	docker build -t dosserver:latest .
	docker tag dosserver:latest ewr.vultrcr.com/cloudregistry/dosserver:latest
	docker push ewr.vultrcr.com/cloudregistry/dosserver:latest
