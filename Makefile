.PHONY: bump
bump:
	@echo "🚀 Bumping Version"
	git tag $(shell svu patch)
	git push --tags

.PHONY: build
build:
	@echo "🚀 Building Version $(shell svu current)"
	go build -o bb main.go

.PHONY: release
release:
	@echo "🚀 Releasing Version $(shell svu current)"
	goreleaser build --id bb --clean --snapshot --single-target --output dist/bb