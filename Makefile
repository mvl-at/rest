GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GORUN=$(GOCMD) run
BUILDDIR?=build
BINARYNAME?=mvl-serve
BINARY=$(BUILDDIR)/$(BINARYNAME)
MAINPACKAGE=./cmd/serve
    
all: test build
prod:
	$(GOBUILD) -o $(BINARYE) -v -ldflags "-s -w" $(MAINPACKAGE)
build:
	$(GOBUILD) -o $(BINARY) -v $(MAINPACKAGE)
test: 
	$(GOTEST) -v ./security_test
live-test:
	$(GOTEST) -v -timeout 0 ./mock_test
clean: 
	$(GOCLEAN)
	rm -rf $(BUILDDIR)
run:
	$(GORUN) $(MAINPACKAGE)/serve.go
deps:
	$(GOGET) ./...
