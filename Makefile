.PHONY: test vet fmt lint godoc

test:
	go test -v ./...

vet:
	go vet ./...

fmt:
	go install mvdan.cc/gofumpt@latest
	gofumpt -l -w .

lint:
	golangci-lint run

godoc:
	go install golang.org/x/tools/cmd/godoc@latest
	@echo "Starting godoc server at http://localhost:6060/pkg/github.com/sirockin/tickers/"
	godoc -http=:6060
