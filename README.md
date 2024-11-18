# Open Industry 4.0 Alliance - Open Edge Computing Service - Go

The oi4-service repository contains all artifacts to build a service compliant to [OEC Development Guideline v1.1.0](docs/OI4_OEC_Development_Guideline_V1.1.0.pdf).
It also covers a base service that will help you build your own compliant OI4 OEC service by providing most of the OI4 state machine.

The oi4-oec-service-go is a community project and is not offered by the OI4 Alliance. It is published under the MIT license
AND most important it is an OI4 open source project that needs your contribution!
If you want to contribute and do not know how, just reach out for the WG leads.

## Prerequisites
The `oi4-oec-service-go` is a Golang application. To get started you need a proper Go installation.
Please install [Go](https://go.dev/doc/install).

The service requires Go version 1.23.1 or higher.

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
There is a demo application that shows how to use the oi4-oec-service-go.

### Preconfigured service
The OI4 OEC guideline defines a set of configuration files that are used to configure the service. 
The demo application comes with a set of preconfigured files that can be used to run the service.
The configuration files are stored in the folder `_demo/testdata`.

### Run the demo service
The demo service can also be run locally. The service is a Golang application and can be started with the following command:

```sh
go run _demo/main.go -tags demo
```

### Building the demo service as binary
The demo service can also be build as a binary. The binaries can be build with a simple make command, depending on the target system.
- Linux: `make build`
- MacOs `make build_mac`
- Windows `make build_windows`

The binaries will be stored in the bin directory. To execute the binary, just run the binary with the runtime flag and the base path to the configuration files.
- Linux: `./bin/main`
- MacOs `./bin/main_mac`
- Windows `./bin/main.exe`
