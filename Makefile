lightcontrol: go.mod go.sum lightcontrol.go
	env GOARM=5 GOARCH=arm GOOS=linux go build -trimpath -ldflags '-s' -o lightcontrol

.PHONY: clean
clean:
	rm -f lightcontrol
