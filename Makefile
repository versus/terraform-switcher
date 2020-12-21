EXE  := tfswitch
PKG  := github.com/versus/terraform-switcher
VER := 0.20.1
PATH := build:$(PATH)
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

ifneq (,$(wildcard ./version))
    include version
    export
endif


$(EXE): go.mod *.go lib/*.go
	go build -v -ldflags "-X main.version=$(VER)" -o ./dist/$@ $(PKG)

.PHONY: release
release: $(EXE) clean gorelease alpine snap

.PHONY: darwin linux
darwin linux:
	GOOS=$@ go build -ldflags "-X main.version=$(VER)" -o ./dist/$(EXE)-$(VER)-$@-$(GOARCH) $(PKG)

.PHONY: clean
clean:
	rm -rf ./dist/

.PHONY: docs
docs:
	cd docs; bundle install --path vendor/bundler; bundle exec jekyll build -c _config.yml; cd ..

.PHONY: snap
snap:
	(multipass stop snapcraft-tfswitch && multipass delete snapcraft-tfswitch && multipass purge && rm -f tfswitch_*.snap) || true  && snapcraft

.PHONY: alpine
alpine:
	cd ./alpine && bash ./build.sh

.PHONY: gorelease
gorelease:
	rm -rf ./dist/
	goreleaser

