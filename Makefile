.PHONY: build build_mac build_windows clean gomodgen

build: gomodgen
	export GO111MODULE=on
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ../bin/main demo/main.go

build_mac: gomodgen
	export GO111MODULE=on
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o ../bin/main_mac demo/main.go

build_windows: gomodgen
	export GO111MODULE=on
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o ../bin/main.exe demo/main.go

clean:
	rm -rf ./bin ./vendor Gopkg.lock

gomodgen:
	go env -w GO111MODULE=on
	cd api && go mod tidy
	cd container && go mod tidy
	cd dnp && go mod tidy
	cd service && go mod tidy
	cd demo && go mod tidy
