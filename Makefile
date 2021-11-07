EXECUTABLE=dyn-https
WINDOWS=$(EXECUTABLE)_windows_amd64.exe
LINUX=$(EXECUTABLE)_linux_amd64
DARWIN=$(EXECUTABLE)_darwin_amd64

build: windows ## Build binaries

windows: $(WINDOWS) ## Build for Windows

linux: $(LINUX) ## Build for Linux

darwin: $(DARWIN) ## Build for Darwin (macOS)

$(WINDOWS):
	go build -i -v -ldflags="-s -w -X 'main.GitHash=$(git describe --always --long --dirty)' -X 'main.Version=$(Get-Date -Format "yy.MM.dd")'"

$(LINUX):
	go build -i -v -ldflags="-s -w -X 'main.GitHash=$(git describe --always --long --dirty)' -X 'main.Version=$(date +'%y.%m.%d')'"

$(DARWIN):
	go build -i -v -ldflags="-s -w"

clean: ## Remove previous build
	rm -f $(EXECUTABLE) $(WINDOWS) $(LINUX) $(DARWIN)

help: ## Display available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

# working version and commit hash build example
# Windows Example
# go build -i -v -ldflags="-s -w -X 'main.GitHash=$(git describe --always --long --dirty)' -X 'main.Version=$(Get-Date -Format "yy.MM.dd")'" github.com/duality-solutions/dyn-https
# Bash Example
# go build -i -v -ldflags="-s -w -X 'main.GitHash=$(git describe --always --long --dirty)' -X 'main.Version=$(date +'%y.%m.%d')'" github.com/duality-solutions/dyn-https