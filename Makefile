.PHONY: create_table

create_table:
	@go run cmd/create_table/create_table.go | gofmt > eucjis2004/table.go

test:
	@go test ./...
