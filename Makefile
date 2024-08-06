.PHONY: bump
bump:
	@echo "Bumping version"
	@git tag $(svu patch)
	@git push --tags