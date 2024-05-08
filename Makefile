.PHONY: build migrate swagger clean run
build:
	@go build -o build/server "./cmd/server/"
	@go build -o build/migrate "./cmd/migrate/"

run: build/server
	@export GIN_MODE=release; \
	./build/server

migrate: build/migrate
	@./build/migrate --env --special --no-ssl

test.prepare: build/migrate
	@./build/migrate --env --no-ssl
	@./build/migrate --env --seed database/fixtures --no-ssl
	@./build/migrate --env --special --no-ssl

clean:
	@go clean
	@rm -fdR build # Doesn't works on windows


swagger:
	@wag init -d .\cmd\server\,.\internal\