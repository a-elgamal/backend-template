MODULE   = $(shell $(GO) list -m)
DATE    ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ" || echo unknown)
VERSION ?= $(shell prefix=$$(echo $(GIT_TAG) | cut -c 1); if [ "$${prefix}" = "v" ]; then echo $(GIT_TAG) | cut -c 2- ; else echo $(GIT_TAG) ; fi)
GIT_TAG ?= $(shell git describe --tags 2>/dev/null || echo unknown)
GIT_COMMIT ?= $(shell git rev-parse HEAD 2>/dev/null || echo unknown)

BUILD_INFO_FIELDS := \
	Version=$(VERSION) \
	GitTag=$(GIT_TAG) \
	GitCommit=$(GIT_COMMIT) \
	BuildDate=$(DATE)

PKGS     = $(or $(PKG),$(shell $(GO) list ./...))
BIN      = bin

GO      = go
V = 0
Q = $(if $(filter 1,$V),,@)
M = $(shell if [ "$$(tput colors 2> /dev/null || echo 0)" -ge 8 ]; then printf "\033[34;1m▶\033[0m"; else printf "▶"; fi)

LDFLAGS += $(foreach entry,$(BUILD_INFO_FIELDS),-X $(MODULE)/internal.$(entry))

.SUFFIXES:
.PHONY: all
all: lint build test $(BIN)

# Tools

$(BIN):
	@mkdir -p $@
$(BIN)/%: | $(BIN) ; $(info $(M) building $(PACKAGE)…)
	$Q env GOBIN=$(abspath $(BIN)) $(GO) install $(PACKAGE)

GOIMPORTS = $(BIN)/goimports
$(BIN)/goimports: PACKAGE=golang.org/x/tools/cmd/goimports@latest

REVIVE = $(BIN)/revive
$(BIN)/revive: PACKAGE=github.com/mgechev/revive@latest

GOTESTSUM = $(BIN)/gotestsum
$(BIN)/gotestsum: PACKAGE=gotest.tools/gotestsum@latest

# Tests
COVERAGE_MODE = count
COVERAGE_MIN = 70
.PHONY: test
test:
test: | $(GOTESTSUM) ; $(info $(M) running coverage tests…) @ ## Run coverage tests locally
	$Q mkdir -p bin/test
	$Q $(GOTESTSUM) -- \
		-coverpkg=$(shell echo $(PKGS) | tr ' ' ',') \
		-covermode=$(COVERAGE_MODE) \
		-coverprofile=bin/test/profile.out $(PKGS)
	$Q $(GO) tool cover -html=bin/test/profile.out -o bin/test/coverage.html
	$Q $(GO) tool cover -func=bin/test/profile.out -o=bin/test/coverage.out
	$Q cat bin/test/coverage.out | grep -i total:
	$Q cat bin/test/coverage.out | gawk '/total:.*statements/ {if (strtonum($$3) < $(COVERAGE_MIN)) {print "ERR: coverage is lower than $(COVERAGE_MIN)"; exit 1}}'

.PHONY: test-ci
test-ci:
test-ci: | $(GOTESTSUM) ; $(info $(M) running coverage tests…) @ ## Run coverage tests in CI
	$Q mkdir -p bin/test
	$Q $(GOTESTSUM) --junitfile bin/test/unit-tests.xml --format github-actions -- \
		-coverpkg=$(shell echo $(PKGS) | tr ' ' ',') \
		-covermode=$(COVERAGE_MODE) \
		-coverprofile=bin/test/profile.out $(PKGS)
	$Q $(GO) tool cover -func=bin/test/profile.out -o=bin/test/coverage.out
	$Q cat bin/test/coverage.out | grep -i total:

.PHONY: lint
lint: | $(REVIVE) ; $(info $(M) running golint…) @ ## Run golint
	$Q $(REVIVE) -formatter friendly -set_exit_status ./...

.PHONY: fmt
fmt: | $(GOIMPORTS) ; $(info $(M) running gofmt…) @ ## Run gofmt on all source files
	$Q $(GOIMPORTS) -local $(MODULE) -w $(shell $(GO) list -f '{{$$d := .Dir}}{{range $$f := .GoFiles}}{{printf "%s/%s\n" $$d $$f}}{{end}}{{range $$f := .CgoFiles}}{{printf "%s/%s\n" $$d $$f}}{{end}}{{range $$f := .TestGoFiles}}{{printf "%s/%s\n" $$d $$f}}{{end}}' $(PKGS))

.PHONY: build
build: ;$(info $(M) building…) ## Runing go build
	$Q $(GO) build -ldflags '$(LDFLAGS)' -o $(BIN)/

.PHONY: run
run: build
run: ;$(info $(M) running…) ## Runing the service
	$Q bin/myservice start

.PHONY: setup
setup: ; $(info $(M) running setup wizard…) @ ## Run the initial setup wizard
	@bash setup.sh

# Misc

.PHONY: clean
clean: ; $(info $(M) cleaning…)	@ ## Cleanup everything
	@rm -rf $(BIN) test

.PHONY: help
help:
	@grep -hE '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-17s\033[0m %s\n", $$1, $$2}'

.PHONY: version
version:
	@echo $(VERSION)
