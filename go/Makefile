GO      	:= go
GOFLAGS 	:= -ldflags="-s -w" # Optimization flags

BIN  			:= ./bin
PKG 			:= ./pkg
CMD 			:= ./cmd
INT				:= ./internal

FILES 		:= $(CMD)/**/*.go $(PKG)/**/*go $(INT)/**/*.go
CMDS 			:= $(notdir $(wildcard $(CMD)/*/))

UDP_FILES := $(CMD)/udp/**/*.go $(PKG)/udp/**/*.go

.PHONY: all
all: $(addprefix $(BIN)/, $(CMDS))

build: $(FILES)
	$(addprefix build-, $(CMDS))

build-%: $(shell find $(CMD)/$* $(PKG)/$* -name '*.go' 2>/dev/null)
	$(GO) build $(GOFLAGS) -o $(BIN)/$* $(CMD)/$*/main.go

run-%: $(shell find $(CMD)/$* $(PKG)/$* -name '*.go' 2>/dev/null)
	CompileDaemon -command="$(GO) run $(CMD)/$*/main.go" -directory=$(PKG)/$* -color=true

test: $(FILES)
	$(GO) test ./... -v -race

test-watch-%: $(shell find $(CMD)/$* $(PKG)/$* -name '*.go' 2>/dev/null)
	PKG_NAME=$*; \
	CompileDaemon \
		-command="$(GO) test $(PKG)/$${PKG_NAME}" \
		-directory="$(PKG)/$${PKG_NAME}" \
		-pattern='(.+\.go)$$' \
		-graceful-kill=true \
		-color=true

test-%: $(shell find $(CMD)/$* $(PKG)/$* -name '*.go' 2>/dev/null)
	$(GO) test $(PKG)/$*

lint: $(FILES)
	golangci-lint run ./...

lint-%: $(shell find $(CMD)/$* $(PKG)/$* -name '*.go' 2>/dev/null)
	golangci-lint run ./$(PKG)/$*

clean:
	$(GO) clean
	rm -rf $(BIN)/*
