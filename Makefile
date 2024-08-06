.PHONY: bump
bump:
	@echo "🚀 Bumping version"
	git tag $(shell svu patch)
	git push --tags

.PHONY: release
release:
	@echo "🚀 Releasing version"
	goreleaser build --id bb --clean --snapshot --single-target --output dist/bb