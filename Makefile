.PHONY: bump
bump:
	@echo "🚀 Bumping version"
	git tag $(shell svu patch)
	git push --tags