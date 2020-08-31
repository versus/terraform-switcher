EXE  := tfswitch
PKG  := github.com/versus/terraform-switcher
VER := 0.9.12
PATH := build:$(PATH)
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

$(EXE): go.mod *.go lib/*.go
	go build -v -ldflags "-X main.version=$(VER)" -o ./dist/$@ $(PKG)

.PHONY: release
release: $(EXE) darwin linux

.PHONY: darwin linux 
darwin linux:
	GOOS=$@ go build -ldflags "-X main.version=$(VER)" -o ./dist/$(EXE)-$(VER)-$@-$(GOARCH) $(PKG)

.PHONY: clean
clean:
	rm -f ./dist/$(EXE) ./dist/$(EXE)-*-*-*

.PHONY: test
test: $(EXE)
	mv ./dist/$(EXE) build
	go test -v ./...

.PHONY: docs
docs:
	cd docs; bundle install --path vendor/bundler; bundle exec jekyll build -c _config.yml; cd ..

