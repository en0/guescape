.PHONY: watch-tests tests build-release


tests:
	@go test ./...

watch-tests:
	@go test ./...
	@while inotifywait -r -e close_write --include '\.go$$' . 2>/dev/null; \
	do \
		echo ""; \
		go test ./...; \
	done

build-release:
	@GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o guescape_windows_amd64.exe .
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o guescape_linux_amd64 .
