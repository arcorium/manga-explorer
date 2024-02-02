.PHONY: build
build:
	@go build -o build/server "./cmd/server/"
	@go build -o build/migrate "./cmd/migrate/"

.PHONY: run
run: build/server
	@export GIN_MODE=release; \
	./build/server

.PHONY: migrate
migrate: build/migrate
	@./build/migrate --env --special --no-ssl

test.prepare: build/migrate
	@./build/migrate --env --no-ssl
	@./build/migrate --env --seed database/fixtures --no-ssl
	@./build/migrate --env --special --no-ssl

.PHONY: clean
clean:
	@go clean
	@rm -fdR build # Doesn't works on windows