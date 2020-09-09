VERSION := $(shell git describe --tags)
LDFLAGS := "-X main.version=$(VERSION)"
GENDIR  := ./gen
CMDDIR  := ./cmd
OBJDIR  := ./obj

# Build commands
.PHONY: build
build: $(OBJDIR)/vert

$(OBJDIR)/%: $(CMDDIR)/%/main.go deps
	go build -ldflags $(LDFLAGS) -o $@ $<

# Clean commands/CSVs
.PHONY: clean
clean:
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
csv: $(OBJDIR)/ucd.csv \
	 $(OBJDIR)/ucdex.csv

$(OBJDIR)/%.csv: $(GENDIR)/gen%csv.go
	@if [ ! -d $(OBJDIR) ]; \
		then echo "mkdir -p $(OBJDIR)"; mkdir -p $(OBJDIR); \
	fi
	go run $< -output $@

# Lint
.PHONY: lint
lint: deps
	go vet ./...
	golint -set_exit_status -min_confidence 0 ./...

# Run tests
.PHONY: test
test: deps
	go test -v -run "Test([^H]|H[^e]|He[^a]|Hea[^v]|Heav[^y])|Example" ./...

.PHONY: alltest
alltest: deps test
	go test -v -timeout 20m -run "Heavy"

#Install dependencies
.PHONY: deps
deps:
	go get github.com/mattn/go-runewidth

