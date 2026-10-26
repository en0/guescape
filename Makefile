.PHONY: watch-tests tests

tests:
	@go test ./...

watch-tests:
	@go test ./...
	@while inotifywait -r -e close_write --include '\.go$$' . 2>/dev/null; \
	do \
		echo ""; \
		go test ./...; \
	done
