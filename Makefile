.PHONY: build check test test-windows clean

build:
	go build ./cmd/... ./internal/...
	npm run build

check:
	gofmt -w $$(find cmd internal -name '*.go' -type f)
	go vet ./cmd/... ./internal/...
	npm run check
	npm run format:check

test:
	go test -race -cover ./cmd/... ./internal/...
	npm test

test-windows:
	go test -tags=windows_integration ./tests/windows/...

clean:
	go clean
	npm run clean --workspace @device-force/web
