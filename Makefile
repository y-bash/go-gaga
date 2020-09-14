VERSION := $(shell git describe --tags)
LDFLAGS := "-X main.version=$(VERSION)"
GENDIR  := ./gen
CMDDIR  := ./cmd
OBJDIR  := ./obj

# Test Function REGEXP
LIGHTW  := "Test([^H]|H[^e]|He[^a]|Hea[^v]|Heav[^y])|Example"
HEAVYW  := "Heavy"
BENCHM  := "Benchmark"


# Build commands
.PHONY: build
build: $(OBJDIR)/vert

$(OBJDIR)/%: $(CMDDIR)/%/main.go deps
	go build -ldflags $(LDFLAGS) -o $@ $<

# Clean commands/CSVs
.PHONY: clean
clean:
	go clean -testcache
	rm -rf $(OBJDIR)

# Install commands
.PHONY: install
install:
	go install -ldflags $(LDFLAGS) ./...

# Uninstall commands
.PHONY: uninstall
uninstall:
	go clean -i ./...

# Generate unichar_tables.go
.PHONY: gen
gen:
	go generate .

# Generate CSVs
.PHONY: csv
csv: dir \
	$(OBJDIR)/ucd.csv \
	$(OBJDIR)/ucdex.csv

$(OBJDIR)/%.csv: $(GENDIR)/gen%csv.go
	go run $< -output $@

# Lint
.PHONY: lint
lint: devdeps
	go vet ./...
	golint -set_exit_status -min_confidence 0 ./...

# Run tests
.PHONY: test
test: devdeps dir
	go test -v -coverprofile=$(OBJDIR)/cover.out -run $(LIGHTW) ./...
	go tool cover -html=$(OBJDIR)/cover.out -o $(OBJDIR)/cover.html

.PHONY: alltests
alltests: test
	go test -v -timeout 20m -run $(HEAVYW)

#Benchmarks
.PHONY: bench
bench: devdeps dir
	go test -v -run $(BENCHM) -bench . -benchmem -o $(OBJDIR)/bench.bin -cpuprofile=$(OBJDIR)/cpu.prof -memprofile=$(OBJDIR)/mem.prof

#Install dependencies
.PHONY: deps
deps:
	go get github.com/mattn/go-runewidth

.PHONY: devdeps
devdeps: deps
	go get golang.org/x/lint
	go get golang.org/x/text

#Make directory
.PHONY: dir
dir:
	@if [ ! -d $(OBJDIR) ]; \
		then echo "mkdir -p $(OBJDIR)"; mkdir -p $(OBJDIR); \
	fi

