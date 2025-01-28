.PHONY: create_table test

create_table:
	@go run cmd/create_table/create_table.go | gofmt > eucjis2004/table.go

test: create_table
	@go test ./...
