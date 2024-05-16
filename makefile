WIN_BIN=dosserver.exe
LINUX_BIN=dosserver
MAC_BIN=dosserver-mac

build:
	GOOS=windows GOARCH=amd64 go build -v -o dist/$(WIN_BIN)
	GOOS=linux GOARCH=amd64 go build -v -o dist/$(LINUX_BIN)
	GOOS=darwin GOARCH=amd64 go build -v -o dist/$(MAC_BIN)

clean:
	rm -rf dist

