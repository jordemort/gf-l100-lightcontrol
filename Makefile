export GOARM ?= 5
export GOARCH ?= arm
export GOOS ?= linux
export GOLDFLAGS ?= -s

lightcontrol: go.mod go.sum lightcontrol.go
	go build -trimpath -ldflags '$(GOLDFLAGS)' -o lightcontrol
	upx --best lightcontrol

.PHONY: clean
clean:
	rm -f lightcontrol
