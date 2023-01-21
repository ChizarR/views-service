run:
	go run ./cmd/main/.

build:
	go build -o ./build/service ./cmd/main/.

run_prod:
	build
	./build/service
