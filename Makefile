.PHONY: gen-proto
gen-proto:
	buf build && buf generate

.PHONY: lint-proto
lint-proto:
	buf lint

.PHONY: test
test:
	go test -v -cover -race -count 1 ./...