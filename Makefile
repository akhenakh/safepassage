LDFLAGS_STATIC = -trimpath -ldflags "-linkmode external -extldflags -static"

targets = safepassage

.PHONY: all lint test clean testnorace testnolint safepassage safepassage-musl

all: test $(targets)

test: testnolint

testnorace:
	go test -v

testnolint:
	go test -race

lint:
	golangci-lint run

safepassage:
	CGO_ENABLED=0 go build -trimpath

safepassage-static: testnorace
	go build -a -v ${LDFLAGS_STATIC}

clean:
	rm -f safepassage

docker:
	docker build . -t akhenakh/safepassage:latest