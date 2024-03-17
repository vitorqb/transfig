.PHONY: default
default: test format

.PHONY: test
test:
	docker compose run --rm test

# test-debug runs golang delve in test mode
.PHONY: test-debug
test-debug:
	docker compose run --rm -it test-debug

.PHONY: format
format:
	docker compose run --rm format
	docker compose run --rm lint
