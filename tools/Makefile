GO_BIN?=$(shell pwd)/.bin
SHELL:=env PATH=$(GO_BIN):$(PATH) $(SHELL)

tools::
	mkdir -p $(GO_BIN)
	go mod tidy
	@cat tools.go | grep '_' | awk -F '"' '{print $$2}' | xargs -tI % sh -c 'GOBIN=${GO_BIN} go install %'
