PROJECT			:= github.com/GlobalCyberAlliance/McScrapy

DEP				:= $(shell which dep)
GO				:= $(shell which go)
GO_BUILD_FLAGS	:=
GO_BUILD		:= $(GO) build $(GO_BUILD_FLAGS)
GO_TEST_FLAGS	:= -v -short
GO_TEST			:= $(GO) test $(GO_TEST_FLAGS)
GO_BENCH_FLAGS	:= -short -bench=. -benchmem
GO_BENCH		:= $(GO) test $(GO_BENCH_FLAGS)

TARGETS			:= bin/mcscrapy

all: deps site bin $(TARGETS)

bin:
	mkdir -p $@

bin/%: $(shell find . -name "*.go" -type f)
	$(GO_BUILD) -o $@ $(PROJECT)/cmd/$*

deps: $(DEP)
	$(DEP) ensure

clean:
	-rm -rf bin site

site:
	mkdir -p $@

test: $(DEP)
	$(GO_TEST) $(shell $(DEP) novendor)

bench: $(DEP)
	$(GO_BENCH) $(shell $(DEP) novendor)

.PHONY: deps clean test bench
