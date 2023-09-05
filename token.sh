#!/bin/zsh

# amd64 (universal)
sqlite3 -line ~/Library/Messages/chat.db "SELECT m.text, m.date FROM message m WHERE text LIKE '%code%' ORDER BY m.ROWID DESC LIMIT 1" | ./bin/token-amd64-darwin
