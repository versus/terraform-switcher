EXE  := tfswitch
PKG  := github.com/versus/terraform-switcher
VER := 0.21.9
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
release: $(EXE) clean tag gorelease alpine 

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
	(mkdir ./dist && multipass stop snapcraft-tfswitch && multipass delete snapcraft-tfswitch && multipass purge && rm -f ./dist/tfswitch_*.snap) || true  && snapcraft && mv tfswitch*.snap ./dist

.PHONY: snap-stop
snap-stop:
	(multipass stop snapcraft-tfswitch && multipass delete snapcraft-tfswitch && multipass purge ) || true
.PHONY: alpine
alpine:
	cd ./alpine && bash ./build.sh

.PHONY: gorelease
gorelease:
	rm -rf ./dist/
	goreleaser

.PHONY: tag
tag:
	git tag -a $(VER) -m "New release"


