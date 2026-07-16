.PHONY: fmt test vet lint vuln check

fmt:
	gofmt -w .

test:
	go test ./...

vet:
	go vet ./...

lint:
	go run honnef.co/go/tools/cmd/staticcheck@latest ./...

vuln:
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

check: test vet lint vuln
