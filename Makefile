# all: run a complete build
all: \
	markdown-lint \
	go-mod-tidy \
	go-generate \
	go-lint \
	go-review \
	go-test \
	git-verify-nodiff \
	git-verify-submodules

export GO111MODULE := on

# clean: remove all generated build files
.PHONY: clean
clean: clean-go-generate
	rm -rf build

.PHONY: clean-go-generate
clean-go-generate:
	find -name '*_string.go' -exec rm {} \+

.PHONY: build
build:
	@git submodule update --init --recursive $@

include build/rules.mk
build/rules.mk: build
	@# included in submodule: build

# markdown-lint: lint Markdown files
.PHONY: markdown-lint
markdown-lint: $(MARKDOWNLINT)
	$(MARKDOWNLINT) --ignore build .

# go-mod-tidy: update Go module files
.PHONY: go-mod-tidy
go-mod-tidy:
	go mod tidy -v

# go-lint: lint Go files
.PHONY: go-lint
go-lint: $(GOLANGCI_LINT)
	# maligned: disabled to skip manual aligment of Scanner state
	$(GOLANGCI_LINT) run --enable-all --disable maligned

# go-test: run Go test suite
.PHONY: go-test
go-test:
	go test -count 1 -cover -race ./...

# go-review: review Go files
.PHONY: go-review
go-review: $(GOREVIEW)
	$(GOREVIEW) -c 1 ./...

# go-generate: generate Go code
.PHONY: go-generate
go-generate: \
	pkg/erb/fixtype_string.go \
	pkg/erb/svtype_string.go

pkg/erb/fixtype_string.go: pkg/erb/fixtype.go $(GOBIN)
	$(GOBIN) -m -run golang.org/x/tools/cmd/stringer \
		-type FixType -trimprefix FixType -output $@ $<

pkg/erb/svtype_string.go: pkg/erb/svtype.go $(GOBIN)
	$(GOBIN) -m -run golang.org/x/tools/cmd/stringer \
		-type SVType -trimprefix SVType -output $@ $<
