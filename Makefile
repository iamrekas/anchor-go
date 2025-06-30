build:
	mkdir ./bin/ 2>/dev/null || true
	go build -o ./bin/anchor-go ./cmd/anchor-go

install:
	go install ./cmd/anchor-go

dummy:
	$(MAKE) build && \
	./bin/anchor-go -src=./example/dummy_idl.json -pkg=dummy -dst=./generated/dummy && \
	go test ./generated/dummy && \
	go test ./example/dummy_test.go

restaking:
	$(MAKE) build && \
	./bin/anchor-go -src=./example/restaking_idl.json -pkg=restaking -dst=./generated/restaking && \
	go test ./generated/restaking && \
	go test ./example/restaking_test.go

test:
	$(MAKE) dummy && $(MAKE) restaking

clean:
	rm -f anchor-go
	rm -rf ./generated/dummy
	rm -rf ./generated/restaking

.PHONY: build install dummy restaking test clean
