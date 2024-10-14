# Open Industry 4.0 Alliance - Open Edge Computing Service - Go

The oi4-service repository contains all artifacts to build a service compliant to [OEC Development Guideline v1.1.0](docs/OI4_OEC_Development_Guideline_V1.1.0.pdf).
It also covers a base service that will help you build your own compliant OI4 OEC service by providing most of the OI4 state machine.

The oi4-oec-service-go is a community project and is not offered by the OI4 Alliance. It is published under the MIT license
AND most important it is an OI4 open source project that needs your contribution!
If you want to contribute and do not know how, just reach out for the WG leads.

## Prerequisites
The `oi4-oec-service-go` is a Golang application. To get started you need a proper Go installation.
Please install [Go](https://go.dev/doc/install).

The service requires Go version 1.23.0 or higher.

The build process uses [Make](https://www.gnu.org/software/make/). Make sure you have it installed.

On Linux you can install it with apt:
```sh
sudo apt-get install make
```

On MacOs you can install it with brew:
```sh
brew install make
```
 On Windows you can install it with [chocolatey](https://chocolatey.org/):
```sh
choco install make
```

## Getting started
The oi4-oec-service-go is deigned to hide the complexity of the OI4 OEC state machine and still provide a flexible way to implement your own business logic.
There is a demo application that shows how to use the oi4-oec-service-go. The demo application can be run as a local application or as a docker container.

The latest version of the demo application can be found at [docker hub](https://hub.docker.com/repository/docker/oi4a/oi4-oec-service-demo-go/general)

### Preconfigured service and docker data
The OI4 OEC guideline defines a set of configuration files that are used to configure the service. 
The demo application comes with a set of preconfigured files that can be used to run the service as a docker container.
The configuration files are stored in the docker_configs.zip file. The docker_configs.zip file can be found in the root directory of the demo application and must be unzipped to the docker_configs directory.

```sh
unzip -d ./docker_configs docker_configs.zip
```

### Run the demo service as docker container
The docker configuration already contains a docker-compose file that can be used to run the demo service as a docker container.

```sh
cd docker_configs/oi4-oec-service-demo-go
docker compose up -d
```

### Run the demo service locally
The demo service can also be run locally. The service is a Golang application and can be started with the following command:

```sh
cd demo
go run main.go -runtime=local base=../docker_configs
```

### Building the demo service as binary
The demo service can also be build as a binary. The binaries can be build with a simple make command, depending on the target system.
- Linux: `make build`
- MacOs `make build_mac`
- Windows `make build_windows`

The binaries will be stored in the bin directory. To execute the binary, just run the binary with the runtime flag and the base path to the configuration files.
- Linux: `./bin/main -runtime=local base=docker_configs`
- MacOs `./bin/main_mac -runtime=local base=docker_configs`
- Windows `./bin/main.exe -runtime=local base=docker_configs`

### Build and run the demo service as docker container
The demo service can also be build as a docker container locally. The docker container can be build with the following make command:

```sh
make docker_build_local
```

## Update go.mod dependencies
The go.mod file is used to manage the dependencies of the project. It is important to keep it up to date.
There is script to simplify the update of the go.mod file. Just run:

```sh
./gomod.sh
```

unzip -d ./docker docker_configs.zip
