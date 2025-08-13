# Root Makefile for building all Lambda functions

.PHONY: all test fmt

# Default target
all: update-spec test

setup:
	@echo "Setting up project..."
	@go mod tidy

test: setup
	@echo "Running tests..."
	@go test ./... -v -race

fmt:
	@echo "Formatting Go code..."
	@gofmt -s -w -l .

update-spec:
	@echo "Downloading latest spec..."
	@curl -sSL -o pkg/isb/spec.yaml.new "https://raw.githubusercontent.com/aws-solutions/innovation-sandbox-on-aws/refs/heads/main/docs/openapi/innovation-sandbox-api.yaml"
	@echo "Comparing with existing spec..."
	@if [ ! -f pkg/isb/spec.yaml ]; then \
		echo "No existing spec found. Updating spec.yaml."; \
		mv pkg/isb/spec.yaml.new pkg/isb/spec.yaml; \
		echo "spec.yaml created."; \
		exit 0; \
	fi
	@awk '/^info:/ {p=1} p && /^  version:/ {next} {print}' pkg/isb/spec.yaml > pkg/isb/spec.yaml.old-noversion
	@awk '/^info:/ {p=1} p && /^  version:/ {next} {print}' pkg/isb/spec.yaml.new > pkg/isb/spec.yaml.new-noversion
	@if diff -q pkg/isb/spec.yaml.old-noversion pkg/isb/spec.yaml.new-noversion >/dev/null; then \
		echo "No changes to spec.yaml."; \
		rm pkg/isb/spec.yaml.new pkg/isb/spec.yaml.old-noversion pkg/isb/spec.yaml.new-noversion; \
	else \
		echo "Spec has changed (other than info.version). Updating spec.yaml."; \
		mv pkg/isb/spec.yaml.new pkg/isb/spec.yaml; \
		rm pkg/isb/spec.yaml.old-noversion pkg/isb/spec.yaml.new-noversion; \
	fi