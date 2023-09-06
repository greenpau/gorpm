.PHONY: build test ctest covdir coverage docs linter qtest clean dep
APP_VERSION:=$(shell cat VERSION | head -1)
GIT_COMMIT:=$(shell git describe --dirty --always)
GIT_BRANCH:=$(shell git rev-parse --abbrev-ref HEAD -- | head -1)
BUILD_USER:=$(shell whoami)
BUILD_DATE:=$(shell date +"%Y-%m-%d")
PROJECT:="github.com/greenpau/gorpm"
BINARY:="gorpm"
CORE_PKG="gorpm"
VERBOSE:=-v
ifdef TEST
	TEST:="-run ${TEST}"
endif

all: build

build:
	@echo "Version: $(APP_VERSION), Branch: $(GIT_BRANCH), Revision: $(GIT_COMMIT)"
	@echo "Build on $(BUILD_DATE) by $(BUILD_USER)"
	@mkdir -p bin/
	@CGO_ENABLED=0 go build -o bin/$(BINARY) $(VERBOSE) \
		-ldflags="-w -s \
		-X main.appName=$(BINARY) \
		-X main.appVersion=$(APP_VERSION) \
		-X main.gitBranch=$(GIT_BRANCH) \
		-X main.gitCommit=$(GIT_COMMIT) \
		-X main.buildUser=$(BUILD_USER) \
		-X main.buildDate=$(BUILD_DATE)" \
		-gcflags="all=-trimpath=$(GOPATH)/src" \
		-asmflags="all=-trimpath $(GOPATH)/src" cmd/$(BINARY)/*
	@echo "Done!"

linter:
	@golint pkg/$(CORE_PKG)/*.go
	@#golint cmd/client/*.go
	@echo "PASS: golint"

test: covdir linter
	@go test $(VERBOSE) -coverprofile=.coverage/coverage.out ./pkg/$(CORE_PKG)/*.go
	@bin/$(BINARY) --version

ctest: covdir linter
	@richgo version || go get -u github.com/kyoh86/richgo
	@time richgo test $(VERBOSE) "${TEST}" -coverprofile=.coverage/coverage.out ./pkg/$(CORE_PKG)/*.go

covdir:
	@mkdir -p .coverage

coverage:
	@go tool cover -html=.coverage/coverage.out -o .coverage/coverage.html

docs:
	@rm -rf .doc/
	@mkdir -p .doc/
	@godoc -html $(PROJECT)/pkg/$(CORE_PKG) > .doc/index.html
	@echo "Run to serve docs:"
	@echo "    godoc -goroot .doc/ -html -http \":8080\""

clean:
	@rm -rf .doc
	@rm -rf .coverage
	@rm -rf bin/

qtest:
	@go test -v -run TestNewSpec ./pkg/$(CORE_PKG)/*.go
	@#go test -v -run TestNewBuild ./pkg/$(CORE_PKG)/*.go

install:
	@sudo rm -rf /usr/local/bin/gorpm
	@sudo cp bin/gorpm /usr/local/bin/gorpm
	@echo "Installed /usr/local/bin/gorpm"

dep:
	@echo "Making dependencies check ..."
	@go install golang.org/x/lint/golint@latest
	@go install github.com/kyoh86/richgo@latest
	@go install github.com/greenpau/versioned/cmd/versioned@latest

release:
	@echo "Making release"
	@go mod tidy
	@go mod verify
	@if [ $(GIT_BRANCH) != "main" ]; then echo "cannot release to non-main branch $(GIT_BRANCH)" && false; fi
	@git diff-index --quiet HEAD -- || ( echo "git directory is dirty, commit changes first" && false )
	@versioned -patch
	@echo "Patched version"
	@git add VERSION
	@git commit -m "released v`cat VERSION | head -1`"
	@git tag -a v`cat VERSION | head -1` -m "v`cat VERSION | head -1`"
	@git push
	@git push --tags
	@@echo "If necessary, run the following commands:"
	@echo "  git push --delete origin v$(APP_VERSION)"
	@echo "  git tag --delete v$(APP_VERSION)"

logo:
	@mkdir -p assets/docs/images
	@convert -background black -fill white -font DejaVu-Sans-Bold -size 640x320! -gravity center -pointsize 96 label:'GoRPM' PNG32:assets/docs/images/logo.png
