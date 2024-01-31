.PHONY: build
build:
	go build -o build/server "./cmd/server/"
	go build -o build/migrate "./cmd/migrate/"

test.prepare: ./build/migrate
	./build/migrate --env test --no-ssl
	./build/migrate --env test --seed database/fixtures --no-ssl
	./build/migrate --env test --special --no-ssl

.PHONY: clean
clean:
	go clean
	rm -fdR build # Doesn't works on windows