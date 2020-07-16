
.PHONY: test
test:
	go test -v -race -cover github.com/replit/database-go
