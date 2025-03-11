init:
	docker compose run --rm --no-deps frontend yarn install
	docker compose run --rm --no-deps backend yarn install

dev:
	docker compose up

.PHONY: init dev
