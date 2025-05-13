
lint:
	$(shell go env GOPATH)/bin/golangci-lint run

lint/fix:
	$(shell go env GOPATH)/bin/golangci-lint run --fix
