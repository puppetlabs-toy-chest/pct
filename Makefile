quality: format lint sec tidy

# Run go mod tidy and check go.sum is unchanged
PHONY+= tidy
tidy:
	@echo "🔘 Checking that go mod tidy does not make a change..."
	@cp go.sum go.sum.bak
	@go mod tidy
	@diff go.sum go.sum.bak && rm go.sum.bak || (echo "🔴 go mod tidy would make a change, exiting"; exit 1)
	@echo "✅ Checking go mod tidy complete"

# Format go code and error if any changes are made
PHONY+= format
format:
	@echo "🔘 Checking that go fmt does not make any changes..."
	@test -z $$(go fmt ./...) || (echo "🔴 go fmt would make a change, exiting"; exit 1)
	@echo "✅ Checking go fmt complete"

PHONY+= lint
lint: $(GOPATH)/bin/golangci-lint
	@echo "🔘 Linting $(1) (`date '+%H:%M:%S'`)"
	@lint=`golint ./...`; \
	if [ "$$lint" != "" ]; \
	then echo "🔴 Lint found by golint"; echo "$$lint"; exit 1;\
	fi
	@lint=`go vet ./...`; \
	if [ "$$lint" != "" ]; \
	then echo "🔴 Lint found by go vet"; echo "$$lint"; exit 1;\
	fi
	@lint=`golangci-lint run`; \
	if [ "$$lint" != "" ]; \
	then echo "🔴 Lint found by golangci-lint"; echo "$$lint"; exit 1;\
	fi
	@echo "✅ Lint-free (`date '+%H:%M:%S'`)"

PHONY+= sec
sec: $(GOPATH)/bin/gosec
	@echo "🔘 Checking for security problems ... (`date '+%H:%M:%S'`)"
	@sec=`gosec -exclude-dir=testutils -quiet ./...`; \
	if [ "$$sec" != "" ]; \
	then echo "🔴 Problems found"; echo "$$sec"; exit 1;\
	else echo "✅ No problems found (`date '+%H:%M:%S'`)"; \
	fi

default: quality

.PHONY: $(PHONY)
