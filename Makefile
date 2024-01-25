.PHONY: build
build:
	go build -o build/server "./cmd/server/"
	go build -o build/migrate "./cmd/migrate/"

.PHONY: test.prepare
test.prepare: ./build/migrate
	./build/migrate

.PHONY: clean
clean:
	go clean
	rm -fdR build