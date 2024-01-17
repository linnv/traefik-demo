CURDIR := $(shell pwd)

GO        := GOPROXY=https://goproxy.cn go

GOBUILD   := GO111MODULE=on GOPATH=$(GOPATH) CGO_ENABLED=0 $(GO) build $(BUILD_FLAG)
GOTEST    := GOPATH=$(GOPATH) CGO_ENABLED=1 $(GO) test -p 3
OSLINUX   := GOARCH=amd64  GOOS=linux
TARGET=jialinwugo
 
ReposGitBranch := $(shell git rev-parse --abbrev-ref HEAD)
GITBRANCH = $(shell if [ "$(CI_COMMIT_BRANCH)" != "" ]; then echo "$(CI_COMMIT_BRANCH)" ; else echo "$(ReposGitBranch)";fi)

LDFLAGS += -X "github.com/linnv/simdog.Version=$(shell git describe --tags --dirty --always)"
LDFLAGS += -X "github.com/linnv/simdog.BuildTime=$(shell date -u '+%Y-%m-%d %I:%M:%S')"
LDFLAGS += -X "github.com/linnv/simdog.GitBranch=$(GITBRANCH)"
# LDFLAGS += -X "github.com/linnv/simdog.GitBranch=$(shell git branch --show-current)"
LDFLAGS += -X "github.com/linnv/simdog.GitHash=$(shell git rev-parse HEAD)"
LDFLAGS += -X "github.com/linnv/simdog.Owner=jialinwu.com"
LDFLAGS += -X "github.com/linnv/simdog.BuildHost=$(shell hostname)"
LDFLAGS += -X "github.com/linnv/simdog.BuildUser=$(shell whoami)"
VERSION = $(shell git describe --tags --dirty --always)

.PHONY: all linux clean pack packall one vulncheck fmt-check gotest fmt vet golangci-lint

all: $(TARGET)

BUILDDIR=$(CURDIR)

CONFIGFILE="config/config.yaml.sample"

$(TARGET): 
	@mkdir -p $(BUILDDIR)
	@echo "--->>building $(BUILDDIR)/$@"
	$(OSLINUX) $(GOBUILD) -ldflags '$(LDFLAGS)' -v -o $(BUILDDIR)/$@  $(CURDIR)/main.go
	
buildtest:
	@rm -rf $(TARGET)
	@mkdir -p $(BUILDDIR)
	@echo "--->>building linux app $(BUILDDIR)/$@"
	$(GOBUILD) -ldflags '$(LDFLAGS)' -v -o $(BUILDDIR)/$(TARGET)  $(CURDIR)/main.go
	# $(BUILDDIR)/$@ v
	@rm -rf $(TARGET)
linux:
	@mkdir -p $(BUILDDIR)
	@echo "--->>building linux app $(BUILDDIR)/$@"
	$(OSLINUX) $(GOBUILD) -ldflags '$(LDFLAGS)' -v -o $(BUILDDIR)/$(TARGET)  $(CURDIR)/main.go
		
clean: 
	@[ -f $(BUILDDIR)/$(TARGET) ] && rm $(BUILDDIR)/$(TARGET) || true

TARGETAPP=$(TARGET)-$(GITBRANCH)-$(VERSION)
pack:$(TARGET)
	@echo "--->>create dir [$(TARGETAPP)]"
	@mkdir -p $(TARGETAPP)
	@echo "--->>copy app from [$(TARGET)] to [$(TARGETAPP)]"
	@cp $(TARGET) $(TARGETAPP)
	@cp ./testSh/localtest.sh $(TARGETAPP)
	# @echo "--->>compress target [$(TARGETAPP)]"
	# tar -czf $(TARGETAPP).tar.gz  $(TARGETAPP)
	# @echo "--->>remove tmp app [$(TARGETAPP)]"
	# @rm -rf $(TARGETAPP)
	# @echo "--->>remove app [$(TARGET)]"
	# @rm $(TARGET)
	# @echo "--->>done: new app built in [$(TARGETAPP).tar.gz]"
packall:$(TARGET)
	@echo "--->>create dir [$(TARGETAPP)]"
	@mkdir -p $(TARGETAPP)
	@echo "--->>copy app from [$(TARGET)] to [$(TARGETAPP)]"
	@cp $(TARGET) $(TARGETAPP)
	@cp ./testSh/localtest.sh $(TARGETAPP)
	@echo "--->>compress target [$(TARGETAPP)]"
	tar -czf $(TARGETAPP).tar.gz  $(TARGETAPP)
	@echo "--->>remove tmp app [$(TARGETAPP)]"
	@rm -rf $(TARGETAPP)
	@echo "--->>remove app [$(TARGET)]"
	@rm $(TARGET)
	@echo "--->>done: new app built in [$(TARGETAPP).tar.gz]"

GOFMT ?= gofmt "-s"
GO ?= go
PACKAGES ?= $(shell go list ./... | grep -v /vendor/)
SOURCES ?= $(shell find . -path ./vendor -prune -o -name "*.go" -type f)

gotest:
	$(GO) test -race $(PACKAGES)

fmt-check:
	@diff=$$($(GOFMT) -d $(SOURCES)); \
	if [ -n "$$diff" ]; then \
		echo "Please run 'make fmt' and commit the result:"; \
		echo "$${diff}"; \
		exit 1; \
	fi;

fmt:
	$(GOFMT) -w $(SOURCES)

vet:
	$(GO) vet $(PACKAGES)

golangci-lint:
	@echo "skip"
	# golangci-lint run
vulncheck:
	# govulncheck $(PACKAGES)
	
linuxarm:
	@mkdir -p $(BUILDDIR)
	@echo "--->>building linux app $(BUILDDIR)/$@"
	GOOS=linux GOARCH=arm $(GOBUILD) -ldflags '$(LDFLAGS)' -v -o $(BUILDDIR)/$(TARGET)  $(CURDIR)/main.go
armv7:
	@mkdir -p $(BUILDDIR)
	@echo "--->>building linux app $(BUILDDIR)/$@"
	GOOS=linux GOARCH=arm GOARM=7  $(GOBUILD) -ldflags '$(LDFLAGS)' -v -o $(BUILDDIR)/$(TARGET)  $(CURDIR)/main.go
armv6:
	@mkdir -p $(BUILDDIR)
	@echo "--->>building linux app $(BUILDDIR)/$@"
	GOOS=linux GOARCH=arm GOARM=6  $(GOBUILD) -ldflags '$(LDFLAGS)' -v -o $(BUILDDIR)/$(TARGET)  $(CURDIR)/main.go
