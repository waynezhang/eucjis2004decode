.PHONY: create_table

create_table:
	@go run cmd/create_table/create_table.go | gofmt > table/table.go

test:
	@go test ./...
