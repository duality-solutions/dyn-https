# Dynmaic HTTPS Proxy REST Server

![#](assets.png)

Connecting HTTP servers and clients on disparate networks using WebRTC and blockchain signaling

## Development Requirements

- [GoLang](https://golang.org/dl/) v1.14 or above.

## Running with dev and debug mode

```go
go run -race . -dev -debug
```

### Building

```go
go build -i -v -ldflags="-X 'main.GitHash=$(git describe --always --long --dirty)' -X 'main.Version=$(date +'%y.%m.%d')'" github.com/duality-solutions/dyn-https
```

```go
# Windows Requires protobuf compiler: https://github.com/protocolbuffers/protobuf/releases
go build -i -v -ldflags="-X 'main.GitHash=$(git describe --always --long --dirty)' -X 'main.Version=$(Get-Date -Format "yy.MM.dd")'" github.com/duality-solutions/dyn-https
```

#### Windows NMake

```shell
nmake /f Makefile
```

#### Linux Make

```bash
make
```

### Diagrams

![General Diagram](docs/diagrams/webbridge-general.png)

![Technical Details Diagram](docs/diagrams/webbridge-tech-details.png)

### License and Copyrights

See [LICENSE.md](./LICENSE.md "LICENSE.md") file for copyright, copying and use information.
