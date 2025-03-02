.PHONY: install-tools

install-tools:
	@echo "Installing templ..."
	@go install github.com/a-h/templ/cmd/templ@latest
	@echo "Installing air..."
	@go install github.com/air-verse/air@latest
	@echo "Tools installed successfully"