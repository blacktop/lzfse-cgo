ifneq ($(filter check-release bump release,$(MAKECMDGOALS)),)
NEXT_VERSION:=$(shell svu patch)
endif

.PHONY: test
test: ## Run tests
	@echo " > Running tests"
	@go test -v ./...

.PHONY: check-release
check-release: ## Check the tree is clean, on master, in sync with origin, and passing tests
	@test -n "${NEXT_VERSION}" || { echo "NEXT_VERSION is empty: install svu or pass it" >&2; exit 1; }
	@test -z "$$(git status --porcelain)" || { echo "working tree is dirty" >&2; exit 1; }
	@test "$$(git branch --show-current)" = master || { echo "not on master" >&2; exit 1; }
	@git fetch --quiet --tags origin master
	@test "$$(git rev-parse HEAD)" = "$$(git rev-parse origin/master)" || \
		{ echo "HEAD differs from origin/master: pull or push first" >&2; exit 1; }
	@tag="$$(git tag --points-at HEAD)"; test -z "$$tag" || \
		{ echo "HEAD is already tagged $$tag: publish or delete it, do not bump again" >&2; exit 1; }
	@go test ./...

.PHONY: bump
bump: check-release ## Tag and push NEXT_VERSION (override with NEXT_VERSION=vX.Y.Z)
	@echo " > Tagging ${NEXT_VERSION}"
	@git tag -a ${NEXT_VERSION} -m "Release ${NEXT_VERSION}"
	@git push origin ${NEXT_VERSION}

.PHONY: release
release: bump ## Run bump, then publish a GitHub release with generated notes
	@echo " > Creating Release ${NEXT_VERSION}"
	@gh release create ${NEXT_VERSION} --verify-tag --generate-notes

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
