# CLI has been removed. This Makefile is kept for potential future Go tooling.

.PHONY: help

## help: Show this help
help:
	@echo "Available targets:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'
