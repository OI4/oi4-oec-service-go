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

docker_prepare:
	docker buildx create --name oi4-builder --driver docker-container --use --platform linux/amd64,linux/arm64,linux/arm/v7,linux/arm/v6,linux/386

docker_build_push:
	docker buildx build --push --platform linux/amd64,linux/arm64,linux/arm/v7,linux/arm/v6 -t oi4a/oi4-oec-service-demo-golang:0.1.0 -t oi4a/oi4-oec-service-demo-golang:latest .

docker_build_local:
	docker build -t oi4-demo-connector-golang:latest .

gomodgen:
	go env -w GO111MODULE=on
	cd api && go mod tidy
	cd container && go mod tidy
	cd dnp && go mod tidy
	cd service && go mod tidy
	cd demo && go mod tidy
