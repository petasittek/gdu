
run:
	go run .

build:
	-mkdir build
	cd build; GOOS=linux GOARCH=arm go build -o gdu-linux-arm ..; tar czf gdu-linux-arm.tgz gdu-linux-arm
	cd build; GOOS=linux GOARCH=amd64 go build -o gdu-linux-amd64 ..; tar czf gdu-linux-amd64.tgz gdu-linux-amd64
	cd build; GOOS=windows GOARCH=amd64 go build -o gdu-windows-amd64.exe ..; zip gdu-windows-amd64.zip gdu-windows-amd64.exe
	cd build; GOOS=darwin GOARCH=amd64 go build -o gdu-darwin-amd64 ..; tar czf gdu-darwin-amd64.tgz gdu-darwin-amd64

test:
	go test -v

coverage:
	go test -v -race -coverprofile=coverage.txt -covermode=atomic

clean:
	-rm coverage.txt
	-rm -r test_dir
	-rm -r build

.PHONY: run build test coverage clean