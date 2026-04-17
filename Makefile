.PHONY: run desktop build build-desktop test release-dist npm-pack

run:
	go run -buildvcs=false ./cmd/server -config configs/config.example.yaml

desktop:
	go run -buildvcs=false ./cmd/desktop

build:
	go build -buildvcs=false ./cmd/server

build-desktop:
	go build -buildvcs=false -o opencode-desktop ./cmd/desktop

test:
	go test ./...

release-dist:
	node ./scripts/release.js

npm-pack:
	npm pack
