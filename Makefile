.PHONY: build build_mac build_windows clean gomodgen

build: gomodgen
	export GO111MODULE=on
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -tags demo -o ../bin/main _demo/main.go

build_mac: gomodgen
	export GO111MODULE=on
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w"  -tags demo -pkgdir demo -o ../bin/main_mac _demo/main.go

build_windows: gomodgen
	export GO111MODULE=on
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -tags demo -o ../bin/main.exe _demo/main.go

clean:
	rm -rf ./bin ./vendor Gopkg.lock
