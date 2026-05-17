BIN := miniCC

.PHONY: build clean run

build:
	go build -o bin/$(BIN) ./cmd/miniCC

run: build
	./bin/$(BIN)

clean:
	rm -rf bin/
