.PHONY: test-db test-ride build install uninstall

# test if my terminal can access chat.db
test-db:
	sqlite3 -line ~/Library/Messages/chat.db "SELECT m.text, m.date FROM message m LIMIT 1"
	sqlite3 -line ~/Library/Messages/chat.db "SELECT m.text, m.date FROM message m ORDER BY m.ROWID DESC LIMIT 1"
	sqlite3 -line ~/Library/Messages/chat.db "SELECT m.text, m.date FROM message m WHERE text LIKE '%code%' ORDER BY m.ROWID DESC LIMIT 1"

# read token (this works only when run in repo root)
test-ride:
	./token.sh

build:
	GOOS=darwin GOARCH=amd64 go build -o bin/token-amd64-darwin main.go
	GOOS=darwin GOARCH=arm64 go build -o bin/token-arm64-darwin main.go

# use universal amd64 bin, it doesnt matter
install: build
	bin/token-amd64-darwin -install

uninstall:
	bin/token-amd64-darwin -uninstall
