# test if my terminal can access chat.db
.PHONY: test-db
test-db:
	sqlite3 -line ~/Library/Messages/chat.db "SELECT m.text, m.date FROM message m LIMIT 1"
	sqlite3 -line ~/Library/Messages/chat.db "SELECT m.text, m.date FROM message m ORDER BY m.ROWID DESC LIMIT 1"
	sqlite3 -line ~/Library/Messages/chat.db "SELECT m.text, m.date FROM message m WHERE text LIKE '%code%' ORDER BY m.ROWID DESC LIMIT 1"

# read token (this works only when run in repo root)
.PHONY: test-ride
test-ride:
	./token.sh

.PHONY: build
build:
	GOOS=darwin GOARCH=amd64 go build -o bin/token-amd64-darwin main.go
	GOOS=darwin GOARCH=arm64 go build -o bin/token-arm64-darwin main.go

# use universal amd64 bin, it doesnt matter
.PHONY: install
install: build
	bin/token-amd64-darwin -install

.PHONY: uninstall
uninstall:
	bin/token-amd64-darwin -uninstall

# update modules to latest minor versions
.PHONY: env
env:
	go get -u
	go mod tidy
	go mod download

.PHONY: test
test:
	go test ./pkg/sms/
	go test ./pkg/symlinks/

.PHONY: lint
lint:
	go fmt
	go vet
	golangci-lint run
