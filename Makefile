VERSION ?= 0.0.0-dev

.PHONY: lint
lint: golangci-lint ## Run golang-ci-lint against code.
	$(GOLANGCI_LINT) run ./...

## Location to create the release
RELEASE_DIR ?= $(shell pwd)/release
$(RELEASE_DIR):
	mkdir -p $(RELEASE_DIR)

.PHONY: release
release: gox gomplate $(RELEASE_DIR) ## Run release artifacts
	$(GOX) -output="dist/{{.Dir}}_{{.OS}}_{{.Arch}}" -os="linux darwin windows" -arch="amd64 arm64"
	VERSION=$(VERSION) $(GOMPLATE) -f hack/provider.yaml.tpl > $(RELEASE_DIR)/provider.yaml
	mv dist/* $(RELEASE_DIR)
	rm -rf dist/

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

## Tool Binaries
GOX ?= $(LOCALBIN)/gox
GOMPLATE ?= $(LOCALBIN)/gomplate
GOLANGCI_LINT ?= $(LOCALBIN)/golangci-lint

## Tool Versions
GOX_VERSION ?= latest
GOMPLATE_VERSION ?= v3.11.5
GOLANGCI_LINT_VERSION ?= v1.52.2

.PHONY: gox
gox: $(GOX) 
$(GOX): $(LOCALBIN)
	GOBIN=$(LOCALBIN) go install github.com/mitchellh/gox@$(GOX_VERSION)

.PHONY: gomplate
gomplate: $(GOMPLATE) 
$(GOMPLATE): $(LOCALBIN)
	GOBIN=$(LOCALBIN) go install github.com/hairyhenderson/gomplate/v3/cmd/gomplate@$(GOMPLATE_VERSION)

.PHONY: golangci-lint
golangci-lint: $(GOLANGCI_LINT)
$(GOLANGCI_LINT): $(GOLANGCI_LINT)
	test -s $(LOCALBIN)/golangci-lint && $(LOCALBIN)/golangci-lint version | grep -q $(GOLANGCI_LINT_VERSION) || \
	GOBIN=$(LOCALBIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)
