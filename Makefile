.PHONY: default
default: test format

.PHONY: test
test:
	docker compose run --rm test

.PHONY: format
format:
	docker compose run --rm format
	docker compose run --rm lint
