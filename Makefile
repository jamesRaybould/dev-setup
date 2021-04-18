all: compile create-dist

compile:
	@echo "Making OSX fat binary"...
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o dev-setup-amd64 main.go
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o dev-setup-arm64 main.go

	lipo -create -output dev-setup dev-setup-amd64 dev-setup-arm64
	@echo "Cleaning up..."
	@rm dev-setup-amd64 dev-setup-arm64

run:
	go run main.go

create-dist:
	@rm -rf dist
	@echo "Producting dist"
	@mkdir -p dist
	@cp -r dev-setup config dist
	@echo "Distributable at ./dist"
