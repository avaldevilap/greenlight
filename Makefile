run:
	go run ./cmd/api

# ==================================================================================== # # BUILD
# ==================================================================================== #
## build/api: build the cmd/api application
.PHONY: build/api
build/api:
	@echo 'Building cmd/api...'
	GOOS=darwin GOARCH=amd64 go build -ldflags='-s' -o=./bin/darwin_amd64/api ./cmd/api
	GOOS=darwin GOARCH=arm64 go build -ldflags='-s' -o=./bin/darwin_arm64/api ./cmd/api
	GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o=./bin/linux_amd64/api ./cmd/api