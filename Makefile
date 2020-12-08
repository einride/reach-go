all: \
	commitlint \
	go-generate \
	go-lint \
	go-review \
	go-test \
	go-mod-tidy \
	git-verify-nodiff

include tools/git-verify-nodiff/rules.mk
include tools/commitlint/rules.mk
include tools/golangci-lint/rules.mk
include tools/goreview/rules.mk
include tools/stringer/rules.mk

.PHONY: clean
clean:
	$(info [$@] cleaning generated files...)
	@find -name '*_string.go' -exec rm {} \+

.PHONY: go-mod-tidy
go-mod-tidy:
	$(info [$@] tidying Go module files...)
	@go mod tidy -v

.PHONY: go-test
go-test:
	$(info [$@] running Go tests...)
	@go test -count 1 -cover -race ./...

.PHONY: go-generate
go-generate:
	$(info [$@] generating Go code...)
	@go generate ./...
