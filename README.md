# Dynamic HTTPS Proxy REST Server

HTTPS Proxy server for Dynamic JSON-RPC

## Development Requirements

- [GoLang](https://golang.org/dl/) v1.14 or above.

## Running with dev and debug mode

```go
go run -race . -dev -debug
```

### Building

```go
# Linux
go build -i -v -ldflags="-X 'main.GitHash=$(git describe --always --long --dirty)' -X 'main.Version=$(date +'%y.%m.%d')'" github.com/duality-solutions/dyn-https
```

```go
# Windows
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

### License and Copyrights

See [LICENSE.md](./LICENSE.md "LICENSE.md") file for copyright, copying and use information.
