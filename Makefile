build:
	go build -o ${GOBIN}/google-workspace-mcp cmd/google-workspace-mcp/main.go

release-dry-run:
	goreleaser release --snapshot --clean --skip=publish

release-snapshot:
	goreleaser build --snapshot --clean